package structs

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vdlbk/playbtg/utils"

	"github.com/vdlbk/playbtg/utils/consts"
)

type GameConfig struct {
	UpperMode         bool `json:"upper-mode" btg:"Upper mode"`
	MixUpperLowerMode bool `json:"upper-lower-mode" btg:"Mixin Upper/Lower mode"`
	NumberMode        bool `json:"number-mode" btg:"Number mode"`
	InfiniteAttempts  bool `json:"infinite-attempts" btg:"Infinite attempts"`
	WordSetMinLength  int  `json:"-" btg:"-"`
	WordSetMaxLength  int  `json:"-" btg:"-"`
}

func (g GameConfig) Render() {
	config := make([][]string, 0)
	rg := reflect.ValueOf(g)

	for i := 0; i < rg.NumField(); i++ {
		if value := rg.Type().Field(i).Tag.Get(consts.TagKey); value != "-" {
			config = append(config, []string{value, fmt.Sprintf("%v", rg.Field(i))})
		}
	}

	utils.PrintTable(config, []string{"Mode", "value"})
}

func (g GameConfig) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(b)
}
