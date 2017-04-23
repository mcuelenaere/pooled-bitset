package pooled_bitset

import (
	"reflect"
	"testing"
)

func TestIterate(t *testing.T) {
	pool := NewFixedCapacityPool(1000)

	bs := pool.Get()
	bs.Set(0)
	bs.Set(100)
	bs.Set(50)
	bs.Set(1)

	data := make([]uint, 4)
	c := 0
	it := bs.Bits()
	for it.Next() {
		if c >= cap(data) {
			t.Fatalf("Iterator has more than %d entries (so far: %v)", cap(data), data)
		}

		data[c] = it.Bit()
		c++
	}

	expected := []uint{0, 1, 50, 100}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v was not expected value %v", data, expected)
	}

	bs.Set(10)
	bs.Set(200)
	data = make([]uint, 6)
	c = 0
	it = bs.Bits()
	for it.Next() {
		data[c] = it.Bit()
		c++
	}

	expected = []uint{0, 1, 10, 50, 100, 200}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v was not expected value %v", data, expected)
	}
}

func BenchmarkIterate(b *testing.B) {
	b.StopTimer()
	pool := NewFixedCapacityPool(10000)
	s := pool.Get()
	for i := 0; i < 10000; i += 3 {
		s.Set(uint(i))
	}
	b.StartTimer()

	for j := 0; j < b.N; j++ {
		it := s.Bits()
		for it.Next() {
			_ = it.Bit()
		}
	}
}

func BenchmarkSparseIterate(b *testing.B) {
	b.StopTimer()
	pool := NewFixedCapacityPool(100000)
	s := pool.Get()
	for i := 0; i < 100000; i += 30 {
		s.Set(uint(i))
	}
	b.StartTimer()

	for j := 0; j < b.N; j++ {
		it := s.Bits()
		for it.Next() {
			_ = it.Bit()
		}
	}
}
