package kll

import (
	"math"

	"github.com/dgryski/go-bits"
	"github.com/dgryski/go-bitstream"
)

// encodes ints without size.
// Should work for uint15
func encodeSmallInts(w *bitstream.BitWriter, xs []uint16) {
	if len(xs) == 0 {
		return
	}
	var m uint16
	for _, x := range xs {
		if x > m {
			m = x
		}
	}
	l := 64 - bits.Clz(uint64(m))
	w.WriteBits(uint64(l), 4)
	for _, x := range xs {
		w.WriteBits(uint64(x), int(l))
	}
}

func decodeSmallInts(r *bitstream.BitReader, xs []uint16) error {
	if len(xs) == 0 {
		return nil
	}
	m, err := r.ReadBits(4)
	if err != nil {
		return err
	}
	for i := range xs {
		v, err := r.ReadBits(int(m))
		if err != nil {
			return err
		}
		xs[i] = uint16(v)
	}
	return nil
}

// encodes floats without size.
func encodeFloats(w *bitstream.BitWriter, fs []float64) {
	if len(fs) == 0 {
		return
	}
	prev := math.Float64bits(fs[0])
	w.WriteBits(prev, 64)
	for _, f := range fs[1:] {
		cur := math.Float64bits(f)
		delta := prev ^ cur
		prev = cur
		se := delta >> 52
		if se == 0 {
			w.WriteBit(bitstream.Zero)
		} else {
			w.WriteBit(bitstream.One)
			w.WriteBits(se, 12)
		}
		m := delta & ((1 << 52) - 1)
		nm := 64 - bits.Clz(m)
		if 52-nm > 7 {
			w.WriteBit(bitstream.One)
			w.WriteBits(52-nm, 6)
		} else {
			w.WriteBit(bitstream.Zero)
			w.WriteBits(52-nm, 3)
		}
		w.WriteBits(m, int(nm))
	}
}

func decodeFloats(r *bitstream.BitReader, fs []float64) error {
	if len(fs) == 0 {
		return nil
	}
	prev, err := r.ReadBits(64)
	if err != nil {
		return err
	}
	fs[0] = math.Float64frombits(prev)

	for i := 1; i < len(fs); i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return err
		}
		var se uint64
		if bit == bitstream.One {
			se, err = r.ReadBits(12)
			if err != nil {
				return err
			}
		}
		bit, err = r.ReadBit()
		if err != nil {
			return err
		}
		var nm uint64
		if bit == bitstream.One {
			nm, err = r.ReadBits(6)
		} else {
			nm, err = r.ReadBits(3)
		}
		if err != nil {
			return err
		}
		nm = 52 - nm
		m, err := r.ReadBits(int(nm))
		if err != nil {
			return err
		}
		m <<= 52 - nm
		delta := (se << 52) | m
		cur := prev ^ delta
		prev = cur
		fs[i] = math.Float64frombits(cur)
	}
	return nil
}
