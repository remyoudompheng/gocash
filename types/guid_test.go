package types

import (
	"testing"
)

func TestGuid(t *testing.T) {
	m := make(map[GUID]bool)
	for i := 0; i < 10000; i++ {
		id := NewGUID()
		if m[id] {
			t.Fatalf("GUID collision: %s", id)
		}
		if len(id) != 32 {
			t.Fatalf("GUID has wrong length: %q", id)
		}
		m[id] = true
	}
}

func BenchmarkGuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewGUID()
	}
}
