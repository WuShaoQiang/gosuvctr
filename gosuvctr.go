package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

var (
	version   string
	config    Configuration
	key       string
	remoteUrl string
)

type Configuration struct {
	Admin      administrator `json:"admin"`
	RemoteAddr string        `json:"remoteAddr"`
	RemotePort string        `json:"remotePort"`
}

type administrator struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JSONResponse struct {
	Status int         `json:"status"`
	Value  interface{} `json:"value"`
}

type Program struct {
	Name         string   `yaml:"name" json:"name"`
	Command      string   `yaml:"command" json:"command"`
	Environ      []string `yaml:"environ" json:"environ"`
	Dir          string   `yaml:"directory" json:"directory"`
	StartAuto    bool     `yaml:"start_auto" json:"startAuto"`
	StartRetries int      `yaml:"start_retries" json:"startRetries"`
	StartSeconds int      `yaml:"start_seconds,omitempty" json:"startSeconds"`
	StopTimeout  int      `yaml:"stop_timeout,omitempty" json:"stopTimeout"`
	retryCount   int
	User         string `yaml:"user,omitempty" json:"user"`
	// Notifications Notifications `yaml:"notifications,omitempty" json:"-"`
	// WebHook       WebHook       `yaml:"webhook,omitempty" json:"-"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	version = "master"
}

func main() {
	defaultConfigPath := filepath.Join(DefaultCtrDir(), "conf/config.json")

	app := cli.NewApp()
	app.Name = "gosuv controller"
	app.Usage = "Controller for gosuv administrator"
	app.Version = version
	app.Before = func(c *cli.Context) error {
		var err error
		configPath := c.GlobalString("conf")
		config, err = readConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}
		getKey()
		remoteUrl = "http://" + config.RemoteAddr + ":" + config.RemotePort
		return nil
	}

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Gavin",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
			Value: defaultConfigPath,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "shutdown",
			Usage:  "Shutdown the remote gosuv",
			Action: shutdownGosuv,
		},
		{
			Name:   "status",
			Usage:  "Check remote gosuv status",
			Action: statusGosuv,
		},
		{
			Name:   "start",
			Usage:  "Start remote gosuv program",
			Action: startProgram,
		},
		{
			Name:   "stop",
			Usage:  "Stop remote gosuv program",
			Action: stopProgram,
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func readConfig(configPath string) (c Configuration, err error) {
	c.RemoteAddr = "127.0.0.1"
	c.RemotePort = "11313"

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getKey() {
	plainKey := config.Admin.Username + ":" + config.Admin.Password
	cryptoKey := base64.StdEncoding.EncodeToString([]byte(plainKey))
	key = "Basic " + cryptoKey
}

func shutdownGosuv(c *cli.Context) error {
	ret, err := postFormWithAuth(remoteUrl+"/api/shutdown", nil, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret.Value)
	return nil
}

func statusGosuv(c *cli.Context) error {
	resp, err := getWithAuth(remoteUrl+"/api/programs", key)
	if err != nil {
		log.Fatal(err)
	}
	var programs = make([]struct {
		Program Program `json:"program"`
		Status  string  `json:"status"`
	}, 0)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &programs)
	if err != nil {
		log.Fatal(err)
	}
	format := "%-23s\t%-8s\n"
	fmt.Printf(format, "PROGRAM NAME", "STATUS")
	for _, p := range programs {
		fmt.Printf(format, p.Program.Name, p.Status)
	}
	return nil
}

func startProgram(c *cli.Context) error {
	programName := c.Args().First()

	success, err := programOperator("start", programName)
	if err != nil {
		log.Fatal(err)
	}
	if success {
		fmt.Println(programName + "   started")
	} else {
		fmt.Println(programName + "   started failed")
	}

	return nil
}

func stopProgram(c *cli.Context) error {
	programName := c.Args().First()

	success, err := programOperator("stop", programName)
	if err != nil {
		log.Fatal(err)
	}
	if success {
		fmt.Println(programName + "   stopped")
	} else {
		fmt.Println(programName + "   stopped failed")
	}

	return nil
}

func programOperator(cmd, programName string) (success bool, err error) {
	resp, err := postWithAuth(remoteUrl+"/api/programs/"+programName+"/"+cmd, "", nil, key)
	if err != nil {
		log.Fatal(err)
	}

	var v = struct {
		Status int `json:"status"`
	}{}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Fatal(err)
	}
	success = (v.Status == 0)
	return
}

func postFormWithAuth(url string, data url.Values, key string) (r JSONResponse, err error) {
	resp, err := postWithAuth(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), key)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, fmt.Errorf("POST %v %v", strconv.Quote(url), string(body))
	}
	return r, nil
}

func postWithAuth(url, contentType string, body io.Reader, key string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", key)
	return http.DefaultClient.Do(req)
}

func getWithAuth(url string, key string) (resp *http.Response, err error) {
	// url := remoteUrl + pathname
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", key)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func DefaultCtrDir() string {
	return filepath.Join(UserHomeDir(), ".gosuvctr")
}
