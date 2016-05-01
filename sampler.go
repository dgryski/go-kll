package kll

import "math/rand"

type sampler struct {
	h uint16
	w uint64
	y float64
}

func (s *sampler) update(x float64, w uint64, to []float64) []float64 {
	ph := uint64(1 << s.h)
	switch {
	case s.w+w <= ph:
		s.w += w
		if rand.Float64()*float64(w) < float64(s.w) {
			s.y = x
		}
		if s.w == ph {
			s.w = 0
			return append(to, s.y)
		}
	case s.w < w:
		if rand.Float64()*float64(w) < float64(ph) {
			return append(to, x)
		}
	default: // W >= w
		s.w = w
		s.y = x
		if rand.Float64()*float64(w) < float64(ph) {
			return append(to, x)
		}
	}
	return to
}

func (s *sampler) grow() {
	s.h++
}
