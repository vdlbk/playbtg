package utils

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func PrintTable(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
}
