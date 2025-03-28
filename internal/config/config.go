package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"connection_user_name"`
}

func Read() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	file := filepath.Join(homeDir, configFileName)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil
	}
	return &config
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	jsonData, err := json.Marshal(*c)
	if err != nil {
		return fmt.Errorf("Error marshalling config: %v", err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting user home directory: %v", err)
	}

	err = os.WriteFile(filepath.Join(homeDir, configFileName), jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing config: %v", err)
	}
	return nil
}
