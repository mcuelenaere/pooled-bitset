package pooled_bitset

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"strings"
)

func getFunctionName(function interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(function).Pointer())
	return strings.Split(fn.Name(), ".")[1]
}

func generateSliceWithRandomData(sliceLength int) []uint64 {
	output := make([]uint64, sliceLength)
	for i := 0; i < sliceLength; i++ {
		output[i] = uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
	}
	return output
}

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

	for _, popcountFn := range popcountSliceVersions {
		for idx, testCase := range testCases {
			t.Run(fmt.Sprintf("version=%s/index=%d", getFunctionName(popcountFn), idx), func(t *testing.T) {
				output := popcountFn(testCase.Input)
				if output != testCase.Expected {
					t.Errorf("%s() gave %d instead of expected %d for input 0x%x", getFunctionName(popcountFn), output, testCase.Expected, testCase.Input)
				}
			})
		}
	}
}

func BenchmarkPopCountSlice(b *testing.B) {
	slice := generateSliceWithRandomData(8192*10)

	for _, popcountFn := range popcountSliceVersions {
		b.Run(fmt.Sprintf("version=%s", getFunctionName(popcountFn)), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				popcountFn(slice)
			}
		})
	}
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

	for _, findFirstSetBitFn := range findFirstSetBitVersions {
		for idx, testCase := range testCases {
			t.Run(fmt.Sprintf("version=%s/index=%d", getFunctionName(findFirstSetBitFn), idx), func(t *testing.T) {
				output := findFirstSetBitFn(testCase.Input)
				if output != testCase.Expected {
					t.Errorf("%s() gave %d instead of expected %d for input 0x%x", getFunctionName(findFirstSetBitFn), output, testCase.Expected, testCase.Input)
				}
			})
		}
	}
}

func BenchmarkFindFirstSetBit(b *testing.B) {
	for _, findFirstSetBitFn := range findFirstSetBitVersions {
		b.Run(fmt.Sprintf("version=%s", getFunctionName(findFirstSetBitFn)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				findFirstSetBitFn(0x12345678)
			}
		})
	}
}

func TestBitOpSlice(t *testing.T) {
	configurations := []struct {
		Operator string
		Funcs []func(dest, a, b []uint64)
	}{
		{"AND", andSliceVersions},
		{"OR", orSliceVersions},
		{"XOR", xorSliceVersions},
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
		for _, bitOpFn := range configuration.Funcs {
			for idx, testCase := range testCases {
				t.Run(fmt.Sprintf("operator=%s/version=%s/index=%d", configuration.Operator, getFunctionName(bitOpFn), idx), func(t *testing.T) {
					output := make([]uint64, len(testCase.InputA))
					expected := testCase.Expected[configuration.Operator]
					bitOpFn(output, testCase.InputA, testCase.InputB)
					if !reflect.DeepEqual(output, expected) {
						t.Errorf("%s() gave %v instead of expected %v for inputs %v and %v", getFunctionName(bitOpFn), output, expected, testCase.InputA, testCase.InputB)
					}
				})
			}
		}
	}
}

func BenchmarkBitOpSlice(b *testing.B) {
	functions := []struct {
		Operator string
		Funcs []func(dest, a, b []uint64)
	}{
		{"AND", andSliceVersions},
		{"OR", orSliceVersions},
		{"XOR", xorSliceVersions},
	}
	sliceLengths := []int{1, 10, 100, 1000, 10000}

	for _, sliceLength := range sliceLengths {
		for _, function := range functions {
			for _, bitOpFn := range function.Funcs {
				b.Run(fmt.Sprintf("operator=%s/version=%s/length=%d", function.Operator, getFunctionName(bitOpFn), sliceLength), func(b *testing.B) {
					b.StopTimer()

					output := make([]uint64, sliceLength)
					sliceA := generateSliceWithRandomData(sliceLength)
					sliceB := generateSliceWithRandomData(sliceLength)

					b.StartTimer()

					// run benchmark
					for n := 0; n < b.N; n++ {
						bitOpFn(output, sliceA, sliceB)
					}
				})
			}
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

	for _, notSliceFn := range notSliceVersions {
		for idx, testCase := range testCases {
			t.Run(fmt.Sprintf("version=%s/index=%d", getFunctionName(notSliceFn), idx), func(t *testing.T) {
				output := make([]uint64, len(testCase.Input))
				notSliceFn(output, testCase.Input)
				if !reflect.DeepEqual(output, testCase.Expected) {
					t.Errorf("%s() gave %v instead of expected %v for input %v", getFunctionName(notSliceFn), output, testCase.Expected, testCase.Input)
				}
			})
		}
	}
}

func BenchmarkNotSlice(b *testing.B) {
	sliceLengths := []int{1, 10, 100, 1000, 10000}

	for _, sliceLength := range sliceLengths {
		for _, notSliceFn := range notSliceVersions {
			b.Run(fmt.Sprintf("version=%s/length=%d", getFunctionName(notSliceFn), sliceLength), func(b *testing.B) {
				b.StopTimer()

				output := make([]uint64, sliceLength)
				input := generateSliceWithRandomData(sliceLength)

				b.StartTimer()

				// run benchmark
				for n := 0; n < b.N; n++ {
					notSliceFn(output, input)
				}
			})
		}
	}
}
