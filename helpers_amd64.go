package pooled_bitset

//go:noescape
func hasPopCount() bool

//go:noescape
func hasAvx2() bool

//go:noescape
func popcountSliceAsm(s []uint64) uint64

//go:noescape
func findFirstSetBitAsm(v uint64) uint64

//go:noescape
func andSliceAvx2(dest, a, b []uint64)

//go:noescape
func orSliceAvx2(dest, a, b []uint64)

//go:noescape
func xorSliceAvx2(dest, a, b []uint64)

//go:noescape
func notSliceAvx2(dest, src []uint64)

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

	if hasAvx2() {
		andSlice = andSliceAvx2
		orSlice = orSliceAvx2
		xorSlice = xorSliceAvx2
		notSlice = notSliceAvx2

		andSliceVersions = []func(dest, a, b []uint64){andSliceGeneric, andSliceSse2, andSliceAvx2}
		orSliceVersions = []func(dest, a, b []uint64){orSliceGeneric, orSliceSse2, orSliceAvx2}
		xorSliceVersions = []func(dest, a, b []uint64){xorSliceGeneric, xorSliceSse2, xorSliceAvx2}
		notSliceVersions = []func(dest, src []uint64){notSliceGeneric, notSliceSse2, notSliceAvx2}
	} else {
		andSlice = andSliceSse2
		orSlice = orSliceSse2
		xorSlice = xorSliceSse2
		notSlice = notSliceSse2

		andSliceVersions = []func(dest, a, b []uint64){andSliceGeneric, andSliceSse2}
		orSliceVersions = []func(dest, a, b []uint64){orSliceGeneric, orSliceSse2}
		xorSliceVersions = []func(dest, a, b []uint64){xorSliceGeneric, xorSliceSse2}
		notSliceVersions = []func(dest, src []uint64){notSliceGeneric, notSliceSse2}
	}
}
