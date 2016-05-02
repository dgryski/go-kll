package kll

import (
	"bytes"
	"math"
	"testing"

	"github.com/dgryski/go-bitstream"
)

func TestEncodeSmallInts(t *testing.T) {
	tcs := [][]uint16{
		nil,
		[]uint16{0},
		[]uint16{0, 0},
		[]uint16{0, 0, 0, 0},
		[]uint16{1},
		[]uint16{0, 1},
		[]uint16{1, 0},
		[]uint16{0, 1, 1, 2, 2, 2},
		[]uint16{1000, 100, 1},
		[]uint16{1<<15 - 1, 0},
		[]uint16{0, 1<<15 - 1, 0},
	}
	for _, tc := range tcs {
		b := &bytes.Buffer{}
		w := bitstream.NewWriter(b)
		encodeSmallInts(w, tc)
		w.Flush(bitstream.One)
		got := make([]uint16, len(tc))
		r := bitstream.NewReader(b)
		err := decodeSmallInts(r, got)
		if err != nil {
			t.Fatal(err)
		}
		eq := len(tc) == len(got)
		if eq {
			for i, x := range tc {
				if got[i] != x {
					eq = false
				}
			}
		}
		if !eq {
			t.Fatalf("expected %v, got %v", tc, got)
		}
	}
}
func TestEncodeFloats(t *testing.T) {
	tcs := [][]float64{
		nil,
		[]float64{0},
		[]float64{0, 0},
		[]float64{1},
		[]float64{0, 1},
		[]float64{1, 0},
		[]float64{0, 1, 1, 2, 2, 2},
		[]float64{-1e6, 1, 1e6, math.Inf(-1), math.Inf(1)},
		[]float64{math.NaN()},
	}
	for _, tc := range tcs {
		b := &bytes.Buffer{}
		w := bitstream.NewWriter(b)
		encodeFloats(w, tc)
		w.Flush(bitstream.One)
		got := make([]float64, len(tc))
		r := bitstream.NewReader(b)
		err := decodeFloats(r, got)
		if err != nil {
			t.Fatal(err)
		}
		eq := len(tc) == len(got)
		if eq {
			for i, x := range tc {
				if got[i] != x {
					if math.IsNaN(got[i]) && math.IsNaN(x) {
						continue
					}
					eq = false
				}
			}
		}
		if !eq {
			t.Fatalf("expected %v, got %v", tc, got)
		}
	}
}
