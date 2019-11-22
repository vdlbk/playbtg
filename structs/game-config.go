package structs

import "encoding/json"

type GameConfig struct {
	UpperMode         bool `json:"upper-mode"`
	MixUpperLowerMode bool `json:"upper-lower-mode"`
}

func (g GameConfig) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(b)
}
