package structs

import (
	"fmt"
	"time"
)

type Event struct {
	Word     string
	Duration time.Duration
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.Duration.String(), e.Word)
}