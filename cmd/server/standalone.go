package main

// This includes the code used when running timeruler
// as a standalone application on your local machine.
// When this is the case, the serverUrl in the config
// file will be `http://localhost:6756`

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"

	"github.com/muesli/go-app-paths"
)

func SetPid(address, port string) error {
	scope := gap.NewScope(gap.User, "timeruler")
	configPath, err := scope.ConfigPath("config.json")
	if err != nil {
		return err
	}
	config := make(map[string]string)
	configFile, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil && os.IsNotExist(err) {
		configFile, err = os.Create(configPath)
		if err != nil {
			return err
		}
		configFile.Close()
	} else if err != nil {
		return err
	} else {
		configBytes, err := io.ReadAll(configFile)
		if err != nil {
			configFile.Close()
			return err
		}
		configFile.Close()

		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}
	}

	pid := syscall.Getpid()
	if _, ok := config["pid"]; ok {
		return errors.New(fmt.Sprintf("Config file shows server already running at pid %s", config["pid"]))
	} else {
		config["pid"] = strconv.Itoa(pid)
	}

	if _, ok := config["server"]; !ok {
		config["server"] = address + ":" + port
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	configFile, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()
	_, err = configFile.Write(configBytes)
	if err != nil {
		return err
	}

	return nil
}
