package structs

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	Word     string
	Duration time.Duration
	Attempts Attempts
	Deletion int
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.Duration.String(), e.Word)
}

type Attempts []string

func (a Attempts) String() string {
	return strings.Join(a, " ")
}

func (e *Event) AnalyzeAttempts() []string {
	var result []string

	for _, a := range e.Attempts {
		fmt.Println(a)
		// TODO: Compute word distance
	}

	return result
}
