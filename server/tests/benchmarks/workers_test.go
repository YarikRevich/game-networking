package benchmarks_test

import (
	"testing"
	// "github.com/franela/goblin"
)

func BenchmarkWorkers(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.Log(10000 * 20000)
	}
	b.StopTimer()
}
