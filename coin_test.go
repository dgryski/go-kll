package kll

import "testing"

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
