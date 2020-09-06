package utils

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func ComputeBounds(set map[int][]string) (int, int) {
	min := len(set)
	max := 0
	for k, _ := range set {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	return min, max
}

func FiftyFifty() bool {
	return rand.Float64() >= 0.5
}

func PrintEmptyLines(n int, writer io.Writer) {
	for i := 0; i < n; i++ {
		if _, err := fmt.Fprintln(writer); err != nil {
			fmt.Println(err)
			break
		}
	}
}
