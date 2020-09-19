package utils

import (
	"testing"
)

func TestPrintColor(t *testing.T) {
	testCases := []struct {
		input          string
		color          string
		expectedOutput string
	}{
		{"foobar", ColorTermRed, ColorTermRed + "foobar" + ColorTermReset},
		{"foobar", "color", "colorfoobar" + ColorTermReset},
		{"", "", ColorTermReset},
	}

	for i, testCase := range testCases {
		t.Run("PrintColor", func(t *testing.T) {
			result := PrintColor(testCase.input, testCase.color)

			if result != testCase.expectedOutput {
				t.Errorf("[%d] Output was incorrect, PrintColor(%v, %v)= %v, want: %v.", i,
					testCase.input,
					testCase.color,
					result,
					testCase.expectedOutput,
				)
			}
		})
	}
}
