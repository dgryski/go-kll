// Package kll implements the KLL streaming quantiles sketch
/*
   http://arxiv.org/pdf/1603.05346v1.pdf
*/
package kll

import (
	"math"
	"math/rand"
	"sort"
)

// Sketch is a streaming quantiles sketch
type Sketch struct {
	compactors []compactor
	k          int
	H          int
	size       int
	maxSize    int
}

// New returns a new Sketch.  k controls the maximum memory used by the stream, which is 3*k + lg(n).
func New(k int) *Sketch {
	s := Sketch{
		k: k,
	}
	s.grow()
	return &s
}

func (s *Sketch) grow() {
	s.compactors = append(s.compactors, compactor{})
	s.H = len(s.compactors)

	s.maxSize = 0
	for h := 0; h < s.H; h++ {
		s.maxSize += s.capacity(h)
	}
}

func (s *Sketch) capacity(h int) int {
	height := float64(s.H - h - 1)
	return int(math.Ceil(float64(s.k)*math.Pow((2.0/3.0), height))) + 1
}

// Update adds x to the stream.
func (s *Sketch) Update(x float64) {
	s.compactors[0] = append(s.compactors[0], x)
	s.size++
	s.compact()
}

func (s *Sketch) compact() {
	for s.size >= s.maxSize {
		for h := 0; h < len(s.compactors); h++ {
			if len(s.compactors[h]) >= s.capacity(h) {
				if h+1 >= s.H {
					s.grow()
				}
				compacted := s.compactors[h].compact()
				s.compactors[h+1] = append(s.compactors[h+1], compacted...)
				s.size = 0
				for _, c := range s.compactors {
					s.size += len(c)
				}
				if s.size < s.maxSize {
					break
				}
			}
		}
	}
}

// Rank estimates the rank of the value x in the stream.
func (s *Sketch) Rank(x float64) int {
	var r int
	for h, c := range s.compactors {
		for _, v := range c {
			if v <= x {
				r += 1 << uint(h)
			}
		}
	}
	return r
}

type compactor []float64

func (c *compactor) compact() []float64 {
	sort.Float64s([]float64(*c))
	dst := make([]float64, 0, len(*c)/2)
	// choose either the evens or the odds
	offs := rand.Intn(2)
	for len(*c) >= 2 {
		l := len(*c) - 2
		dst = append(dst, (*c)[l+offs])
		*c = (*c)[:l]
	}

	return dst
}
