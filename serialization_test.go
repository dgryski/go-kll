package kll

import (
	"math/rand"
	"testing"
)

var stateBlackhole State

func BenchmarkGetState(b *testing.B) {
	const k = 1000
	r := New(k)
	for i := 0; i < 100*k; i++ {
		r.Update(rand.NormFloat64())
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stateBlackhole = r.State()
	}
}
