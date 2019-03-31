package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/urfave/cli"
)

// JSONResponse is for
type JSONResponse struct {
	Status int         `json:"status"`
	Value  interface{} `json:"value"`
}

// Program is for
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

func shutdownGosuv(c *cli.Context) error {
	ret, err := postFormWithAuth(remoteURL+"/api/shutdown", nil, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret.Value)
	return nil
}

func statusGosuv(c *cli.Context) error {
	resp, err := getWithAuth(remoteURL+"/api/programs", key)
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
	resp, err := postWithAuth(remoteURL+"/api/programs/"+programName+"/"+cmd, "", nil, key)
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
