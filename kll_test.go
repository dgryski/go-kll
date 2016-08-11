package kll

import (
	"math/rand"
	"testing"
)

func benchmarkAdd(b *testing.B, cons func() float64, k int) {
	// generate the random data
	values := make([]float64, b.N)
	for i := range values {
		values[i] = cons()
	}
	r := New(k)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Update(values[i])
	}
}

func BenchmarkAddNormal_1(b *testing.B) {
	benchmarkAdd(b, rand.NormFloat64, 1)
}

func BenchmarkAddNormal_5(b *testing.B) {
	benchmarkAdd(b, rand.NormFloat64, 5)
}

func BenchmarkAddNormal_10(b *testing.B) {
	benchmarkAdd(b, rand.NormFloat64, 10)
}

func BenchmarkAddNormal_100(b *testing.B) {
	benchmarkAdd(b, rand.NormFloat64, 100)
}

func BenchmarkAddNormal_1000(b *testing.B) {
	benchmarkAdd(b, rand.NormFloat64, 1000)
}
