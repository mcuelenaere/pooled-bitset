package pooled_bitset

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
)

func TestPopCountSlice(t *testing.T) {
	testCases := []struct {
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

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("version=generic/index=%d", idx), func (t *testing.T) {
			output := popcountSliceGeneric(testCase.Input)
			if output != testCase.Expected {
				t.Errorf("popCountSliceGeneric() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
			}
		})

		t.Run(fmt.Sprintf("version=asm/index=%d", idx), func (t *testing.T) {
			output := popcountSlice(testCase.Input)
			if output != testCase.Expected {
				t.Errorf("popcountSliceAsm() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
			}
		})
	}
}

// this is necessary because otherwise the inliner would unfairly give the non-assembly function an advantage
var indirectPopCountSliceGeneric = popcountSliceGeneric

func BenchmarkPopCountSlice(b *testing.B) {
	// generate slice containing increasing numbers
	slice := make([]uint64, 8192*10)
	for n := 0; n < len(slice); n++ {
		slice[n] = uint64(n)
	}

	b.Run("version=generic", func (b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			indirectPopCountSliceGeneric(slice)
		}
	})

	b.Run("version=asm", func (b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			popcountSlice(slice)
		}
	})

}

func TestFindFirstSetBit(t *testing.T) {
	testCases := []struct {
		Input    uint64
		Expected uint64
	}{
		{0x1, 0},
		{0x2, 1},
		{0x44, 2},
		{0x8000000000000000, 63},
		{0xFFFFFFFFFFFFFFFF, 0},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("version=generic/index=%d", idx), func (t *testing.T) {
			output := findFirstSetBitGeneric(testCase.Input)
			if output != testCase.Expected {
				t.Errorf("findFirstSetBitGeneric() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
			}
		})

		t.Run(fmt.Sprintf("version=asm/index=%d", idx), func (t *testing.T) {
			output := findFirstSetBit(testCase.Input)
			if output != testCase.Expected {
				t.Errorf("findFirstSetBitAsm() gave %d instead of expected %d for input 0x%x", output, testCase.Expected, testCase.Input)
			}
		})
	}
}

// this is necessary because otherwise the inliner would unfairly give the non-assembly function an advantage
var indirectFindFirstSetBitGeneric = findFirstSetBitGeneric

func BenchmarkFindFirstSetBit(b *testing.B) {
	b.Run("version=generic", func (b *testing.B) {
		for n := 0; n < b.N; n++ {
			indirectFindFirstSetBitGeneric(0x12345678)
		}
	})

	b.Run("version=asm", func (b *testing.B) {
		for n := 0; n < b.N; n++ {
			findFirstSetBit(0x12345678)
		}
	})

}

func TestBitOpSlice(t *testing.T) {
	configurations := []struct {
		Operator string
		Version string
		Func func(dest, a, b []uint64)
	}{
		{"AND", "generic", andSliceGeneric},
		{"AND", "asm", andSlice},
		{"OR", "generic", orSliceGeneric},
		{"OR", "asm", orSlice},
		{"XOR", "generic", xorSliceGeneric},
		{"XOR", "asm", xorSlice},
	}

	testCases := []struct {
		InputA      []uint64
		InputB      []uint64
		Expected    map[string][]uint64
	}{
		{[]uint64{}, []uint64{}, map[string][]uint64{"AND": []uint64{}, "OR": []uint64{}, "XOR": []uint64{}}},
		{[]uint64{0x1}, []uint64{0x0}, map[string][]uint64{"AND": []uint64{0x0}, "OR": []uint64{0x1}, "XOR": []uint64{0x1}}},
		{[]uint64{0x1, 0x1, 0x1, 0x1}, []uint64{0x0, 0x1, 0x0, 0x0}, map[string][]uint64{"AND": []uint64{0x0, 0x1, 0x0, 0x0}, "OR": []uint64{0x1, 0x1, 0x1, 0x1}, "XOR": []uint64{0x1, 0x0, 0x1, 0x1}}},
	}

	for _, configuration := range configurations {
		for idx, testCase := range testCases {
			t.Run(fmt.Sprintf("operator=%s/version=%s/index=%d", configuration.Operator, configuration.Version, idx), func(t *testing.T) {
				output := make([]uint64, len(testCase.InputA))
				expected := testCase.Expected[configuration.Operator]
				configuration.Func(output, testCase.InputA, testCase.InputB)
				if !reflect.DeepEqual(output, expected) {
					t.Errorf("%s gave %v instead of expected %v for inputs %v and %v", runtime.FuncForPC(reflect.ValueOf(configuration.Func).Pointer()).Name(), output, expected, testCase.InputA, testCase.InputB)
				}
			})
		}
	}
}

func BenchmarkBitOpSlice(b *testing.B) {
	functions := []struct {
		Operator string
		Version string
		Func func(dest, a, b []uint64)
	}{
		{"AND", "generic", andSliceGeneric},
		{"AND", "asm", andSlice},
		{"OR", "generic", orSliceGeneric},
		{"OR", "asm", orSlice},
		{"XOR", "generic", xorSliceGeneric},
		{"XOR", "asm", xorSlice},
	}
	sliceLengths := []int{1, 10, 100, 1000, 10000}

	for _, function := range functions {
		for _, sliceLength := range sliceLengths {
			b.Run(fmt.Sprintf("operator=%s/version=%s/length=%d", function.Operator, function.Version, sliceLength), func(b *testing.B) {
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
					function.Func(output, sliceA, sliceB)
				}
			})
		}
	}
}

func TestNotSlice(t *testing.T) {
	testCases := []struct {
		Input    []uint64
		Expected []uint64
	}{
		{[]uint64{}, []uint64{}},
		{[]uint64{0x1}, []uint64{0xfffffffffffffffe}},
		{[]uint64{0x1, 0xffff, 0xffff0000, 0x0}, []uint64{0xfffffffffffffffe, 0xffffffffffff0000, 0xffffffff0000ffff, 0xffffffffffffffff}},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("version=generic/index=%d", idx), func(t *testing.T) {
			output := make([]uint64, len(testCase.Input))
			notSliceGeneric(output, testCase.Input)
			if !reflect.DeepEqual(output, testCase.Expected) {
				t.Errorf("notSliceGeneric() gave %v instead of expected %v for input %v", output, testCase.Expected, testCase.Input)
			}
		})

		t.Run(fmt.Sprintf("version=asm/index=%d", idx), func(t *testing.T) {
			output := make([]uint64, len(testCase.Input))
			notSlice(output, testCase.Input)
			if !reflect.DeepEqual(output, testCase.Expected) {
				t.Errorf("notSliceGeneric() gave %v instead of expected %v for input %v", output, testCase.Expected, testCase.Input)
			}
		})
	}
}

func BenchmarkNotSlice(b *testing.B) {
	functions := []struct {
		Version string
		Func func(dest, src []uint64)
	}{
		{"generic", notSliceGeneric},
		{"asm", notSlice},
	}
	sliceLengths := []int{1, 10, 100, 1000, 10000}

	for _, function := range functions {
		for _, sliceLength := range sliceLengths {
			b.Run(fmt.Sprintf("version=%s/length=%d", function.Version, sliceLength), func(b *testing.B) {
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
					function.Func(output, input)
				}
			})
		}
	}
}
