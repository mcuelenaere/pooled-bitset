package pooled_bitset

import (
	"testing"
)

func TestPoolGet(t *testing.T) {
	bs := NewFixedLengthPool(64).Get()
	if bs == nil {
		t.Fatal("Got a nil BitSet from pool")
	}
	if bs.Len() != 64 {
		t.Fatal("Got a BitSet of different length from pool")
	}
}

func TestPoolPut(t *testing.T) {
	pool := NewFixedLengthPool(64)
	bs1 := pool.Get()
	pool.Put(bs1)
	bs2 := pool.Get()
	if bs1 != bs2 {
		t.Fatal("Got a different BitSet instance from pool when previous one was released")
	}
}