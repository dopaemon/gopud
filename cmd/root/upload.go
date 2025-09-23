package root

import (
	"gopud/internal/app"

	"github.com/spf13/cobra"
)

const (
	cmdUploadUse   = "upload"
	cmdUploadShort = "With that command you can upload files"
)

var uploadCmd = &cobra.Command{
	Use:   cmdUploadUse,
	Short: cmdUploadShort,
	RunE:  app.RunUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolP("verbose", "v", true, "Show more information after an upload (Anonymous, ID, URL)")
}
