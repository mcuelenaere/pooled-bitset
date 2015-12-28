package pooled_bitset

// the wordSize of a bit set
const wordSize = uint(64)

// log2WordSize is lg(wordSize)
const log2WordSize = uint(6)

type BitSet struct {
	pool   *BitSetPool
	length uint
	set    []uint64
}

// wordsNeeded calculates the number of words needed for i bits
func wordsNeeded(i uint) int {
	if i > ((^uint(0)) - wordSize + 1) {
		return int((^uint(0)) >> log2WordSize)
	}
	return int((i + (wordSize - 1)) >> log2WordSize)
}

// Release gives the BitSet instance back to the pool
func (b *BitSet) Release() {
	b.pool.Put(b)
}

// Len returns the length of the BitSet in words
func (b *BitSet) Len() uint {
	return b.length
}

func (b *BitSet) And(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	andSlice(result.set, b.set, other.set)
	return result
}

func (b *BitSet) Or(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	orSlice(result.set, b.set, other.set)
	return result
}

func (b *BitSet) Xor(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	xorSlice(result.set, b.set, other.set)
	return result
}

func (b *BitSet) Not() *BitSet {
	// TODO: verify other is from this pool
	result := b.pool.Get()
	notSlice(result.set, b.set)
	return result
}

func (b *BitSet) Count() uint {
	return uint(popcountSlice(b.set))
}

func (b *BitSet) Contains(i uint) bool {
	var mask uint64 = 1 << (i & (wordSize-1))
	return (b.set[i >> log2WordSize] & mask) != 0
}

func (b *BitSet) Flip(i uint) {
	b.set[i >> log2WordSize] ^= 1 << (i & (wordSize - 1))
}

func (b *BitSet) Set(i uint) {
	b.set[i >> log2WordSize] |= 1 << (i & (wordSize - 1))
}

func (b *BitSet) Clear(i uint) {
	b.set[i >> log2WordSize] &^= 1 << (i & (wordSize - 1))
}

// ClearAll clears the entire BitSet
func (b *BitSet) ClearAll() {
	for i := range b.set {
		b.set[i] = 0
	}
}

func (b *BitSet) Clone() *BitSet {
	c := b.pool.Get()
	copy(c.set, b.set)
	return c
}

type Iterator struct {
	first      bool
	setIdx     int
	currentBit uint
	bitSet     *BitSet
}

func (i *Iterator) Bit() uint {
	return uint(i.setIdx) * wordSize + i.currentBit
}

func (i *Iterator) Next() bool {
	var currentValue uint64
	if i.first {
		currentValue = i.bitSet.set[i.setIdx]
		i.first = false
	} else {
		currentValue = i.bitSet.set[i.setIdx] &^ (1 << (i.currentBit + 1) - 1)
	}

	for currentValue == 0 {
		i.setIdx++
		if i.setIdx >= len(i.bitSet.set) {
			return false
		}
		currentValue = i.bitSet.set[i.setIdx]
		i.currentBit = 0
	}
	i.currentBit = uint(findFirstSetBit(currentValue))
	return true
}

func (b *BitSet) Bits() Iterator {
	return Iterator {
		first: true,
		bitSet: b,
	}
}