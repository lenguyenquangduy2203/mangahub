package config

import (
	"encoding/json"
	"log"
	"os"
)

type ClientConfig struct {
	ServerURL string `json:"server_url"`
	WsURL     string `json:"ws_url"`
}

func LoadClientConfig() *ClientConfig {
	// Defaults
	cfg := &ClientConfig{
		ServerURL: "http://localhost:3000/api/v1",
		WsURL:     "ws://localhost:3000/api/v1/ws",
	}

	file, err := os.Open("client_config.json")
	if err == nil {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("Error closing config file: %v", err)
			}
		}(file)
		_ = json.NewDecoder(file).Decode(cfg)
	}

	return cfg
}
