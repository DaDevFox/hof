## Higher Order Functions in Go 1.23+ using Iterators and Generics

A checklist of JavaScript-style higher-order functions implemented in Go.

(Added `Square[E Number]([]E) iter.Seq[E]` and `Cube[E Number]([]E) iter.Seq[E]`)

---

## Core Array Methods

- [x] **`Map[E, T]([]E, func(E) T) iter.Seq[T]`** — Transform each element
- [x] **`Filter[E]([]E, func(E) bool) iter.Seq[E]`** — Keep elements that satisfy a condition
- [x] **`Reduce[T, U]([]T, func(U, T) U, init U) U`** — Accumulate values into one
- [x] **`ForEach[E]([]E, func(E))`** — Apply side-effects (printing, logging, etc.)
- [x] **`Find[E]([]E, func(E) bool) (E, bool)`** — Return first element satisfying condition
- [x] **`Some[E]([]E, func(E) bool) bool`** — Return `true` if any element matches
- [x] **`Every[E]([]E, func(E) bool) bool`** — Return `true` if all elements match

---

## Numeric Helpers

- [x] **`Square[E Number]([]E) iter.Seq[E]`** — Find the square of each element
- [x] **`Cube[E Number]([]E) iter.Seq[E]`** — Find the cube of each element
- [x] **`Sum[E Number]([]E) E`** — Add all numbers
- [x] **`Average[E Number]([]E) float64`** — Compute mean
- [x] **`Max[E Number]([]E) E`** — Find maximum value
- [x] **`Min[E Number]([]T) E`** — Find minimum value

---

## Collection Utilities

- [x] **`GroupBy[T, K comparable]([]T, func(T) K) map[K][]T`** — Cluster elements by key
- [x] **`Partition[T]([]T, func(T) bool) ([]T, []T)`** — Split into matching/non-matching
- [x] **`Unique[T comparable]([]T) []T`** — Remove duplicates
- [x] **`Zip[A, B]([]A, []B) [][2]any`** — Combine two slices
- [x] **`Unzip[A, B]([][2]any) ([]A, []B)`** — Split pairs
- [x] **`FlatMap[T, U]([]T, func(T) []U) []U`** — Map + flatten in one step
- [x] **`Chunk[T]([]T, size int) [][]T`** — Split slice into groups

---

## Functional Composition

- [x] **`Compose[A, B, C](f func(B) C, g func(A) B) func(A) C`** — Compose functions (right-to-left)
- [x] **`Pipe[A, B, C](f func(A) B, g func(B) C) func(A) C`** — Compose functions (left-to-right)
- [x] **`Curry` patterns using closures** — Turn multi-arg func into chain of funcs

## TODO: 

- [ ] parallel version of all existing functions
- [ ] N-way partition (int valued function)


---
