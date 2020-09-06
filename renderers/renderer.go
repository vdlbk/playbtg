package renderers

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/vdlbk/playbtg/utils/consts"
)

type Renderer interface {
	RenderTitle(int, string)
	RenderTable([][]string, []string)
	RenderEmptyLines(int)
	RenderText(string)
	RenderExitMessage()
}

func InitFile(folderPath string) (io.Writer, *string) {
	now := time.Now()
	now.Format("20060102_150405")
	fileName := fmt.Sprintf("%s_result_%s.md", consts.GlobalAppName, now.Format("20060102_150405"))
	if !strings.HasSuffix(folderPath, "/") {
		fileName = "/" + fileName
	}
	fileName = folderPath + fileName
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(consts.ErrCannotOpenFile)
		return nil, nil
	}
	return f, &fileName
}
