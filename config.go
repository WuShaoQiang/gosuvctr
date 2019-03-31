package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Configuration record the config.json data
type Configuration struct {
	Admin      administrator `json:"admin"`
	RemoteAddr string        `json:"remoteAddr"`
	RemotePort string        `json:"remotePort"`
}

type administrator struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserHomeDir return $HOME directory
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

// DefaultCtrDir return the config.json file's filepath
func DefaultCtrDir() string {
	return filepath.Join(UserHomeDir(), ".gosuvctr")
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
