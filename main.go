package main

import (
	"fmt"
	"os"

	"gopud/cmd/root"
	"gopud/internal/config"
)

var SecKey string

func realMain() int {
	root.Execute()
	return 0
}

func main() {
	// go build -v -ldflags "-X 'main.SecKey=$(go run ./cmd/genkey/main.go)'"

	if SecKey == "" {
		fmt.Println("Read main.go for learn build !!!")
	}

	config.SECKey = SecKey
	os.Exit(realMain())
}
