package hof_test

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/DaDevFox/hof"
)

func BenchmarkExtraneousMapStreaming(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for _, val := range inputArr {
			hof.Map(hof.Stream(val), func(j int) string {
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

func BenchmarkMap(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for _, val := range inputArr {
			hof.MapArray(val, func(j int) string {
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

func BenchmarkMapChainStreaming(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for _, val := range inputArr {
			slices.Collect(hof.Map(hof.Filter(hof.MapArray(val, func(j int) string {
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
			}),
				func(a string) bool {
					return len(a) < 10
				}),
				func(b string) string {
					return "padded: " + b
				}))
		}
	}
}

func BenchmarkMapChain(b *testing.B) {
	for b.Loop() {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
		}

		for _, val := range inputArr {
			hof.MapArray(slices.Collect(hof.FilterArray(slices.Collect(hof.MapArray(val, func(j int) string {
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
			})),
				func(a string) bool {
					return len(a) < 10
				})),
				func(b string) string {
					return "padded: " + b
				})
		}
	}
}

func TestMapArray(t *testing.T) {
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
			for i := range hof.MapArray(val, func(j int) string {
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
			}) {
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
			for sq := range hof.MapArray(val, func(x int) int { return x * x }) {
				req = append(req, sq)
			}
			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter coffee", func(t *testing.T) {
		inputArr := [][]string{
			{"tea", "soup", "coffee", "coke"},
			{"pan", "ghee", "mug", "spatula"},
			{"mustard", "egg", "meat"},
		}

		outputArr := [][]string{
			{"soup", "coke"},
			{"ghee"},
			{},
		}
		filterArr := []func(string) bool{
			func(eat string) bool {
				if eat == "soup" || eat == "coke" {
					return true
				}
				return false
			},
			func(utensil string) bool {
				return utensil == "ghee"
			},
			func(food string) bool {
				return food == "chicken"
			},
		}
		for idx, val := range inputArr {
			req := []string{}
			for v := range hof.FilterArray(val, filterArr[idx]) {
				req = append(req, v)
			}

			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("TestReduceSumInts", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		got := hof.Reduce(nums, func(acc, v int) int { return acc + v }, 0)
		want := 15
		if got != want {
			t.Errorf("Sum of ints: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceSumFloats", func(t *testing.T) {
		nums := []float64{1.2, 2.3, 3.4}
		got := hof.Reduce(nums, func(acc, v float64) float64 { return acc + v }, 0)
		want := 6.9
		if !floatAlmostEqual(got, want, 1e-9) {
			t.Errorf("Sum of floats: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceConcatStrings", func(t *testing.T) {
		words := []string{"Go", " ", "is", " ", "fun"}
		got := hof.Reduce(words, func(acc, v string) string { return acc + v }, "")
		want := "Go is fun"
		if got != want {
			t.Errorf("Concat strings: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceProductInts", func(t *testing.T) {
		nums := []int{1, 2, 3, 4}
		got := hof.Reduce(nums, func(acc, v int) int { return acc * v }, 1)
		want := 24
		if got != want {
			t.Errorf("Product of ints: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceCountCondition", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5, 6}
		got := hof.Reduce(nums, func(acc, v int) int {
			if v%2 == 0 {
				return acc + 1
			}
			return acc
		}, 0)
		want := 3 // even numbers: 2, 4, 6
		if got != want {
			t.Errorf("Count evens: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceStructs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{"A", 20},
			{"B", 25},
			{"C", 30},
		}

		got := hof.Reduce(people, func(acc int, p Person) int { return acc + p.Age }, 0)
		want := 75
		if got != want {
			t.Errorf("Sum ages: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduceEmptySlice", func(t *testing.T) {
		nums := []int{}
		got := hof.Reduce(nums, func(acc, v int) int { return acc + v }, 10)
		want := 10
		if got != want {
			t.Errorf("Empty slice: got %v, want %v", got, want)
		}
	})

	t.Run("TestReduce_StringCount", func(t *testing.T) {
		fruits := []string{"apple", "banana", "apple", "orange"}

		got := hof.Reduce(fruits, func(acc map[string]int, val string) map[string]int {
			acc[val]++
			return acc
		}, map[string]int{})

		want := map[string]int{"apple": 2, "banana": 1, "orange": 1}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reduce() mismatch\n got:  %v\n want: %v", got, want)
		}
	})
}

func floatAlmostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestForEach(t *testing.T) {
	t.Run("test array transformation", func(t *testing.T) {
		var buf bytes.Buffer

		inpArr := []int{1, 2, 3, 4, 5}
		hof.ForEachArray(inpArr, func(x int) {
			buf.WriteString(fmt.Sprintf("%d\n", x))
		})

		got := buf.String()
		want := "1\n2\n3\n4\n5\n"

		if got != want {
			t.Errorf("ForEach() = %q, want %q", got, want)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("find cherry", func(t *testing.T) {
		inputArr := []string{"banana", "apple", "cherry", "orange"}
		cherry, found := hof.Find(inputArr, func(fruit string) bool {
			return fruit == "cherry"
		})

		if !found {
			t.Error("not found")
		}

		if cherry != "cherry" {
			t.Errorf("got:%s\nwant:%s", cherry, "cherry")
		}
	})
	t.Run("find cherry not", func(t *testing.T) {
		inputArr := []string{"banana", "apple", "orange"}
		cherry, found := hof.Find(inputArr, func(fruit string) bool {
			return fruit == "cherry"
		})

		if found {
			t.Error("incorrect found")
		}

		if cherry != "" {
			t.Errorf("got:%s\nwant:%s", cherry, "")
		}
	})
}

func TestSome(t *testing.T) {
	t.Run("greater than 5", func(t *testing.T) {
		inputArr := []int{1, 2, 3, 6, 6, 7, 8}
		want := true
		if out := !hof.Some(inputArr, func(x int) bool {
			return x > 5
		}); out {
			t.Errorf("got:%v\nwant:%v", out, want)
		}
	})

	t.Run("greater than 5 not", func(t *testing.T) {
		inputArr := []int{1, 2, 3, 4}
		want := false
		if out := hof.Some(inputArr, func(x int) bool {
			return x > 5
		}); out {
			t.Errorf("got:%v\nwant:%v", out, want)
		}
	})
}

func TestEvery(t *testing.T) {
	t.Run("greater than 5", func(t *testing.T) {
		inputArr := []int{6, 6, 7, 8}
		want := true
		if out := !hof.Every(inputArr, func(x int) bool {
			return x > 5
		}); out {
			t.Errorf("got:%v\nwant:%v", out, want)
		}
	})

	t.Run("greater than 5 not", func(t *testing.T) {
		inputArr := []int{5, 7, 8, 1, 3}
		want := false
		if out := hof.Every(inputArr, func(x int) bool {
			return x > 5
		}); out {
			t.Errorf("got:%v\nwant:%v", out, want)
		}
	})
}

func TestSquare(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-4, 0, 69, 12},
			{0},
			{},
		}

		outputArr := [][]int{
			{1, 4, 9, 16, 25},
			{16, 0, 4761, 144},
			{0},
			{},
		}

		for idx, val := range inputArr {
			req := []int{}
			for sq := range hof.Square(val) {
				req = append(req, sq)
			}

			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
	t.Run("floats", func(t *testing.T) {
		inputArr := [][]float64{
			{1.9, 2.2, 3.2, 4.0, 5.1},
			{-4.9, 0.1, 69.42, 12.7},
			{0.0},
			{},
		}

		outputArr := [][]float64{
			{3.61, 4.84, 10.24, 16, 26.01},
			{24.01, 0.01, 4819.1364, 161.29},
			{0},
			{},
		}

		for idx, val := range inputArr {
			var req []float64
			for sq := range hof.Square(val) {
				req = append(req, sq)
			}

			if !slicesAlmostEqual(req, outputArr[idx], 1e-9) {
				t.Errorf("case %d:\ngot:  %v\nwant: %v", idx, req, outputArr[idx])
			}
		}
	})
}

func slicesAlmostEqual(a, b []float64, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > eps {
			return false
		}
	}
	return true
}

func TestCube(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		inputArr := [][]int{
			{1, 2, 3, 4, 5},
			{-2, 0, 3, -4},
			{0},
			{},
		}

		outputArr := [][]int{
			{1, 8, 27, 64, 125},
			{-8, 0, 27, -64},
			{0},
			{},
		}

		for idx, val := range inputArr {
			req := []int{}
			for cube := range hof.Cube(val) {
				req = append(req, cube)
			}

			if !reflect.DeepEqual(req, outputArr[idx]) {
				t.Errorf("got:%v\nwant:%v", req, outputArr[idx])
			}
		}
	})
	t.Run("floats", func(t *testing.T) {
		inputArr := [][]float64{
			{1.0, 2.0, 3.0},
			{-2.0, 0.5, 1.5},
			{0.0},
			{},
		}

		outputArr := [][]float64{
			{1.0, 8.0, 27.0},
			{-8.0, 0.125, 3.375},
			{0.0},
			{},
		}

		for idx, val := range inputArr {
			var req []float64
			for cube := range hof.Cube(val) {
				req = append(req, cube)
			}

			if !slicesAlmostEqual(req, outputArr[idx], 1e-9) {
				t.Errorf("case %d:\ngot:  %v\nwant: %v", idx, req, outputArr[idx])
			}
		}
	})
}

func TestSum(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testCases := []struct {
			input []int
			want  int
		}{
			{[]int{1, 2, 3, 4, 5}, 15},
			{[]int{-1, -2, -3}, -6},
			{[]int{0, 0, 0}, 0},
			{[]int{100}, 100},
			{[]int{}, 0},
			{[]int{-5, 5, -10, 10}, 0},
		}

		for _, tc := range testCases {
			got := hof.Sum(tc.input)
			if got != tc.want {
				t.Errorf("Sum(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})

	t.Run("floats", func(t *testing.T) {
		testCases := []struct {
			input []float64
			want  float64
		}{
			{[]float64{1.5, 2.5, 3.0}, 7.0},
			{[]float64{-1.1, -2.2, -3.3}, -6.6},
			{[]float64{0.0, 0.0}, 0.0},
			{[]float64{3.14}, 3.14},
			{[]float64{}, 0.0},
		}

		for _, tc := range testCases {
			got := hof.Sum(tc.input)
			if !floatAlmostEqual(got, tc.want, 1e-9) {
				t.Errorf("Sum(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})
}

func TestAverage(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testCases := []struct {
			input []int
			want  float64
		}{
			{[]int{1, 2, 3, 4, 5}, 3.0},
			{[]int{-2, -4, -6}, -4.0},
			{[]int{0, 0, 0}, 0.0},
			{[]int{100}, 100.0},
			{[]int{10, 20}, 15.0},
			{[]int{1, 3, 5, 7, 9}, 5.0},
		}

		for _, tc := range testCases {
			got := hof.Average(tc.input)
			if !floatAlmostEqual(got, tc.want, 1e-9) {
				t.Errorf("Average(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})

	t.Run("floats", func(t *testing.T) {
		testCases := []struct {
			input []float64
			want  float64
		}{
			{[]float64{1.0, 2.0, 3.0}, 2.0},
			{[]float64{-1.5, -2.5}, -2.0},
			{[]float64{5.5}, 5.5},
			{[]float64{0.1, 0.2, 0.3}, 0.2},
			{[]float64{10.5, 20.5, 30.0}, 20.333333333333332},
		}

		for _, tc := range testCases {
			got := hof.Average(tc.input)
			if !floatAlmostEqual(got, tc.want, 1e-9) {
				t.Errorf("Average(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testCases := []struct {
			input []int
			want  int
		}{
			{[]int{1, 2, 3, 4, 5}, 1},
			{[]int{5, 4, 3, 2, 1}, 1},
			{[]int{-1, -2, -3}, -3},
			{[]int{0, -5, 10}, -5},
			{[]int{100}, 100},
			{[]int{7, 7, 7}, 7},
			{[]int{}, 0}, // empty slice returns zero value
		}

		for _, tc := range testCases {
			got := hof.Min(tc.input)
			if got != tc.want {
				t.Errorf("Min(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})

	t.Run("floats", func(t *testing.T) {
		testCases := []struct {
			input []float64
			want  float64
		}{
			{[]float64{1.5, 2.5, 0.5}, 0.5},
			{[]float64{-1.1, -2.2, -0.5}, -2.2},
			{[]float64{3.14}, 3.14},
			{[]float64{5.0, 5.0, 5.0}, 5.0},
			{[]float64{}, 0.0}, // empty slice returns zero value
		}

		for _, tc := range testCases {
			got := hof.Min(tc.input)
			if !floatAlmostEqual(got, tc.want, 1e-9) {
				t.Errorf("Min(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testCases := []struct {
			input []int
			want  int
		}{
			{[]int{1, 2, 3, 4, 5}, 5},
			{[]int{5, 4, 3, 2, 1}, 5},
			{[]int{-1, -2, -3}, -1},
			{[]int{0, -5, 10}, 10},
			{[]int{100}, 100},
			{[]int{7, 7, 7}, 7},
		}

		for _, tc := range testCases {
			got := hof.Max(tc.input)
			if got != tc.want {
				t.Errorf("Max(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})

	t.Run("floats", func(t *testing.T) {
		testCases := []struct {
			input []float64
			want  float64
		}{
			{[]float64{1.5, 2.5, 0.5}, 2.5},
			{[]float64{-1.1, -2.2, -0.5}, -0.5},
			{[]float64{3.14}, 3.14},
			{[]float64{5.0, 5.0, 5.0}, 5.0},
		}

		for _, tc := range testCases {
			got := hof.Max(tc.input)
			if !floatAlmostEqual(got, tc.want, 1e-9) {
				t.Errorf("Max(%v) = %v, want %v", tc.input, got, tc.want)
			}
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		nums := []int{}
		got := hof.Max(nums)
		want := 0 // zero value for int
		if got != want {
			t.Errorf("Max(%v) = %v, want %v", nums, got, want)
		}
	})
}

func TestGroupBy(t *testing.T) {
	t.Run("group by length", func(t *testing.T) {
		words := []string{"cat", "dog", "bird", "fish", "elephant"}
		got := hof.GroupBy(words, func(s string) int { return len(s) })
		want := map[int][]string{
			3: {"cat", "dog"},
			4: {"bird", "fish"},
			8: {"elephant"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GroupBy() = %v, want %v", got, want)
		}
	})

	t.Run("group by first letter", func(t *testing.T) {
		words := []string{"apple", "banana", "apricot", "blueberry"}
		got := hof.GroupBy(words, func(s string) byte { return s[0] })
		want := map[byte][]string{
			'a': {"apple", "apricot"},
			'b': {"banana", "blueberry"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GroupBy() = %v, want %v", got, want)
		}
	})

	t.Run("group by even/odd", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6}
		got := hof.GroupBy(numbers, func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		})
		want := map[string][]int{
			"even": {2, 4, 6},
			"odd":  {1, 3, 5},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GroupBy() = %v, want %v", got, want)
		}
	})
}

func TestPartition(t *testing.T) {
	t.Run("partition even/odd", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		matched, rest := hof.Partition(numbers, func(n int) bool { return n%2 == 0 })

		wantMatched := []int{2, 4, 6, 8}
		wantRest := []int{1, 3, 5, 7}

		if !reflect.DeepEqual(matched, wantMatched) {
			t.Errorf("Partition() matched = %v, want %v", matched, wantMatched)
		}
		if !reflect.DeepEqual(rest, wantRest) {
			t.Errorf("Partition() rest = %v, want %v", rest, wantRest)
		}
	})

	t.Run("partition strings by length", func(t *testing.T) {
		words := []string{"cat", "elephant", "dog", "butterfly"}
		matched, rest := hof.Partition(words, func(s string) bool { return len(s) > 3 })

		wantMatched := []string{"elephant", "butterfly"}
		wantRest := []string{"cat", "dog"}

		if !reflect.DeepEqual(matched, wantMatched) {
			t.Errorf("Partition() matched = %v, want %v", matched, wantMatched)
		}
		if !reflect.DeepEqual(rest, wantRest) {
			t.Errorf("Partition() rest = %v, want %v", rest, wantRest)
		}
	})

	t.Run("partition empty slice", func(t *testing.T) {
		numbers := []int{}
		matched, rest := hof.Partition(numbers, func(n int) bool { return n > 0 })

		if len(matched) != 0 || len(rest) != 0 {
			t.Errorf("Partition() on empty slice should return empty slices, got matched=%v, rest=%v", matched, rest)
		}
	})
}

func TestUnique(t *testing.T) {
	t.Run("integers with duplicates", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 1, 4, 3, 5}
		got := hof.Unique(input)
		want := []int{1, 2, 3, 4, 5}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Unique(%v) = %v, want %v", input, got, want)
		}
	})

	t.Run("strings with duplicates", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "cherry", "banana"}
		got := hof.Unique(input)
		want := []string{"apple", "banana", "cherry"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Unique(%v) = %v, want %v", input, got, want)
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := hof.Unique(input)
		want := []int{1, 2, 3, 4, 5}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Unique(%v) = %v, want %v", input, got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := hof.Unique(input)
		want := []int{}

		if len(got) != 0 {
			t.Errorf("Unique(%v) = %v, want %v", input, got, want)
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("same length slices", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []string{"a", "b", "c"}
		got := hof.Zip(a, b)
		want := [][2]any{{1, "a"}, {2, "b"}, {3, "c"}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Zip(%v, %v) = %v, want %v", a, b, got, want)
		}
	})

	t.Run("different length slices", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []string{"a", "b", "c"}
		got := hof.Zip(a, b)
		want := [][2]any{{1, "a"}, {2, "b"}, {3, "c"}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Zip(%v, %v) = %v, want %v", a, b, got, want)
		}
	})

	t.Run("empty slices", func(t *testing.T) {
		a := []int{}
		b := []string{}
		got := hof.Zip(a, b)
		want := [][2]any{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Zip(%v, %v) = %v, want %v", a, b, got, want)
		}
	})
}

func TestUnzip(t *testing.T) {
	t.Run("unzip pairs", func(t *testing.T) {
		pairs := [][2]any{{1, "a"}, {2, "b"}, {3, "c"}}
		gotA, gotB := hof.Unzip[int, string](pairs)
		wantA := []int{1, 2, 3}
		wantB := []string{"a", "b", "c"}

		if !reflect.DeepEqual(gotA, wantA) {
			t.Errorf("Unzip() first slice = %v, want %v", gotA, wantA)
		}
		if !reflect.DeepEqual(gotB, wantB) {
			t.Errorf("Unzip() second slice = %v, want %v", gotB, wantB)
		}
	})

	t.Run("empty pairs", func(t *testing.T) {
		pairs := [][2]any{}
		gotA, gotB := hof.Unzip[int, string](pairs)

		if len(gotA) != 0 || len(gotB) != 0 {
			t.Errorf("Unzip() on empty pairs should return empty slices, got A=%v, B=%v", gotA, gotB)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("split strings into characters", func(t *testing.T) {
		words := []string{"cat", "dog"}
		got := hof.FlatMap(words, func(s string) []string {
			chars := make([]string, len(s))
			for i, r := range s {
				chars[i] = string(r)
			}
			return chars
		})
		want := []string{"c", "a", "t", "d", "o", "g"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("FlatMap() = %v, want %v", got, want)
		}
	})

	t.Run("duplicate numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		got := hof.FlatMap(numbers, func(n int) []int {
			return []int{n, n}
		})
		want := []int{1, 1, 2, 2, 3, 3}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("FlatMap() = %v, want %v", got, want)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		got := hof.FlatMap(numbers, func(n int) []int {
			return []int{}
		})

		if len(got) != 0 {
			t.Errorf("FlatMap() = %v, want empty slice", got)
		}
	})
}

func TestChunk(t *testing.T) {
	t.Run("chunk into size 2", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7}
		got := hof.Chunk(input, 2)
		want := [][]int{{1, 2}, {3, 4}, {5, 6}, {7}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Chunk(%v, 2) = %v, want %v", input, got, want)
		}
	})

	t.Run("chunk into size 3", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		got := hof.Chunk(input, 3)
		want := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Chunk(%v, 3) = %v, want %v", input, got, want)
		}
	})

	t.Run("chunk size larger than slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := hof.Chunk(input, 5)
		want := [][]int{{1, 2, 3}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Chunk(%v, 5) = %v, want %v", input, got, want)
		}
	})

	t.Run("chunk with zero size", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := hof.Chunk(input, 0)

		if got != nil {
			t.Errorf("Chunk(%v, 0) = %v, want nil", input, got)
		}
	})

	t.Run("chunk with negative size", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := hof.Chunk(input, -1)

		if got != nil {
			t.Errorf("Chunk(%v, -1) = %v, want nil", input, got)
		}
	})

	t.Run("chunk empty slice", func(t *testing.T) {
		input := []int{}
		got := hof.Chunk(input, 2)

		if len(got) != 0 {
			t.Errorf("Chunk(%v, 2) = %v, want empty slice", input, got)
		}
	})
}

func TestCompose(t *testing.T) {
	t.Run("compose int functions", func(t *testing.T) {
		addOne := func(x int) int { return x + 1 }
		multiplyTwo := func(x int) int { return x * 2 }

		composed := hof.Compose(multiplyTwo, addOne)
		got := composed(5)
		want := 12 // (5 + 1) * 2

		if got != want {
			t.Errorf("Compose(multiplyTwo, addOne)(5) = %v, want %v", got, want)
		}
	})

	t.Run("compose string functions", func(t *testing.T) {
		addExclamation := func(s string) string { return s + "!" }
		toUpper := func(s string) string { return "UPPER:" + s }

		composed := hof.Compose(toUpper, addExclamation)
		got := composed("hello")
		want := "UPPER:hello!"

		if got != want {
			t.Errorf("Compose(toUpper, addExclamation)(\"hello\") = %v, want %v", got, want)
		}
	})
}

func TestPipe(t *testing.T) {
	t.Run("pipe int functions", func(t *testing.T) {
		addOne := func(x int) int { return x + 1 }
		multiplyTwo := func(x int) int { return x * 2 }

		piped := hof.Pipe(addOne, multiplyTwo)
		got := piped(5)
		want := 12 // (5 + 1) * 2

		if got != want {
			t.Errorf("Pipe(addOne, multiplyTwo)(5) = %v, want %v", got, want)
		}
	})

	t.Run("pipe string functions", func(t *testing.T) {
		addExclamation := func(s string) string { return s + "!" }
		toUpper := func(s string) string { return "UPPER:" + s }

		piped := hof.Pipe(addExclamation, toUpper)
		got := piped("hello")
		want := "UPPER:hello!"

		if got != want {
			t.Errorf("Pipe(addExclamation, toUpper)(\"hello\") = %v, want %v", got, want)
		}
	})
}

func TestCurry(t *testing.T) {
	t.Run("curry add function", func(t *testing.T) {
		add := func(a, b int) int { return a + b }
		curriedAdd := hof.Curry(add)

		addFive := curriedAdd(5)
		got := addFive(3)
		want := 8

		if got != want {
			t.Errorf("Curry(add)(5)(3) = %v, want %v", got, want)
		}
	})

	t.Run("curry string concat", func(t *testing.T) {
		concat := func(a, b string) string { return a + b }
		curriedConcat := hof.Curry(concat)

		addHello := curriedConcat("Hello, ")
		got := addHello("World!")
		want := "Hello, World!"

		if got != want {
			t.Errorf("Curry(concat)(\"Hello, \")(\"World!\") = %v, want %v", got, want)
		}
	})

	t.Run("curry multiply", func(t *testing.T) {
		multiply := func(a, b float64) float64 { return a * b }
		curriedMultiply := hof.Curry(multiply)

		double := curriedMultiply(2.0)
		got := double(3.5)
		want := 7.0

		if !floatAlmostEqual(got, want, 1e-9) {
			t.Errorf("Curry(multiply)(2.0)(3.5) = %v, want %v", got, want)
		}
	})
}

// Early Termination Tests for Iterator Functions

func TestMapEarlyTermination(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	var result []string

	// Break after processing 2 elements to trigger early return in Map
	for val := range hof.MapArray(input, func(i int) string {
		return fmt.Sprintf("item_%d", i)
	}) {
		result = append(result, val)
		if len(result) >= 2 {
			break // This will cause yield to return false
		}
	}

	expected := []string{"item_1", "item_2"}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFilterEarlyTermination(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var result []int

	// Break after finding 2 even numbers to trigger early return in Filter
	for val := range hof.FilterArray(input, func(i int) bool {
		return i%2 == 0
	}) {
		result = append(result, val)
		if len(result) >= 2 {
			break // This will cause yield to return false
		}
	}

	expected := []int{2, 4}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSquareEarlyTermination(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	var result []int

	// Break after processing 3 elements to trigger early return in Square
	for val := range hof.Square(input) {
		result = append(result, val)
		if len(result) >= 3 {
			break // This will cause yield to return false
		}
	}

	expected := []int{1, 4, 9}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCubeEarlyTermination(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	var result []int

	// Break after processing 3 elements to trigger early return in Cube
	for val := range hof.Cube(input) {
		result = append(result, val)
		if len(result) >= 3 {
			break // This will cause yield to return false
		}
	}

	expected := []int{1, 8, 27}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
