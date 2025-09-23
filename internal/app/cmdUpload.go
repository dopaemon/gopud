package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopud/internal/config"

	"github.com/ManuelReschke/go-pd/pkg/pd"
	"github.com/imroc/req/v3"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	_ "github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func RunUpload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return PrintGlamourError(errors.New("please add a file to your upload request"))
	}

	verboseFlag := cmd.Flags().Changed("verbose")

	for idx, file := range args {
		if _, err := os.Stat(filepath.FromSlash(file)); errors.Is(err, os.ErrNotExist) {
			return PrintGlamourError(errors.New("one of the given files does not exist"))
		}

		r := &pd.RequestUpload{
			PathToFile: file,
			Anonymous:  false,
			Auth: pd.Auth{
				APIKey: config.APIRawKey,
			},
		}

		c := pd.New(nil, nil)
		var bar *progressbar.ProgressBar
		startTime := time.Now()
		lastUpdate := time.Now()

		if verboseFlag {
			bar = progressbar.NewOptions(
				100,
				progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetWidth(30),
				progressbar.OptionSetRenderBlankState(true),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}),
			)

			c.SetUploadCallback(func(info req.UploadInfo) {
				if info.FileSize == 0 {
					return
				}

				progress := int(float64(info.UploadedSize) / float64(info.FileSize) * 100)
				bar.Set(progress)

				if time.Since(lastUpdate) >= 100*time.Millisecond {
					lastUpdate = time.Now()
					elapsed := time.Since(startTime).Seconds()
					if elapsed < 0.001 {
						elapsed = 0.001
					}
					speed := formatBytes(float64(info.UploadedSize)/elapsed) + "/s"
					bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset] Uploading %s... %s", idx+1, len(args), file, speed))
				}
			})
		}

		rsp, err := c.UploadPOST(r)
		if err != nil {
			return PrintGlamourError(err)
		}

		if verboseFlag && bar != nil {
			bar.Finish()
			fmt.Println()
			time.Sleep(50 * time.Millisecond)
		}

		msg := fmt.Sprintf("Uploaded [%d/%d] %s | URL: %s", idx+1, len(args), file, rsp.GetFileURL())
		fmt.Println(RenderGlamour(msg))
	}

	return nil
}

func formatBytes(b float64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", b/KB)
	default:
		return fmt.Sprintf("%.0f B", b)
	}
}

