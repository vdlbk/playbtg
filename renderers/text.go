package renderers

import (
	"fmt"
	"io"
	"strings"

	"github.com/vdlbk/playbtg/utils"
)

type TextRenderer struct {
	Writer io.Writer
}

func (t *TextRenderer) RenderTitle(level int, text string) {
	if level == 0 {
		level = 1
	} else if level > 6 {
		level = 6
	}

	title := strings.Repeat("#", level) + " " + text
	t.RenderText(title)
}

func (t *TextRenderer) RenderTable(data [][]string, headers []string) {
	utils.PrintTable(t.Writer, data, headers, false)
}

func (t *TextRenderer) RenderEmptyLines(n int) {
	utils.PrintEmptyLines(n, t.Writer)
}

func (t *TextRenderer) RenderText(text string) {
	fmt.Fprintln(t.Writer, text)
}

func (t *TextRenderer) RenderExitMessage() {
	// Do nothing
}
