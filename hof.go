// Package hof: Higher Order Functions in Golang Using Iterators and Generics
package hof

import (
	"iter"
	"slices"
)

// Core Array Methods

// Filter : Keep elements that satisfy a condition
func Filter[E any](arr []E, filter func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range arr {
func Stream[T any](arr []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range arr {
			proceed := yield(v)
			if !proceed {
				return
			}
		}
	}
}

// Map : Transform each element
// Collected (array) input, streaming output in serial, forwards direction
func MapArray[E, T any](arr []E, transform func(E) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range arr {
			if !yield(transform(v)) {
				return
			}
		}
	}
}

// Map : Transform each element
// Collected (array) input + output in serial, forwards direction
func MapToArray[E, T any](arr []E, transform func(E) T) []T {
	return slices.Collect(MapArray(arr, transform))
}

// Map : Transform each element
// Bidi streaming in serial, forwards direction
func Map[E, T any](stream iter.Seq[E], transform func(E) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range stream {
			if !yield(transform(v)) {
				return
			}
		}
	}
}
			if filter(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Reduce : Accumulate values into one
func Reduce[E any, T any](arr []E, fn func(T, E) T, init T) T {
	acc := init
	for _, v := range arr {
		acc = fn(acc, v)
	}
	return acc
}

// ForEach : Apply side-effects (printing, logging, etc.)
func ForEach[E any](arr []E, fn func(E)) {
	for _, v := range arr {
		fn(v)
	}
}

// Find : Return first element satisfying condition
func Find[E any](arr []E, fn func(E) bool) (E, bool) {
	var out E
	for _, v := range arr {
		if fn(v) {
			return v, true
		}
	}
	return out, false
}

// Some : Return true if any element matches
func Some[E any](arr []E, fn func(E) bool) bool {
	return slices.ContainsFunc(arr, fn)
}

// Every : Return true if all elements match
func Every[E any](arr []E, fn func(E) bool) bool {
	for _, v := range arr {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Numeric Helpers

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Square : Find the square of each element
func Square[E Number](arr []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range arr {
			if !yield(v * v) {
				return
			}
		}
	}
}

// Cube : Find the cube of each element
func Cube[E Number](arr []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range arr {
			if !yield(v * v * v) {
				return
			}
		}
	}
}

// Sum : Add all numbers
func Sum[E Number](arr []E) E {
	var sum E
	for _, v := range arr {
		sum += v
	}
	return sum
}

// Average : Compute mean
func Average[E Number](arr []E) float64 {
	var sum E
	for _, v := range arr {
		sum += v
	}
	return float64(sum) / float64(len(arr))
}

// Min : Find maximum value
func Min[E Number](arr []E) E {
	switch len(arr) {
	case 0:
		var min_ E
		return min_
	case 1:
		return arr[0]
	default:
		min_ := arr[0]
		for _, v := range arr {
			if min_ > v {
				min_ = v
			}
		}
		return min_
	}
}

// Max : Find minimum value
func Max[E Number](arr []E) E {
	switch len(arr) {
	case 0:
		var max_ E
		return max_
	case 1:
		return arr[0]
	default:
		max_ := arr[0]
		for _, v := range arr {
			if max_ < v {
				max_ = v
			}
		}
		return max_
	}
}

// Collection Utilities

// GroupBy : Cluster elements by key
func GroupBy[T any, K comparable](arr []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, v := range arr {
		key := keyFn(v)
		groups[key] = append(groups[key], v)
	}
	return groups
}

// Partition : Split into matching/non-matching
func Partition[T any](arr []T, fn func(T) bool) ([]T, []T) {
	var matched, rest []T
	for _, v := range arr {
		if fn(v) {
			matched = append(matched, v)
		} else {
			rest = append(rest, v)
		}
	}
	return matched, rest
}

// Unique : Remove duplicates
func Unique[T comparable](arr []T) []T {
	seen := make(map[T]struct{})
	var result []T
	for _, v := range arr {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Zip : Combine two slices
func Zip[A, B any](a []A, b []B) [][2]any {
	n := len(a)
	n = min(n, len(b))
	result := make([][2]any, n)
	for i := 0; i < n; i++ {
		result[i] = [2]any{a[i], b[i]}
	}
	return result
}

// Unzip : Split pairs
func Unzip[A, B any](pairs [][2]any) ([]A, []B) {
	var aVals []A
	var bVals []B
	for _, pair := range pairs {
		aVals = append(aVals, pair[0].(A))
		bVals = append(bVals, pair[1].(B))
	}
	return aVals, bVals
}

// FlatMap : Map + flatten in one step
func FlatMap[T any, U any](arr []T, fn func(T) []U) []U {
	var result []U
	for _, v := range arr {
		result = append(result, fn(v)...)
	}
	return result
}

// Chunk : Split slice into groups
func Chunk[T any](arr []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var chunks [][]T
	for i := 0; i < len(arr); i += size {
		end := i + size
		end = min(end, len(arr))
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// Functional Composition

// Compose : Compose functions (right-to-left)
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
	return func(x A) C {
		return f(g(x))
	}
}

// Pipe : Compose functions (left-to-right)
func Pipe[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(x A) C {
		return g(f(x))
	}
}

// Curry : Turn multi-arg func into chain of funcs
func Curry[A, B, C any](fn func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C {
			return fn(a, b)
		}
	}
}
