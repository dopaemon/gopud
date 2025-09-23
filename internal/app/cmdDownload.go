package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopud/internal/config"

	"github.com/ManuelReschke/go-pd/pkg/pd"
	_ "github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func RunDownload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return PrintGlamourError(errors.New("please add a pixeldrain URL or file id to your download request"))
	}

	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return PrintGlamourError(errors.New("please add a valid path where you want to save the files"))
	}
	if path == "" {
		path, _ = os.Getwd()
	}

	apiKey := config.APIRawKey

	verboseFlag := cmd.Flags().Changed("verbose")

	for idx, file := range args {
		fileID := file
		if strings.ContainsAny(file, pd.BaseURL) {
			fileID = filepath.Base(file)
		}

		reqInfo := &pd.RequestFileInfo{
			ID: fileID,
		}
		if apiKey != "" {
			reqInfo.Auth.APIKey = apiKey
		}

		c := pd.New(nil, nil)
		rspInfo, err := c.GetFileInfo(reqInfo)
		if err != nil {
			return PrintGlamourError(err)
		}

		reqDL := &pd.RequestDownload{
			ID:         fileID,
			PathToSave: filepath.FromSlash(path + "/" + rspInfo.Name),
		}
		if apiKey != "" {
			reqDL.Auth.APIKey = apiKey
		}

		startTime := time.Now()
		rspDL, err := c.Download(reqDL)
		if err != nil {
			return PrintGlamourError(err)
		}

		msg := ""
		if rspDL.Success {
			if verboseFlag {
				elapsed := time.Since(startTime).Round(time.Millisecond)
				msg = fmt.Sprintf("Downloaded [%d/%d] %s | ID: %s | Stored to: %s | Time: %s",
					idx+1, len(args), rspDL.FileName, reqDL.ID, reqDL.PathToSave, elapsed)
			} else {
				msg = fmt.Sprintf("%s", reqDL.PathToSave)
			}
		} else {
			msg = fmt.Sprintf("Failed! ID: %s | Value: %s | Message: %s", reqDL.ID, rspDL.Value, rspDL.Message)
		}

		fmt.Println(RenderGlamour(msg))
	}

	return nil
}
