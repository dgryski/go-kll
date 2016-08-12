package kll

import (
	"math"
	"testing"
)

func TestComputeHeight(t *testing.T) {
	for i := range heightsCache {
		computed := math.Pow((2.0 / 3.0), float64(i))
		if heightsCache[i] != computed {
			t.Fatalf("cache bad: %v: %v != %v", i, heightsCache[i], computed)
		}
	}
}
