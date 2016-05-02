package kll

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"sort"

	"github.com/dgryski/go-bitstream"
)

type serializable struct {
	Compactors compactors
	// Compactors []compactor // -- would disable custom serialzation if uncommented.
	K       int
	H       int
	Size    int
	MaxSize int
	Count   int
}

type flat struct {
	v []float64
	h []uint16
}

func (q flat) Len() int { return len(q.v) }

func (q flat) Less(i int, j int) bool { return q.v[i] < q.v[j] }

func (q flat) Swap(i int, j int) {
	q.v[i], q.v[j] = q.v[j], q.v[i]
	q.h[i], q.h[j] = q.h[j], q.h[i]
}

func (c compactors) toFlat() flat {
	var f flat
	for h, c := range c {
		for _, v := range c {
			f.v = append(f.v, v)
			f.h = append(f.h, uint16(h))
		}
	}
	sort.Sort(f)
	return f
}

func (f flat) fromFlat() compactors {
	var cs compactors
	for i, uh := range f.h {
		h := int(uh)
		for len(cs) <= int(uh) {
			cs = append(cs, compactor{})
		}
		cs[h] = append(cs[h], f.v[i])
	}
	return cs
}

func (cs compactors) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	f := cs.toFlat()
	binary.Write(buf, binary.LittleEndian, uint16(len(f.h)))
	w := bitstream.NewWriter(buf)
	encodeSmallInts(w, f.h)
	encodeFloats(w, f.v)
	w.Flush(bitstream.Zero)
	return buf.Bytes(), nil
}

func (cs *compactors) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	var ncs uint16
	err := binary.Read(buf, binary.LittleEndian, &ncs)
	if err != nil {
		return err
	}
	n := int(ncs)
	f := flat{
		h: make([]uint16, n),
		v: make([]float64, n),
	}
	r := bitstream.NewReader(buf)
	err = decodeSmallInts(r, f.h)
	if err != nil {
		return err
	}
	err = decodeFloats(r, f.v)
	if err != nil {
		return err
	}
	*cs = f.fromFlat()
	return nil
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
