package main

import (
	"os"
	"path"
	"os/user"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Username string `json:"username"`
	Path string `json:"-"`
}

func LoadConfig(configPath string) (*Config, error) {
	buf, e := ioutil.ReadFile(configPath)
	if e != nil {
		return nil, e
	}

	var config Config
	if e = json.Unmarshal(buf, &config); e != nil {
		return nil, e
	}

	config.Path = configPath
	return &config, nil
}

func PromptConfig(configPath string) (*Config, error) {
	var config Config
	config.Path = configPath

	username, e := PromptUsername()
	if e == nil {
		config.Username = username
		buf, _ := json.Marshal(config)
		e = ioutil.WriteFile(configPath, buf, 0644)
	}

	if e != nil{
		return nil, e
	}
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
		return nil, e
	}
	return config, nil
}
