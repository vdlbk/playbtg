package structs

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Word     string
	Duration time.Duration
	Attempts Attempts
	Deletion int
	Errors   int
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.Duration.String(), e.Word)
}

type Attempts []string

func (a Attempts) String() string {
	tmp := make([]string, len(a))
	copy(tmp, a)
	for idx, text := range a {
		if text == "" || text == " " {
			tmp[idx] = "<SPACE>"
		}
	}
	return strings.Join(tmp, " ")
}

func (e *Event) AnalyzeAttempts(result *AnalysisMap) {
	m := *result

	for _, attempt := range e.Attempts {
		// TODO: Compute word distance

		for idx := range e.Word {
			if idx < len(attempt) && e.Word[idx] != attempt[idx] {
				if analysis, ok := m[rune(e.Word[idx])]; ok {
					analysis.Chars[rune(attempt[idx])]++
				} else {
					chars := make(map[rune]int)
					chars[rune(attempt[idx])]++
					m[rune(e.Word[idx])] = Analysis{Chars: chars}
				}
			}
		}
	}

	result = &m
}

type AnalysisMap map[rune]Analysis

func (a *AnalysisMap) ToCharSwaps() CharSwaps {
	var result CharSwaps
	for origin, analysis := range *a {
		for c, count := range analysis.Chars {
			result = append(result, CharSwap{Origin: origin, Swap: c, Count: count})
		}
	}

	return result
}

type Analysis struct {
	Chars map[rune]int
}

type CharSwap struct {
	Origin rune
	Swap   rune
	Count  int
}

func (c *CharSwap) ToStrings() []string {
	return []string{fmt.Sprintf("%c", c.Origin), fmt.Sprintf("%c", c.Swap), strconv.Itoa(c.Count)}
}

type CharSwaps []CharSwap

func (c CharSwaps) Len() int {
	return len(c)
}

// Reversed by default
func (c CharSwaps) Less(i, j int) bool {
	return c[i].Count > c[j].Count
}

func (c CharSwaps) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
