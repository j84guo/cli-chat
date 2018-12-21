package main

import (
	"os"
	"fmt"
	"path"
	"os/user"
	"strings"
	"io/ioutil"
)

type Config struct {
	username string
	path string
}

func askUsername() (string, error) {
	return promptLine("Enter a username:\n")
}

func promptLine(msg string) (string, error) {
	var line string
	fmt.Print(msg)
	if _, e := fmt.Scanln(&line); e != nil {
		return "", e
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func loadConfig() (*Config, error) {
	user, e := user.Current()
	if e != nil {
		return nil, e
	}

	var config Config
	config.path = path.Join(user.HomeDir, ".chatconfig")

	buf, e := ioutil.ReadFile(config.path)
	if e != nil{
		return nil, e
	}

	config.username = string(buf)
	return &config, nil
}

func makeConfig() (*Config, error) {
	user, e := user.Current()
	if e != nil {
		return nil, e
	}

	var config Config
	config.path = path.Join(user.HomeDir, ".chatconfig")
	config.username, e = askUsername()
	if e == nil {
		e = ioutil.WriteFile(config.path, []byte(config.username), 0644)
	}

	if e != nil{
		return nil, e
	}
	return &config, nil
}

func GetConfig() (*Config, error) {
	config, e := loadConfig()
	if os.IsNotExist(e) {
		config, e = makeConfig()
	}
	if e != nil {
		return nil, nil
	}
	return config, nil
}
