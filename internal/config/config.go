package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	filepath := getConfigFilePath()
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var config Config
	err = json.Unmarshal(configData, &config)

	if err != nil {
		log.Fatalf("Error Umarshaling: %v", err)
	}

	return config
}

func (c *Config) SetUser(currentUser string) {
	c.CurrentUserName = currentUser
	err := write(*c)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

func getConfigFilePath() (filePath string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error finding home directory: %v", err)
	}
	filePath = homeDir + "/" + configFileName
	return filePath
}

func write(cfg Config) error {
	jsonData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling: %v", err)
	}
	filepath := getConfigFilePath()
	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	return err
}
