package utils

import (
	"encoding/base64"
	"net/http"
)

func VerifyPixelDrainAPIKey(apiKey string) bool {
	url := "https://pixeldrain.com/api/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	token := base64.StdEncoding.EncodeToString([]byte(":" + apiKey))
	req.Header.Set("Authorization", "Basic "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
