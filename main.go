package main

import (
	"fmt"

	"encoding/base64"


	"gopud/internal/flags"
	"gopud/internal/config"
	"gopud/internal/security"
)

var (
	YouAPIKey	string
)

func main() {
	_, _ = flags.Flags()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config file not exist !!!")
	}

	if cfg.APIKey == "" {
		if config.SECKey == "" {
			// config.SECKey = "vSECKEY"
			config.SECKey = "12345678901234567890123456789012"
		}
		types64, err := security.EncryptData([]byte("Hello World !!!"), []byte(config.SECKey))
		if err != nil {
			fmt.Println("Base64 APIKey Error: ", err)
		}
		cfg.APIKey = base64.StdEncoding.EncodeToString(types64)
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Can't create config File !!!")
		}
	}

	fmt.Println("Hello World !!!")
}
