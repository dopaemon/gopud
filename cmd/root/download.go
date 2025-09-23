package root

import (
	"gopud/internal/app"

	"github.com/spf13/cobra"
)

const (
	cmdDownloadUse   = "download"
	cmdDownloadShort = "With that command you can download a file"
)

var downloadCmd = &cobra.Command{
	Use:   cmdDownloadUse,
	Short: cmdDownloadShort,
	RunE:  app.RunDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("path", "p", "", "Path where the files are stored")
	downloadCmd.Flags().BoolP("verbose", "v", true, "Show more information after an upload (Anonymous, ID, URL)")
}
