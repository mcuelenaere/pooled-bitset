// +build !amd64

package pooled_bitset

func init() {
	popcountSlice = popcountSliceGeneric
	findFirstSetBit = findFirstSetBitGeneric
	andSlice = andSliceGeneric
	orSlice = orSliceGeneric
	xorSlice = xorSliceGeneric
	notSlice = notSliceGeneric
}
