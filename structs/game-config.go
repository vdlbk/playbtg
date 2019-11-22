package structs

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/vdlbk/playbtg/utils/consts"
	"os"
	"reflect"
)

type GameConfig struct {
	UpperMode         bool `json:"upper-mode" btg:"Upper mode"`
	MixUpperLowerMode bool `json:"upper-lower-mode" btg:"Mixin Upper/Lower mode"`
	NumberMode        bool `json:"number-mode" btg:"Number mode"`
	WordSetMinLength  int  `json:"-" btg:"-"`
	WordSetMaxLength  int  `json:"-" btg:"-"`
}

func (g GameConfig) Render() {
	config := make([][]string, 0)
	rg := reflect.ValueOf(g)

	for i := 0; i < rg.NumField(); i++ {
		if v := rg.Type().Field(i).Tag.Get(consts.TagKey); v != "-" {
			config = append(config, []string{v, fmt.Sprintf("%v", rg.Field(i))})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Configuration key", "value"})
	table.SetCenterSeparator("|")
	table.AppendBulk(config)
	table.Render()
}

func (g GameConfig) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(b)
}
