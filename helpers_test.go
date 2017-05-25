package pooled_bitset

import (
	"math/rand"
	"reflect"
	"runtime"
	"testing"
)

var popCountSliceTestCases = []struct {
	Input    []uint64
	Expected uint64
}{
	{[]uint64{}, 0},
	{[]uint64{0x1}, 1},
	{[]uint64{0x2}, 1},
	{[]uint64{0x44}, 2},
	{[]uint64{0x8000000000000000}, 1},
	{[]uint64{0xFFFFFFFFFFFFFFFF}, 64},

	{[]uint64{0x1, 0x2}, 2},
	{[]uint64{0x8000000000000000, 0xFFFFFFFFFFFFFFFF}, 65},
}

func TestPopCountSliceGeneric(t *testing.T) {
	for _, testCase := range popCountSliceTestCases {
		output := popcountSliceGeneric(testCase.Input)
		if output != testCase.Expected {
			t.Errorf("popCountSliceGeneric() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
		}
	}
}

func TestPopCountSlice(t *testing.T) {
	for _, testCase := range popCountSliceTestCases {
		output := popcountSlice(testCase.Input)
		if output != testCase.Expected {
			t.Errorf("popcountSliceAsm() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
		}
	}
}

// this is necessary because otherwise the inliner would unfairly give the non-assembly function an advantage
var indirectPopCountSliceGeneric = popcountSliceGeneric

func generatePopCountSlice() []uint64 {
	slice := make([]uint64, 8192*10)
	for n := 0; n < len(slice); n++ {
		slice[n] = uint64(n)
	}
	return slice
}

func BenchmarkPopCountSliceGeneric(b *testing.B) {
	slice := generatePopCountSlice()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		indirectPopCountSliceGeneric(slice)
	}
}

func BenchmarkPopCountSlice(b *testing.B) {
	slice := generatePopCountSlice()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		popcountSlice(slice)
	}
}

var findFirstSetBitTestCases = []struct {
	Input    uint64
	Expected uint64
}{
	{0x1, 0},
	{0x2, 1},
	{0x44, 2},
	{0x8000000000000000, 63},
	{0xFFFFFFFFFFFFFFFF, 0},
}

func TestFindFirstSetBitGeneric(t *testing.T) {
	for _, testCase := range findFirstSetBitTestCases {
		output := findFirstSetBitGeneric(testCase.Input)
		if output != testCase.Expected {
			t.Errorf("findFirstSetBitGeneric() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
		}
	}
}

func TestFindFirstSetBit(t *testing.T) {
	for _, testCase := range findFirstSetBitTestCases {
		output := findFirstSetBit(testCase.Input)
		if output != testCase.Expected {
			t.Errorf("findFirstSetBitAsm() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
		}
	}
}

// this is necessary because otherwise the inliner would unfairly give the non-assembly function an advantage
var indirectFindFirstSetBitGeneric = findFirstSetBitGeneric

func BenchmarkFindFirstSetBitGeneric(b *testing.B) {
	for n := 0; n < b.N; n++ {
		indirectFindFirstSetBitGeneric(0x12345678)
	}
}

func BenchmarkFindFirstSetBitAsm(b *testing.B) {
	for n := 0; n < b.N; n++ {
		findFirstSetBit(0x12345678)
	}
}

var bitOpsSliceTestCases = []struct {
	InputA      []uint64
	InputB      []uint64
	ExpectedAnd []uint64
	ExpectedOr  []uint64
	ExpectedXor []uint64
}{
	{[]uint64{}, []uint64{}, []uint64{}, []uint64{}, []uint64{}},
	{[]uint64{0x1}, []uint64{0x0}, []uint64{0x0}, []uint64{0x1}, []uint64{0x1}},
	{[]uint64{0x1, 0x1, 0x1, 0x1}, []uint64{0x0, 0x1, 0x0, 0x0}, []uint64{0x0, 0x1, 0x0, 0x0}, []uint64{0x1, 0x1, 0x1, 0x1}, []uint64{0x1, 0x0, 0x1, 0x1}},
}

const (
	OP_AND = iota
	OP_OR
	OP_XOR
)

func testBitOpSlice(t *testing.T, bitOpFunc func(dest, a, b []uint64), bitOp int) {
	for _, testCase := range bitOpsSliceTestCases {
		output := make([]uint64, len(testCase.InputA))
		var expected []uint64
		switch bitOp {
		case OP_AND:
			expected = testCase.ExpectedAnd
		case OP_OR:
			expected = testCase.ExpectedOr
		case OP_XOR:
			expected = testCase.ExpectedXor
		}

		bitOpFunc(output, testCase.InputA, testCase.InputB)
		if !reflect.DeepEqual(output, expected) {
			t.Errorf("%s gave %v instead of expected %v for inputs %v and %v", runtime.FuncForPC(reflect.ValueOf(bitOpFunc).Pointer()).Name(), output, expected, testCase.InputA, testCase.InputB)
		}
	}
}

func TestAndSliceGeneric(t *testing.T) { testBitOpSlice(t, andSliceGeneric, OP_AND) }
func TestAndSliceAsm(t *testing.T)     { testBitOpSlice(t, andSlice, OP_AND) }

func TestOrdSliceGeneric(t *testing.T) { testBitOpSlice(t, orSliceGeneric, OP_OR) }
func TestOrdSliceAsm(t *testing.T)     { testBitOpSlice(t, orSlice, OP_OR) }

func TestXorSliceGeneric(t *testing.T) { testBitOpSlice(t, xorSliceGeneric, OP_XOR) }
func TestXorSliceAsm(t *testing.T)     { testBitOpSlice(t, xorSlice, OP_XOR) }

func benchmarkBitOp(b *testing.B, sliceLength int, bitOpFunc func(dest, a, b []uint64)) {
	b.StopTimer()

	output := make([]uint64, sliceLength)
	sliceA := make([]uint64, sliceLength)
	sliceB := make([]uint64, sliceLength)

	// fill input with random values
	for i := 0; i < sliceLength; i++ {
		sliceA[i] = uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
		sliceB[i] = uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
	}

	b.StartTimer()

	// run benchmark
	for n := 0; n < b.N; n++ {
		bitOpFunc(output, sliceA, sliceB)
	}
}

func BenchmarkAndSlice1Generic(b *testing.B) { benchmarkBitOp(b, 1, andSliceGeneric) }
func BenchmarkAndSlice1(b *testing.B)        { benchmarkBitOp(b, 1, andSlice) }

func BenchmarkAndSlice10Generic(b *testing.B) { benchmarkBitOp(b, 10, andSliceGeneric) }
func BenchmarkAndSlice10(b *testing.B)        { benchmarkBitOp(b, 10, andSlice) }

func BenchmarkAndSlice100Generic(b *testing.B) { benchmarkBitOp(b, 100, andSliceGeneric) }
func BenchmarkAndSlice100(b *testing.B)        { benchmarkBitOp(b, 100, andSlice) }

func BenchmarkAndSlice1000Generic(b *testing.B) { benchmarkBitOp(b, 1000, andSliceGeneric) }
func BenchmarkAndSlice1000(b *testing.B)        { benchmarkBitOp(b, 1000, andSlice) }

func BenchmarkAndSlice10000Generic(b *testing.B) { benchmarkBitOp(b, 10000, andSliceGeneric) }
func BenchmarkAndSlice10000(b *testing.B)        { benchmarkBitOp(b, 10000, andSlice) }

func BenchmarkOrSlice1Generic(b *testing.B) { benchmarkBitOp(b, 1, orSliceGeneric) }
func BenchmarkOrSlice1(b *testing.B)        { benchmarkBitOp(b, 1, orSlice) }

func BenchmarkOrSlice10Generic(b *testing.B) { benchmarkBitOp(b, 10, orSliceGeneric) }
func BenchmarkOrSlice10(b *testing.B)        { benchmarkBitOp(b, 10, orSlice) }

func BenchmarkOrSlice100Generic(b *testing.B) { benchmarkBitOp(b, 100, orSliceGeneric) }
func BenchmarkOrSlice100(b *testing.B)        { benchmarkBitOp(b, 100, orSlice) }

func BenchmarkOrSlice1000Generic(b *testing.B) { benchmarkBitOp(b, 1000, orSliceGeneric) }
func BenchmarkOrSlice1000(b *testing.B)        { benchmarkBitOp(b, 1000, orSlice) }

func BenchmarkOrSlice10000Generic(b *testing.B) { benchmarkBitOp(b, 10000, orSliceGeneric) }
func BenchmarkOrSlice10000(b *testing.B)        { benchmarkBitOp(b, 10000, orSlice) }

func BenchmarkXorSlice1Generic(b *testing.B) { benchmarkBitOp(b, 1, xorSliceGeneric) }
func BenchmarkXorSlice1(b *testing.B)        { benchmarkBitOp(b, 1, xorSlice) }

func BenchmarkXorSlice10Generic(b *testing.B) { benchmarkBitOp(b, 10, xorSliceGeneric) }
func BenchmarkXorSlice10(b *testing.B)        { benchmarkBitOp(b, 10, xorSlice) }

func BenchmarkXorSlice100Generic(b *testing.B) { benchmarkBitOp(b, 100, xorSliceGeneric) }
func BenchmarkXorSlice100(b *testing.B)        { benchmarkBitOp(b, 100, xorSlice) }

func BenchmarkXorSlice1000Generic(b *testing.B) { benchmarkBitOp(b, 1000, xorSliceGeneric) }
func BenchmarkXorSlice1000(b *testing.B)        { benchmarkBitOp(b, 1000, xorSlice) }

func BenchmarkXorSlice10000Generic(b *testing.B) { benchmarkBitOp(b, 10000, xorSliceGeneric) }
func BenchmarkXorSlice10000(b *testing.B)        { benchmarkBitOp(b, 10000, xorSlice) }

var notSliceTestCases = []struct {
	Input    []uint64
	Expected []uint64
}{
	{[]uint64{}, []uint64{}},
	{[]uint64{0x1}, []uint64{0xfffffffffffffffe}},
	{[]uint64{0x1, 0xffff, 0xffff0000, 0x0}, []uint64{0xfffffffffffffffe, 0xffffffffffff0000, 0xffffffff0000ffff, 0xffffffffffffffff}},
}

func testNotSlice(t *testing.T, bitOpFunc func(dest, src []uint64)) {
	for _, testCase := range notSliceTestCases {
		output := make([]uint64, len(testCase.Input))
		bitOpFunc(output, testCase.Input)
		if !reflect.DeepEqual(output, testCase.Expected) {
			t.Errorf("%s gave %v instead of expected %v for input %v", runtime.FuncForPC(reflect.ValueOf(bitOpFunc).Pointer()).Name(), output, testCase.Expected, testCase.Input)
		}
	}
}

func TestNotSliceGeneric(t *testing.T) { testNotSlice(t, notSliceGeneric) }
func TestNotSliceAsm(t *testing.T)     { testNotSlice(t, notSlice) }

func benchmarkNot(b *testing.B, sliceLength int, bitOpFunc func(dest, src []uint64)) {
	b.StopTimer()

	output := make([]uint64, sliceLength)
	input := make([]uint64, sliceLength)

	// fill input with random values
	for i := 0; i < sliceLength; i++ {
		input[i] = uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
	}

	b.StartTimer()

	// run benchmark
	for n := 0; n < b.N; n++ {
		bitOpFunc(output, input)
	}
}

func BenchmarkNotSlice1Generic(b *testing.B) { benchmarkNot(b, 1, notSliceGeneric) }
func BenchmarkNotSlice1(b *testing.B)        { benchmarkNot(b, 1, notSlice) }

func BenchmarkNotSlice10Generic(b *testing.B) { benchmarkNot(b, 10, notSliceGeneric) }
func BenchmarkNotSlice10(b *testing.B)        { benchmarkNot(b, 10, notSlice) }

func BenchmarkNotSlice100Generic(b *testing.B) { benchmarkNot(b, 100, notSliceGeneric) }
func BenchmarkNotSlice100(b *testing.B)        { benchmarkNot(b, 100, notSlice) }

func BenchmarkNotSlice1000Generic(b *testing.B) { benchmarkNot(b, 1000, notSliceGeneric) }
func BenchmarkNotSlice1000(b *testing.B)        { benchmarkNot(b, 1000, notSlice) }

func BenchmarkNotSlice10000Generic(b *testing.B) { benchmarkNot(b, 10000, notSliceGeneric) }
func BenchmarkNotSlice10000(b *testing.B)        { benchmarkNot(b, 10000, notSlice) }
