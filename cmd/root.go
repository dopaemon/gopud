package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
                Use:   "gopud [flags]",
                Short: "Run gopud Program !!!",
                Long: `PixelDrain Download / Upload CLI

It doesn’t really do anything, but that’s the point.™`,
                Example: `
# Run it:
gopud

# Run it with short argument:
gopud upload "example.zip image.jpg"

# Run it with arguments:
gopud download "aaeb11j"
`,
}

func Execute() {
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
