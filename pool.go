package pooled_bitset

import (
	"sync"
)

type BitSetPool struct {
	capacity uint
	pool     sync.Pool
}

// NewFixedCapacityPool creates a pool, returning BitSet's of a fixed capacity
func NewFixedCapacityPool(capacity uint) *BitSetPool {
	p := &BitSetPool{
		capacity: capacity,
	}

	p.pool = sync.Pool{
		New: func() interface{} {
			return &BitSet{
				pool: p,
				set:  make([]uint64, wordsNeeded(capacity)),
			}
		},
	}

	return p
}

func (p *BitSetPool) fastGet() *BitSet {
	return p.pool.Get().(*BitSet)
}

// Get returns an empty BitSet from the pool
func (p *BitSetPool) Get() *BitSet {
	bs := p.fastGet()
	bs.ClearAll()
	return bs
}

// Put gives back the given BitSet to the pool
func (p *BitSetPool) Put(bitSet *BitSet) {
	p.pool.Put(bitSet)
}

// BitSetCapacity returns the capacity of BitSets returned by this pool
func (p *BitSetPool) BitSetCapacity() uint {
	return p.capacity
}
