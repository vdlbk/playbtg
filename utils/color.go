package utils

import (
	"fmt"
)

const (
	ColorTermRed    = "\033[31m"
	ColorTermGreen  = "\033[32m"
	ColorTermYellow = "\033[33m"
	ColorTermBlue   = "\033[34m"
	ColorTermReset  = "\033[0m"
)

func PrintYellow(text string) string {
	return PrintColor(text, ColorTermYellow)
}

func PrintRed(text string) string {
	return PrintColor(text, ColorTermRed)
}

func PrintColor(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorTermReset)
}
