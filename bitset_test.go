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
	it := bs.Iterator()
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
	it = bs.Iterator()
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
		c := uint(0)
		it := s.Iterator()
		for it.Next() {
			c += it.Bit()
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
		c := uint(0)
		it := s.Iterator()
		for it.Next() {
			c += it.Bit()
		}
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	f()
}

func TestOperationsOnBitSetsOfDifferentSizes(t *testing.T) {
	pool1 := NewFixedCapacityPool(64)
	pool2 := NewFixedCapacityPool(32)

	bs1 := pool1.Get()
	bs2 := pool2.Get()

	if (bs1.IsEqual(bs2) || bs2.IsEqual(bs1)) {
		t.Error("BitSets aren't supposed to be equal")
	}

	assertPanic(t, func() { bs1.And(bs2) })
	assertPanic(t, func() { bs2.And(bs1) })

	assertPanic(t, func() { bs1.Or(bs2) })
	assertPanic(t, func() { bs2.Or(bs1) })

	assertPanic(t, func() { bs1.Xor(bs2) })
	assertPanic(t, func() { bs2.Xor(bs1) })
}

func TestIterateUsingCallback(t *testing.T) {
	pool := NewFixedCapacityPool(1000)

	bs := pool.Get()
	bs.Set(0)
	bs.Set(100)
	bs.Set(50)
	bs.Set(1)

	data := make([]uint, 4)
	c := 0
	bs.Walk(func(i uint) {
		if c >= cap(data) {
			t.Fatalf("Iterator has more than %d entries (so far: %v)", cap(data), data)
		}

		data[c] = i
		c++
	})

	expected := []uint{0, 1, 50, 100}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v was not expected value %v", data, expected)
	}

	bs.Set(10)
	bs.Set(200)
	data = make([]uint, 6)
	c = 0
	bs.Walk(func(i uint) {
		data[c] = i
		c++
	})

	expected = []uint{0, 1, 10, 50, 100, 200}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v was not expected value %v", data, expected)
	}
}

func BenchmarkIterateUsingCallback(b *testing.B) {
	b.StopTimer()
	pool := NewFixedCapacityPool(10000)
	s := pool.Get()
	for i := 0; i < 10000; i += 3 {
		s.Set(uint(i))
	}
	b.StartTimer()

	for j := 0; j < b.N; j++ {
		c := uint(0)
		s.Walk(func (i uint) {
			c += i
		})
	}
}

func BenchmarkSparseIterateUsingCallback(b *testing.B) {
	b.StopTimer()
	pool := NewFixedCapacityPool(100000)
	s := pool.Get()
	for i := 0; i < 100000; i += 30 {
		s.Set(uint(i))
	}
	b.StartTimer()

	for j := 0; j < b.N; j++ {
		c := uint(0)
		s.Walk(func (i uint) {
			c += i
		})
	}
}
