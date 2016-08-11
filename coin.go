package kll

import "math/rand"

// 64-bit xorshift multiply rng from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf
func xorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}

// coin is a simple struct to let us get random bools and make minimum calls
// to the random number generator.
type coin struct {
	st   uint64
	val  uint64
	bits int
}

func newCoin() coin {
	return coin{
		st: uint64(rand.Int63()),
	}
}

// v is either 0 or 1
func (c *coin) toss() (v int) {
	if c.bits == 0 {
		c.st = xorshiftMult64(c.st)
		c.val = c.st
		c.bits = 64
	}
	c.bits--
	v = int(c.val & 1)
	c.val >>= 1
	return v
}
