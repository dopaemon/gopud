package main

import (
	"fmt"
	"os"

	"gopud/cmd/root"
	"gopud/internal/config"
)

var SecKey string = "58398300406949656740850173328796"

func realMain() int {
	root.Execute()
	return 0
}

func main() {
	if SecKey == "" {
		fmt.Println("Read main.go for learn build !!!")
	}

	config.SECKey = SecKey
	os.Exit(realMain())
}
