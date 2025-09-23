package main

import (
	"fmt"
	"os"

	"encoding/base64"


	"gopud/cmd"
	"gopud/internal/config"
	"gopud/internal/security"
	"gopud/internal/utils"

	"github.com/charmbracelet/huh"
)

var (
	YouAPIKey	string
	tokenInput	string
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config file not exist !!!")
	}

	if config.SECKey == "" {
		// config.SECKey = "vSECKEY"
		config.SECKey = "12345678901234569876432126789013"
	}

	if len(cfg.APIKey) < 5 || utils.VerifyPixelDrainAPIKey(GetAPIKey(cfg.APIKey)) == false {
		for {
			form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
						Title("PixelDrain API KEY: ").
						Placeholder("Keep null, or type exit for Exit.").
						Password(true).
						Value(&tokenInput),
					),
				)

			if err = form.Run(); err != nil {
				if err == huh.ErrUserAborted {
					fmt.Println(err)
					os.Exit(0)
				}
				fmt.Println(err)
				return
			}

			if tokenInput == "" || tokenInput == "exit" { os.Exit(0) }
			if utils.VerifyPixelDrainAPIKey(tokenInput) { break }
			tokenInput = ""
		}

		types64, err := security.EncryptData([]byte(tokenInput), []byte(config.SECKey))
		if err != nil {
			fmt.Println("Base64 APIKey Error: ", err)
		}
		cfg.APIKey = base64.StdEncoding.EncodeToString(types64)
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Can't create config File !!!")
		}
	}

	config.APIRawKey = GetAPIKey(cfg.APIKey)
	cmd.Execute()
}

func GetAPIKey(apikey string) string {
	encrypt, err := base64.StdEncoding.DecodeString(apikey)
	if err != nil {
		fmt.Println("Can't Decode Base64 API Key")
	}
	decrypted, err := security.DecryptData(encrypt, []byte(config.SECKey))
	if err != nil {
		fmt.Println("Can't Decode Byte Code, Rerun gopud")
		_ = config.DeleteConfig()
		os.Exit(1)
	}

	return string(decrypted)
}
