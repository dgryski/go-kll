package kll

import (
	"math/rand"
	"testing"
	"time"
)

func TestCoin(t *testing.T) {
	// set up a coin that should return alternating bits
	c := coin{
		st:   0xaaaaaaaaaaaaaaaa,
		mask: 1,
	}

	for i := 0; i < 64; i++ {
		if v := c.toss(); v != i&1 {
			t.Fatalf("toss %d: %d != %d", i, v, i&1)
		}
	}
}

func TestCoinMany(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := coin{
		st:   uint64(rng.Int63()),
		mask: 0,
	}
	t.Logf("state: 0x%016x", c.st)

	pos := 0
	for i := 0; i < 1000; i++ {
		v := c.toss()
		if v != 0 && v != 1 {
			t.Fatal("invalid value from coin:", v)
		}
		if v == 1 {
			pos++
		}
	}

	t.Logf("pos: %v", pos)

	// someone can do the binomial/normal, but i expect somewhere between
	// 400 and 600 will never fail.
	if pos < 400 || pos > 600 {
		t.Fatal("abornmal bias in flips:", pos)
	}
}
