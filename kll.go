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
	compactors compactors
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
				s.compactors[h+1] = s.compactors[h].compact(s.compactors[h+1])
				s.updateSize()
				if s.size < s.maxSize {
					break
				}
			}
		}
	}
}

func (s *Sketch) updateSize() {
	s.size = 0
	for _, c := range s.compactors {
		s.size += len(c)
	}
}

// Merge merges a second sketch into this one
func (s *Sketch) Merge(t *Sketch) {
	for s.H < t.H {
		s.grow()
	}

	for h, c := range t.compactors {
		s.compactors[h] = append(s.compactors[h], c...)
	}

	s.updateSize()
	s.compact()
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

func (s *Sketch) Count() int {
	var n int
	for h, c := range s.compactors {
		n += len(c) * (1 << uint(h))
	}
	return n
}

// Quantile estimates the quantile of the value x in the stream.
func (s *Sketch) Quantile(x float64) float64 {
	var r, n int
	for h, c := range s.compactors {
		for _, v := range c {
			w := 1 << uint(h)
			if v <= x {
				r += w
			}
			n += w
		}
	}
	return float64(r) / float64(n)
}

type CDF []Quantile

func (q CDF) Len() int { return len(q) }

func (q CDF) Less(i int, j int) bool { return q[i].V < q[j].V }

func (q CDF) Swap(i int, j int) { q[i], q[j] = q[j], q[i] }

type Quantile struct {
	Q float64
	V float64
}

func (s *Sketch) CDF() CDF {
	q := make(CDF, 0, s.size)

	var totalW float64
	for h, c := range s.compactors {
		weight := float64(int(1 << uint(h)))
		for _, v := range c {
			q = append(q, Quantile{Q: weight, V: v})
		}
		totalW += float64(len(c)) * weight
	}

	sort.Sort(q)

	var curW float64
	for i := range q {
		curW += q[i].Q
		q[i].Q = curW / totalW
	}

	return q
}

// Quantile estimates the quantile of the value x in the stream.
func (q CDF) Quantile(x float64) float64 {
	idx := sort.Search(len(q), func(i int) bool { return q[i].V >= x })
	if idx == 0 {
		return 0
	}
	return q[idx-1].Q
}

// Query estimates the value given quantile p.
func (q CDF) Query(p float64) float64 {
	idx := sort.Search(len(q), func(i int) bool { return q[i].Q >= p })
	if idx == len(q) {
		return q[len(q)-1].V
	}
	return q[idx].V
}

// QuantileLI estimates the quantile of the value x in the stream using linear interpolation.
func (q CDF) QuantileLI(x float64) float64 {
	idx := sort.Search(len(q), func(i int) bool { return q[i].V >= x })
	if idx == len(q) {
		return 1
	}
	if idx == 0 {
		return 0
	}
	// a < x <= b
	a, aq := q[idx-1].V, q[idx-1].Q
	b, bq := q[idx].V, q[idx].Q
	return ((a-x)*bq + (x-b)*aq) / (a - b)
}

// QueryLI estimates the value given quantile p using linear interpolation.
func (q CDF) QueryLI(p float64) float64 {
	idx := sort.Search(len(q), func(i int) bool { return q[i].Q >= p })
	if idx == len(q) {
		return q[len(q)-1].V
	}
	if idx == 0 {
		return q[0].V
	}
	// aq < p <= b
	a, aq := q[idx-1].V, q[idx-1].Q
	b, bq := q[idx].V, q[idx].Q
	return ((aq-p)*b + (p-bq)*a) / (aq - bq)
}

type compactor []float64

func (c *compactor) compact(dst []float64) []float64 {
	sort.Float64s([]float64(*c))
	free := cap(dst) - len(dst)
	if free < len(*c)/2 {
		extra := len(*c)/2 - free
		newdst := make([]float64, len(dst), cap(dst)+extra)
		copy(newdst, dst)
		dst = newdst
	}

	// choose either the evens or the odds
	offs := rand.Intn(2)
	for len(*c) >= 2 {
		l := len(*c) - 2
		dst = append(dst, (*c)[l+offs])
		*c = (*c)[:l]
	}

	return dst
}

type compactors []compactor
