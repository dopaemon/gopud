package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"gopud/internal/config"
	"gopud/internal/security"
	"gopud/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:   "gopud [flags]",
	Short: "Run gopud Program !!!",
	Long: `PixelDrain Download / Upload CLI

It doesn’t really do anything, but that’s the point.™`,
	Example: `
# Run it:
gopud

# Run it with setup API Key:
gopud api "your_api_key"

# Run it with short argument:
gopud upload "example.zip image.jpg"

# Run it with arguments:
gopud download "aaeb11j"
`,
}

var (
	tokenInput string
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "api [API_KEY]",
		Short: "Set PixelDrain API Key",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var key string
			if len(args) == 1 {
				key = args[0]
			} else {
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("PixelDrain API KEY: ").
							Placeholder("Keep null, or type exit for Exit.").
							Password(true).
							Value(&tokenInput),
					),
				)
				if err := form.Run(); err != nil || tokenInput == "" || tokenInput == "exit" {
					fmt.Println("Hủy thiết lập API Key")
					os.Exit(0)
				}
				key = tokenInput
			}

			if err := saveAPIKey(key); err != nil {
				fmt.Println("Lỗi lưu API Key:", err)
				os.Exit(1)
			}
			fmt.Println("API Key saved successfully!")
		},
	})
}

func saveAPIKey(key string) error {
	if !utils.VerifyPixelDrainAPIKey(key) {
		return fmt.Errorf("invalid API key")
	}

	if config.SECKey == "" {
		config.SECKey = "12345678901234569876432126789013"
	}

	encrypted, err := security.EncryptData([]byte(key), []byte(config.SECKey))
	if err != nil {
		return err
	}

	cfg, _ := config.LoadConfig()
	cfg.APIKey = base64.StdEncoding.EncodeToString(encrypted)

	return config.SaveConfig(cfg)
}

func ensureAPIKey() {
	cfg, _ := config.LoadConfig()

	if config.SECKey == "" {
		config.SECKey = "12345678901234569876432126789013"
	}

	if len(cfg.APIKey) < 5 || !utils.VerifyPixelDrainAPIKey(GetAPIKey(cfg.APIKey)) {
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

			if err := form.Run(); err != nil || tokenInput == "" || tokenInput == "exit" {
				fmt.Println("Thoát chương trình")
				os.Exit(0)
			}

			if err := saveAPIKey(tokenInput); err == nil {
				break
			}
			tokenInput = ""
		}
		cfg, _ = config.LoadConfig()
	}

	config.APIRawKey = GetAPIKey(cfg.APIKey)
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

func Execute() {
	if len(os.Args) > 1 && os.Args[1] == "api" {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	}

	ensureAPIKey()

	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
