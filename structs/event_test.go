package structs

import (
	"reflect"
	"sort"
	"testing"
)

func TestAnalysisMap_ToCharSwaps(t *testing.T) {
	testCases := []struct {
		input          AnalysisMap
		expectedResult CharSwaps
	}{
		{nil, nil},
		{map[rune]Analysis{
			'a': {
				Chars: map[rune]int{'b': 1, 'c': 2},
			},
			'q': {
				Chars: map[rune]int{'a': 4},
			},
		}, CharSwaps{
			{'a', 'b', 1},
			{'a', 'c', 2},
			{'q', 'a', 4},
		}},
		{map[rune]Analysis{
			'a': {},
			'q': {
				Chars: map[rune]int{'a': 4},
			},
		}, CharSwaps{
			{'q', 'a', 4},
		}},
		{map[rune]Analysis{
			'a': {},
			'q': {
				Chars: map[rune]int{},
			},
		}, nil,
		},
	}

	for i, testCase := range testCases {
		t.Run("ToCharSwaps", func(t *testing.T) {
			result := testCase.input.ToCharSwaps()

			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Errorf("[%d] Output was incorrect, ToCharSwaps()= %v; want: %v.", i,
					result,
					testCase.expectedResult,
				)
			}
		})
	}
}

func TestCharSwaps_Sort(t *testing.T) {
	testCases := []struct {
		input          CharSwaps
		expectedResult CharSwaps
	}{
		{nil, nil},
		{CharSwaps{}, CharSwaps{}},
		{CharSwaps{
			{'a', 'b', 1},
			{'a', 'c', 2},
			{'q', 'a', 4},
		}, CharSwaps{
			{'q', 'a', 4},
			{'a', 'c', 2},
			{'a', 'b', 1},
		}},
		{CharSwaps{
			{'a', 'b', 1},
			{'q', 'a', 4},
			{'a', 'c', 2},
			{'w', 'z', 9},
			{'z', 'w', 5},
		}, CharSwaps{
			{'w', 'z', 9},
			{'z', 'w', 5},
			{'q', 'a', 4},
			{'a', 'c', 2},
			{'a', 'b', 1},
		}},
	}

	for i, testCase := range testCases {
		t.Run("ToCharSwaps", func(t *testing.T) {
			sort.Sort(testCase.input)

			if !reflect.DeepEqual(testCase.input, testCase.expectedResult) {
				t.Errorf("[%d] Output was incorrect, ToCharSwaps()= %v; want: %v.", i,
					testCase.input,
					testCase.expectedResult,
				)
			}
		})
	}
}
