package service

import (
	"encoding/json"
	"os"
)

type Config struct {
	Login             string   `json:"login"`
	Password          string   `json:"password"`
	Port              string   `json:"port"`
	StoragePath       string   `json:"storage_path"`
	AllowedExtensions []string `json:"allowed_extensions"`
	AuthEnabled       bool     `json:"auth_enabled"`
}

func GetConfig() Config {
	file, err := os.Open("config/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}
