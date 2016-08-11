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
	mask uint64
}

// v is either 0 or 1
func (c *coin) toss() (v int) {
	if c.mask == 0 {
		if c.st == 0 {
			c.st = uint64(rand.Int63())
		}
		c.st = xorshiftMult64(c.st)
		c.mask = 1
	}
	if c.st&c.mask > 0 {
		v = 1
	}
	c.mask <<= 1
	return v
}
