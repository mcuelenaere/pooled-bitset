package pooled_bitset

import (
	"testing"
)

func TestPoolGet(t *testing.T) {
	bs := NewFixedCapacityPool(64).Get()
	if bs == nil {
		t.Fatal("Got a nil BitSet from pool")
	}
	if bs.Cap() != 64 {
		t.Fatal("Got a BitSet of different length from pool")
	}
}

func TestPoolGetReturnsZeroedBitSet(t *testing.T) {
	pool := NewFixedCapacityPool(64)

	// allocate bitset, fill it and return it back to the pool
	bs := pool.Get()
	for i := uint(0); i < 64; i++ {
		bs.Set(i)
	}
	bs.Release()

	// allocate another bitset and check if it's zeroed out
	bs = pool.Get()
	if bs.Len() > 0 {
		t.Fatalf("Bitset is %v, expected it to be 0!", bs.Len())
	}
	bs.Release()
}