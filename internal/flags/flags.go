package flags

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func Flags() (string, string) {
	var upload string
	var download string

	cmd := &cobra.Command{
		Use:   "gopud [flags]",
		Short: "Run gopud Program !!!",
		Long: `PixelDrain Download / Upload CLI

It doesn’t really do anything, but that’s the point.™`,
		Example: `
# Run it:
gopud

# Run it with short argument:
gopud --upload "example.zip image.jpg"

# Run it with arguments:
gopud --download "aaeb11j"
`,
		RunE: func(c *cobra.Command, args []string) error {
			if c.Flags().Changed("upload") {
				return fmt.Errorf("upload")
			} else if c.Flags().Changed("download") {
				return fmt.Errorf("download")
			} else {
				c.Println("You ran the root command. Now try --help.")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&upload, "upload", "u", "", "Upload file to PixelDrain")
	cmd.Flags().StringVarP(&download, "download", "d", "", "Download file from PixelDrain")

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.CompletionOptions.DisableDefaultCmd = true

	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return download, upload
}
