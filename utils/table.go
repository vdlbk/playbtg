package utils

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func PrintTable(writer io.Writer, data [][]string, headers []string, markdown bool) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	if markdown {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	}
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
}
