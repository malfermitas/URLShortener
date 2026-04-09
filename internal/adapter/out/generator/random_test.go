package generator

import (
	"testing"
)

func BenchmarkURLGenerator_Generate(b *testing.B) {
	gen := NewURLGenerator()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		gen.Generate()
	}
}

func BenchmarkURLGenerator_Generate_Parallel(b *testing.B) {
	gen := NewURLGenerator()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gen.Generate()
		}
	})
}

func TestSnowflakeGenerator_Uniqueness(t *testing.T) {
	gen := NewURLGenerator()
	seen := make(map[string]bool)

	for i := 0; i < 100000; i++ {
		key := gen.Generate()
		if seen[key] {
			t.Fatalf("Duplicate key found: %s at iteration %d", key, i)
		}
		seen[key] = true
	}
}
