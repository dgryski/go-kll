package kll

import (
	"bytes"
	"encoding/gob"
)

type serializable struct {
	Compactors []compactor
	K          int
	H          int
	Size       int
	MaxSize    int
	Count      int
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *Sketch) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(
		serializable{
			Compactors: s.compactors,
			K:          s.k,
			H:          s.H,
			Size:       s.size,
			MaxSize:    s.maxSize,
		})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Sketch) UnmarshalBinary(data []byte) error {
	var r serializable
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&r)
	if err != nil {
		return err
	}
	s.compactors = r.Compactors
	s.k = r.K
	s.H = r.H
	s.size = r.Size
	s.maxSize = r.MaxSize
	return nil
}
