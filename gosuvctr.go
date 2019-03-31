package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	version   string
	config    Configuration
	key       string
	remoteURL string
)

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
		remoteURL = "http://" + config.RemoteAddr + ":" + config.RemotePort
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

func getKey() {
	plainKey := config.Admin.Username + ":" + config.Admin.Password
	cryptoKey := base64.StdEncoding.EncodeToString([]byte(plainKey))
	key = "Basic " + cryptoKey
}
