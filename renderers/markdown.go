package renderers

import (
	"fmt"
	"io"
	"strings"

	"github.com/vdlbk/playbtg/utils"
)

type MarkdownRenderer struct {
	Writer   io.Writer
	FileName *string
}

func (m *MarkdownRenderer) RenderTitle(level int, text string) {
	if level == 0 {
		level = 1
	} else if level > 6 {
		level = 6
	}

	title := strings.Repeat("#", level) + " " + text
	m.RenderText(title)
}

func (m *MarkdownRenderer) RenderTable(data [][]string, headers []string) {
	utils.PrintTable(m.Writer, data, headers, true)
}

func (m *MarkdownRenderer) RenderEmptyLines(n int) {
	utils.PrintEmptyLines(n, m.Writer)
}

func (m *MarkdownRenderer) RenderText(text string) {
	fmt.Fprintln(m.Writer, text)
}

func (m *MarkdownRenderer) RenderExitMessage() {
	if m.FileName != nil {
		fmt.Println("Results were successfully writen to ", *m.FileName)
	}
}
