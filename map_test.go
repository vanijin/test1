package main

import (
	"testing"
)

func TestIntMinBasic(t *testing.T) {
	ans := IntMix(2, -2)
	if ans != 0 {
		t.Errorf("IntMin(2, -2) = %d; want 0", ans)
	}
}

func BenchmarkIntMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntMix(1, 2)
	}
}
