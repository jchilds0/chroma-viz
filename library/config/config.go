package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conn struct {
	Name    string
	Address string
	Port    int
	Type    string
}

type Config struct {
	HubAddr          string
	HubPort          int
	PreviewDirectory string
	PreviewName      string
	Connections      []Conn
}

func ImportConfig(file string) (*Config, error) {
	var conf Config

	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file (%s)", err)
	}

	err = json.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file (%s)", err)
	}

	return &conf, nil
}
