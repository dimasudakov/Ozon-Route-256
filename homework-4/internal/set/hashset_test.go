package set

import "testing"

func Test_HashSet(t *testing.T) {
	set := NewHashSet[int]()

	N := 20

	for i := 0; i < N; i += 2 {
		set.Add(i)
		set.Add(i)
	}
	if set.Size() != uint32(N/2) {
		t.Errorf("Set length should be %d", N/2)
	}
	for i := 0; i < N; i += 2 {
		if !set.Contains(i) {
			t.Errorf("Expected set to contain %d", i)
		}
		if set.Contains(i + 1) {
			t.Errorf("Set shouldn't contain %d", i+1)
		}
	}

	set.Remove(2)
	set.Remove(10)
	if set.Contains(2) || set.Contains(10) {
		t.Errorf("Expected set to not contain %d or %d after removal", 2, 10)
	}

	set.Add(2)
	set.Add(3)
	set.RemoveAll()
	if set.Size() != 0 {
		t.Errorf("Expected set to be empty after RemoveAll")
	}
}

func Benchmark_HashSetAdd(b *testing.B) {
	set := NewHashSet[int]()
	for i := 0; i < b.N; i++ {
		set.Add(i)
	}
}

func Benchmark_HashSetContains(b *testing.B) {
	N := int(1e6)
	set := NewHashSet[int]()
	for i := 0; i < N; i++ {
		set.Add(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Contains(i)
	}
}

func Benchmark_HashSetRemove(b *testing.B) {
	N := int(1e6)
	set := NewHashSet[int]()
	for i := 0; i < N; i++ {
		set.Add(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Remove(i)
	}
}

func Benchmark_HashSetRemoveAll(b *testing.B) {
	N := int(1e6)
	set := NewHashSet[int]()
	for i := 0; i < N; i++ {
		set.Add(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.RemoveAll()
	}
}
