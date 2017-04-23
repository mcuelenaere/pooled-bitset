package pooled_bitset

//go:noescape
func hasPopCount() bool

//go:noescape
func hasAvx() bool

//go:noescape
func popcountSliceAsm(s []uint64) uint64

//go:noescape
func findFirstSetBitAsm(v uint64) uint64

//go:noescape
func andSliceAvx(dest, a, b []uint64)

//go:noescape
func orSliceAvx(dest, a, b []uint64)

//go:noescape
func xorSliceAvx(dest, a, b []uint64)

//go:noescape
func andSliceSse2(dest, a, b []uint64)

//go:noescape
func orSliceSse2(dest, a, b []uint64)

//go:noescape
func xorSliceSse2(dest, a, b []uint64)

//go:noescape
func notSliceSse2(dest, src []uint64)

func init() {
	if hasPopCount() {
		popcountSlice = popcountSliceAsm
	} else {
		popcountSlice = popcountSliceGeneric
	}

	findFirstSetBit = findFirstSetBitAsm
	notSlice = notSliceSse2

	if hasAvx() {
		andSlice = andSliceAvx
		orSlice = orSliceAvx
		xorSlice = xorSliceAvx
	} else {
		andSlice = andSliceSse2
		orSlice = orSliceSse2
		xorSlice = xorSliceSse2
	}
}
