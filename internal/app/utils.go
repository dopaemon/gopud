package app

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func RenderGlamour(text string) string {
	out, err := glamour.Render(text, "dark")
	if err != nil {
		return text
	}
	return out
}

func PrintGlamourError(err error) error {
	fmt.Println(RenderGlamour(fmt.Sprintf("❌ Error: %s", err.Error())))
	return err
}
