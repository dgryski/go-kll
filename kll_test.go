package kll

import (
	"math/rand"
	"sort"
	"testing"
	"time"
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

func TestCompactorInsertionSort(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, dup := range []bool{false, true} {
		for _, l := range []int{0, 1, 2, 3, 5, 8, 1 << 5, 1 << 10} {
			for i := 0; i < 100; i++ {
				c := make(compactor, l)
				for i := range c {
					if dup && i%2 == 1 {
						c[i] = c[i-1]
					} else {
						c[i] = rng.NormFloat64()
					}
				}
				cp := make(compactor, l)
				copy(cp, c)
				c.insertionSort()
				if !sort.Float64sAreSorted([]float64(c)) {
					t.Fatalf("failed to sort: %v", c)
				}
				sort.Float64s(cp)
				for i, v := range c {
					if v != cp[i] {
						t.Fatalf("failed to sort: %f!=%f @%d\nexpected: %v, got %v", cp[i], v, i, cp, c)
					}
				}
			}
		}
	}
}
