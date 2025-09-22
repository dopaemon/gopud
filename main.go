package main

import (
	"fmt"

	"gopud/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config file not exist !!!")
	}

	if cfg.APIKey == "" {
		cfg.APIKey = "MyAPIKey"
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Can't create config File !!!")
		}
	}

	fmt.Println("Hello World !!!")
}
