package main

import (
	"os"
	"path"
	"os/user"
	"io/ioutil"
)

type Config struct {
	username string
	path string
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	config.path = configPath

	buf, e := ioutil.ReadFile(config.path)
	if e != nil{
		return nil, e
	}

	config.username = string(buf)
	return &config, nil
}

func PromptConfig(configPath string) (*Config, error) {
	var config Config
	config.path = configPath

	username, e := PromptUsername()
	if e == nil {
		e = ioutil.WriteFile(configPath, []byte(username), 0644)
	}
	if e != nil{
		return nil, e
	}

	config.username = username
	return &config, nil
}

func LoadOrPromptConfig() (*Config, error) {
	usr, e := user.Current()
	if e != nil {
		return nil, e
	}

	configPath := path.Join(usr.HomeDir, ".chatconfig")
	config, e := LoadConfig(configPath)

	if os.IsNotExist(e) {
		config, e = PromptConfig(configPath)
	}
	if e != nil {
		return nil, nil
	}

	return config, nil
}
