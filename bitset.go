package pooled_bitset

// the wordSize of a bit set
const wordSize = uint(64)

// log2WordSize is lg(wordSize)
const log2WordSize = uint(6)

type BitSet struct {
	pool     *BitSetPool
	set      []uint64
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

// Cap returns the capacity of the BitSet in bits
func (b *BitSet) Cap() uint {
	return b.pool.capacity
}

// Len returns the number of set bits
func (b *BitSet) Len() uint {
	return uint(popcountSlice(b.set))
}

// And returns a new BitSet containing all bits AND'ed with the given BitSet
func (b *BitSet) And(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	andSlice(result.set, b.set, other.set)
	return result
}

// Or returns a new BitSet containing all bits OR'ed with the given BitSet
func (b *BitSet) Or(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	orSlice(result.set, b.set, other.set)
	return result
}

// Xor returns a new BitSet containing all bits XOR'ed with the given BitSet
func (b *BitSet) Xor(other *BitSet) *BitSet {
	// TODO: verify other is from this pool
	result := other.Clone()
	xorSlice(result.set, b.set, other.set)
	return result
}

// Is the length an exact multiple of word sizes?
func (b *BitSet) isWordAligned() bool {
	return b.Cap() % wordSize == 0
}

// Clean last word by setting unused bits to 0
func (b *BitSet) cleanLastWord() {
	if !b.isWordAligned() {
		// Mask for cleaning last word
		const allBits uint64 = 0xffffffffffffffff
		b.set[wordsNeeded(b.Cap()) - 1] &= allBits >> (wordSize - b.Cap() % wordSize)
	}
}

// Not returns a new BitSet containing all bits NOT'ed with themselves
func (b *BitSet) Not() *BitSet {
	// TODO: verify other is from this pool
	result := b.pool.Get()
	notSlice(result.set, b.set)
	result.cleanLastWord()
	return result
}

// IsEqual returns true whether the given BitSet is equals to ourself
func (b *BitSet) IsEqual(other *BitSet) bool {
	// TODO: verify other is from this pool
	for i, word := range b.set {
		otherWord := other.set[i]
		if word != otherWord {
			return false
		}
	}
	return true
}

// Contains returns true when the given bit is set
func (b *BitSet) Contains(i uint) bool {
	var mask uint64 = 1 << (i & (wordSize-1))
	return (b.set[i >> log2WordSize] & mask) != 0
}

// Flip inverts the bit at the given index
func (b *BitSet) Flip(i uint) {
	b.set[i >> log2WordSize] ^= 1 << (i & (wordSize - 1))
}

// Set sets the bit at the given index to 1
func (b *BitSet) Set(i uint) {
	b.set[i >> log2WordSize] |= 1 << (i & (wordSize - 1))
}

// Clear sets the bit at the given index to 0
func (b *BitSet) Clear(i uint) {
	b.set[i >> log2WordSize] &^= 1 << (i & (wordSize - 1))
}

// ClearAll clears the entire BitSet
func (b *BitSet) ClearAll() {
	for i := range b.set {
		b.set[i] = 0
	}
}

// Clone creates a copy of the BitSet
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
	// TODO: convert this to a callback interface
	return Iterator {
		first: true,
		bitSet: b,
	}
}