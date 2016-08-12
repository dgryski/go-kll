package kll

import (
	"bytes"
	"encoding/gob"
	"unsafe"
)

// we know that compactor is really a []float64, and we want to refer to them
// in the state, so we can just unsafely convert them.

func compactorsAsFloats(c []compactor) [][]float64 {
	return *(*[][]float64)(unsafe.Pointer(&c))
}

func floatsAsCompactors(f [][]float64) []compactor {
	return *(*[]compactor)(unsafe.Pointer(&f))
}

// State represents the state of the Sketch. It is used for serializing and
// deserializing to disk.
type State struct {
	Compactors [][]float64
	K          int
	H          int
	Size       int
	MaxSize    int
	Count      int
}

// State returns the current state of the Sketch. The state is invalid if any
// other methods of the Sketch are called, and it must not be mutated.
func (s *Sketch) State() State {
	return State{
		Compactors: compactorsAsFloats(s.compactors),
		K:          s.k,
		H:          s.H,
		Size:       s.size,
		MaxSize:    s.maxSize,
	}
}

// SetState sets the state of the Sketch to the passed State. The memory is
// shared, so the passed State is invalid to be read from or written to after
// this call.
func (s *Sketch) SetState(state State) {
	s.compactors = floatsAsCompactors(state.Compactors)
	s.k = state.K
	s.H = state.H
	s.size = state.Size
	s.maxSize = state.MaxSize
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *Sketch) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s.State())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Sketch) UnmarshalBinary(data []byte) error {
	var r State
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&r)
	if err != nil {
		return err
	}
	s.SetState(r)
	return nil
}
