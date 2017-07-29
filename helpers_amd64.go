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
		popcountSliceVersions = []func([]uint64) uint64{popcountSliceGeneric, popcountSliceAsm}
	} else {
		popcountSlice = popcountSliceGeneric
		popcountSliceVersions = []func([]uint64) uint64{popcountSliceGeneric}
	}

	findFirstSetBit = findFirstSetBitAsm
	findFirstSetBitVersions = []func(uint64) uint64{findFirstSetBitGeneric, findFirstSetBitAsm}
	notSlice = notSliceSse2
	notSliceVersions = []func(dest, src []uint64){notSliceGeneric, notSliceSse2}

	if hasAvx() {
		andSlice = andSliceAvx
		orSlice = orSliceAvx
		xorSlice = xorSliceAvx

		andSliceVersions = []func(dest, a, b []uint64){andSliceGeneric, andSliceSse2, andSliceAvx}
		orSliceVersions = []func(dest, a, b []uint64){orSliceGeneric, orSliceSse2, orSliceAvx}
		xorSliceVersions = []func(dest, a, b []uint64){xorSliceGeneric, xorSliceSse2, xorSliceAvx}
	} else {
		andSlice = andSliceSse2
		orSlice = orSliceSse2
		xorSlice = xorSliceSse2

		andSliceVersions = []func(dest, a, b []uint64){andSliceGeneric, andSliceSse2}
		orSliceVersions = []func(dest, a, b []uint64){orSliceGeneric, orSliceSse2}
		xorSliceVersions = []func(dest, a, b []uint64){xorSliceGeneric, xorSliceSse2}
	}
}
