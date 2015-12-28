package pooled_bitset

var (
	popcountSlice func([]uint64) uint64
	findFirstSetBit func(uint64) uint64
	andSlice func(dest, a, b []uint64)
	orSlice func(dest, a, b []uint64)
	xorSlice func(dest, a, b []uint64)
	notSlice func(dest, src []uint64)
)

// bit population count, take from
// https://code.google.com/p/go/issues/detail?id=4988#c11
// credit: https://code.google.com/u/arnehormann/
func popcountSliceGeneric(s []uint64) (n uint64) {
	cnt := uint64(0)
	for _, x := range s {
		x -= (x >> 1) & 0x5555555555555555
		x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
		x += x >> 4
		x &= 0x0f0f0f0f0f0f0f0f
		x *= 0x0101010101010101
		cnt += x >> 56
	}
	return cnt
}

var deBruijn = [...]byte{
	0, 1, 56, 2, 57, 49, 28, 3, 61, 58, 42, 50, 38, 29, 17, 4,
	62, 47, 59, 36, 45, 43, 51, 22, 53, 39, 33, 30, 24, 18, 12, 5,
	63, 55, 48, 27, 60, 41, 37, 16, 46, 35, 44, 21, 52, 32, 23, 11,
	54, 26, 40, 15, 34, 20, 31, 10, 25, 14, 19, 9, 13, 8, 7, 6,
}

func findFirstSetBitGeneric(v uint64) uint64 {
	return uint64(deBruijn[((v&-v)*0x03f79d71b4ca8b09)>>58])
}

func andSliceGeneric(dest, a, b []uint64) {
	for i, word := range a {
		dest[i] = word & b[i]
	}
}

func orSliceGeneric(dest, a, b []uint64) {
	for i, word := range a {
		dest[i] = word | b[i]
	}
}

func xorSliceGeneric(dest, a, b []uint64) {
	for i, word := range a {
		dest[i] = word ^ b[i]
	}
}

func notSliceGeneric(dest, src []uint64) {
	for i, word := range src {
		dest[i] = ^word
	}
}