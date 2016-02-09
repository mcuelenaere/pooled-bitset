package pooled_bitset

import (
	"sync"
)

type BitSetPool struct {
	length uint
	pool   sync.Pool
}

func NewFixedLengthPool(length uint) *BitSetPool {
	p := &BitSetPool {
		length: length,
	}

	p.pool = sync.Pool {
		New: func() interface{} {
			return &BitSet {
				pool: p,
				length: length,
				set: make([]uint64, wordsNeeded(length)),
			}
		},
	}

	return p
}

func (p *BitSetPool) Get() *BitSet {
	return p.pool.Get().(*BitSet);
}

func (p *BitSetPool) Put(bitSet *BitSet) {
	p.pool.Put(bitSet)
}

func (p *BitSetPool) BitSetLen() uint {
	return p.length
}