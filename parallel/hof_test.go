package parallel_test

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	// "slices"
	"testing"

	hof "github.com/DaDevFox/hof"
	hofparallel "github.com/DaDevFox/hof/parallel"
)

var THREADS_TO_TEST = 2

var MAXARRSIZE = 1000
var UNEVEN_WORK_FREQUENCY = 5
var UNEVEN_WORK_INTENSITY = 1600

func BenchmarkMapParallel_Uneven_Strided(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for range MAXARRSIZE {
			arr := []int{}
			for i := range MAXARRSIZE {
				arr = append(arr, i)
			}
			inputArr = append(inputArr, arr)
		}

		for _, val := range inputArr {
			hofparallel.MapParallel_Strided(val, func(j int) string {
				switch {
				case j > 0:
					if 10*j < UNEVEN_WORK_FREQUENCY {
						return fmt.Sprintf("special sum is %d: ", hof.Sum(hof.MapToArray(make([]int, UNEVEN_WORK_INTENSITY), func(k int) int {
							return 5
						})))
					}
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					return fmt.Sprintf("incorrect input got %d", j)
				}
			}, THREADS_TO_TEST)
		}
	}
}

func BenchmarkMapParallel_Uneven_CacheLinear(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for range MAXARRSIZE {
			arr := []int{}
			for i := range MAXARRSIZE {
				arr = append(arr, i)
			}
			inputArr = append(inputArr, arr)
		}

		for _, val := range inputArr {
			hofparallel.MapParallel_CacheLinear(val, func(j int) string {
				switch {
				case j > 0:
					if j < UNEVEN_WORK_FREQUENCY {
						arr := make([]int, UNEVEN_WORK_INTENSITY)
						for i := range UNEVEN_WORK_INTENSITY {
							arr[i] = 4
						}
						return fmt.Sprintf("special sum is %d: ", hof.Sum(hof.MapToArray(arr, func(k int) int {
							return 5
						})))
					}
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					return fmt.Sprintf("incorrect input got %d", j)
				}
			}, THREADS_TO_TEST)
		}
	}
}

// func BenchmarkMapParallel_Gregwebs(b *testing.B) {
// 	for b.Loop() {
// 		inputArr := [][]int{
// 			{1, 2, 3, 4, 5},
// 			{-4, 0, 69, 12},
// 			{0},
// 		}
//
// 		for range MAXARRSIZE {
// 			arr := []int{}
// 			for range MAXARRSIZE {
// 				arr = append(arr, rand.IntN(1000))
// 			}
// 			inputArr = append(inputArr, arr)
// 		}
//
// 		for _, val := range inputArr {
// 			hofparallel.MapParallel_Gregwebs(val, func(j int) string {
// 				switch {
// 				case j > 0:
// 					return fmt.Sprintf("this is %d", j)
// 				case j == 0:
// 					return "this is empty"
// 				case j < 0:
// 					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
// 				default:
// 					return fmt.Sprintf("incorrect input got %d", j)
// 				}
// 			}, 2)
// 		}
// 	}
// }

// func BenchmarkMapParallelStreaming(b *testing.B) {
// 	for b.Loop() {
// 		inputArr := [][]int{
// 			{1, 2, 3, 4, 5},
// 			{-4, 0, 69, 12},
// 			{0},
// 		}
//
// 		for range MAXARRSIZE {
// 			arr := []int{}
// 			for range MAXARRSIZE {
// 				arr = append(arr, rand.IntN(1000))
// 			}
// 			inputArr = append(inputArr, arr)
// 		}
//
// 		for _, val := range inputArr {
// 			slices.Collect(hofparallel.MapParallelStreaming(hof.Stream(val), func(j int) string {
// 				switch {
// 				case j > 0:
// 					return fmt.Sprintf("this is %d", j)
// 				case j == 0:
// 					return "this is empty"
// 				case j < 0:
// 					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
// 				default:
// 					return fmt.Sprintf("incorrect input got %d", j)
// 				}
// 			}, THREADS_TO_TEST))
// 		}
// 	}
// }

func BenchmarkMapParallel(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for range MAXARRSIZE {
			arr := []int{}
			for range MAXARRSIZE {
				arr = append(arr, rand.IntN(1000))
			}
			inputArr = append(inputArr, arr)
		}

		for _, val := range inputArr {
			hofparallel.MapParallel(val, func(j int) string {
				switch {
				case j > 0:
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					return fmt.Sprintf("incorrect input got %d", j)
				}
			}, THREADS_TO_TEST)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for range MAXARRSIZE {
			arr := []int{}
			for range MAXARRSIZE {
				arr = append(arr, rand.IntN(1000))
			}
			inputArr = append(inputArr, arr)
		}

		for _, val := range inputArr {
			hof.MapToArray(val, func(j int) string {
				switch {
				case j > 0:
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					return fmt.Sprintf("incorrect input got %d", j)
				}
			})
		}
	}
}

func TestMapParallelStreaming(t *testing.T) {
	t.Run("int to main", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		outputArr := [][]string{
			{"this is 1", "this is 2", "this is 3", "this is 4", "this is 5"},
			{"this is a neg 4", "this is empty", "this is 69", "this is 12"},
			{"this is empty"},
		}

		for idx, val := range inputArr {
			var req []string
			for i := range hofparallel.MapParallelStreaming(hof.Stream(val), func(j int) string {
				fmt.Printf("hei: %d\n", j)
				switch {
				case j > 0:
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					t.Errorf("incorrect input got %d", j)
				}
				return ""
			}, 4) {
				req = append(req, i)
			}

			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})

	t.Run("square", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		outputArr := [][]int{
			{1, 4, 9, 16, 25},
			{16, 0, 4761, 144},
			{0},
		}

		for idx, val := range inputArr {
			var req []int
			for _, sq := range hofparallel.MapParallel(val, func(x int) int { return x * x }, 4) {
				req = append(req, sq)
			}
			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
}

func TestMapParallel(t *testing.T) {
	t.Run("int to main", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		outputArr := [][]string{
			{"this is 1", "this is 2", "this is 3", "this is 4", "this is 5"},
			{"this is a neg 4", "this is empty", "this is 69", "this is 12"},
			{"this is empty"},
		}

		for idx, val := range inputArr {
			var req []string
			for _, i := range hofparallel.MapParallel(val, func(j int) string {
				switch {
				case j > 0:
					return fmt.Sprintf("this is %d", j)
				case j == 0:
					return "this is empty"
				case j < 0:
					return fmt.Sprintf("this is a neg %d", int(math.Abs(float64(j))))
				default:
					t.Errorf("incorrect input got %d", j)
				}
				return ""
			}, 4) {
				req = append(req, i)
			}

			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})

	t.Run("square", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		outputArr := [][]int{
			{1, 4, 9, 16, 25},
			{16, 0, 4761, 144},
			{0},
		}

		for idx, val := range inputArr {
			var req []int
			for _, sq := range hofparallel.MapParallel(val, func(x int) int { return x * x }, 4) {
				req = append(req, sq)
			}
			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
}
