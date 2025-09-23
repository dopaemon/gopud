package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopud/internal/config"

	"github.com/ManuelReschke/go-pd/pkg/pd"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

func RunUpload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please add a file to your upload request")
	}

	verboseFlag := cmd.Flags().Changed("verbose")

	for _, file := range args {
		if _, err := os.Stat(filepath.FromSlash(file)); errors.Is(err, os.ErrNotExist) {
			return errors.New("one of the given files does not exist")
		}

		r := &pd.RequestUpload{
			PathToFile: file,
			Anonymous:  true,
		}

		r.Anonymous = false
		r.Auth.APIKey = config.APIRawKey

		c := pd.New(nil, nil)
		if verboseFlag {
			c.SetUploadCallback(func(info req.UploadInfo) {
				if info.FileSize > 0 {
					fmt.Printf("%q uploaded %.2f%%\n", info.FileName, float64(info.UploadedSize)/float64(info.FileSize)*100.0)
				} else {
					fmt.Printf("%q uploaded 0%% (file size is zero)\n", info.FileName)
				}
			})
		}
		rsp, err := c.UploadPOST(r)
		if err != nil {
			return err
		}

		msg := ""
		if verboseFlag {
			msg = fmt.Sprintf("Successful! Anonymous upload: %v | ID: %s | URL: %s", r.Anonymous, rsp.ID, rsp.GetFileURL())
		} else {
			msg = fmt.Sprintf("%s", rsp.GetFileURL())
		}

		fmt.Println(msg)
	}

	return nil
}
