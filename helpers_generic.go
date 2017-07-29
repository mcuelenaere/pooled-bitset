// +build !amd64

package pooled_bitset

func init() {
	popcountSlice = popcountSliceGeneric
	findFirstSetBit = findFirstSetBitGeneric
	andSlice = andSliceGeneric
	orSlice = orSliceGeneric
	xorSlice = xorSliceGeneric
	notSlice = notSliceGeneric

	popcountSliceVersions = []func([]uint64) uint64{popcountSliceGeneric}
	findFirstSetBitVersions = []func(uint64) uint64{findFirstSetBitGeneric}
	andSliceVersions = []func(dest, a, b []uint64){andSliceGeneric}
	orSliceVersions = []func(dest, a, b []uint64){orSliceGeneric}
	xorSliceVersions = []func(dest, a, b []uint64){xorSliceGeneric}
	notSliceVersions = []func(dest, src []uint64){notSliceGeneric}
}
