package utils

import (
	"bytes"
	"testing"
)

func TestComputeBounds(t *testing.T) {
	testCases := []struct {
		input       map[int][]string
		expectedMin int
		expectedMax int
	}{
		{nil, 0, 0},
		{map[int][]string{
			3: []string{"B", "D", "F"},
			6: []string{"A"},
		}, 2, 6},
		{map[int][]string{
			3: []string{"B", "D", "F"},
			6: []string{"A"},
			9: []string{"C", "E"},
		}, 3, 9},
		{map[int][]string{
			42: []string{"B", "D", "F"},
			23: []string{"A"},
			16: []string{"C", "E"},
			40: []string{"C", "E"},
			1:  []string{"C", "E"},
			29: []string{"C", "E"},
		}, 1, 42},
	}

	for i, testCase := range testCases {
		t.Run("ComputeBounds", func(t *testing.T) {
			min, max := ComputeBounds(testCase.input)

			if min != testCase.expectedMin || max != testCase.expectedMax {
				t.Errorf("[%d] Output was incorrect, ComputeBounds()= %v, %v; want: %v %v.", i,
					min,
					max,
					testCase.expectedMin,
					testCase.expectedMax,
				)
			}
		})
	}
}

func TestPrintEmptyLines(t *testing.T) {
	var buffer bytes.Buffer

	PrintEmptyLines(3, &buffer)

	result := buffer.String()
	if result != "\n\n\n" {
		t.Errorf("Output was incorrect, PrintEmptyLines()= |%d| want = %v", len(result), 3)
	}
}
