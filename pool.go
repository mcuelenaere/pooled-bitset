package pooled_bitset

import (
	"sync"
)

type BitSetPool struct {
	capacity uint
	pool     sync.Pool
}

func NewFixedCapacityPool(capacity uint) *BitSetPool {
	p := &BitSetPool {
		capacity: capacity,
	}

	p.pool = sync.Pool {
		New: func() interface{} {
			return &BitSet {
				pool: p,
				set: make([]uint64, wordsNeeded(capacity)),
			}
		},
	}

	return p
}

// Get returns a BitSet from the pool (which could or could not be a reused instance)
func (p *BitSetPool) Get() *BitSet {
	return p.pool.Get().(*BitSet);
}

// Put gives back the given BitSet to the pool
func (p *BitSetPool) Put(bitSet *BitSet) {
	p.pool.Put(bitSet)
}

// BitSetCapacity returns the capacity of BitSets returned by this pool
func (p *BitSetPool) BitSetCapacity() uint {
	return p.capacity
}