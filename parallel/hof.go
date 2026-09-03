package parallel

import (
	"fmt"
	"iter"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quintans/toolkit/latch"

	dgravesa "github.com/dgravesa/go-parallel/parallel"
	gregwebs "github.com/gregwebs/go-parallel"
)

func MapParallelStreaming[E, T any](stream iter.Seq[E], transform func(E) T, threads int) iter.Seq[T] {
	fmt.Println("starting call")
	working := []E{}
	completed := []T{}
	cancel := make(chan bool)

	val := false
	done := &val

	var completionCursor int32
	completionCursor = -1

	var mu sync.Mutex

	for i := range threads {
		go func() {
			cursor := i
			for {
				if cursor >= len(working) {
					select {
					case <-cancel: // check cancel channel for completion signal
						fmt.Printf("t%d cancelling\n", i)
						return
					default:
					}

					continue
				}

				mu.Lock()
				if cursor >= len(working) || cursor >= len(completed) {
					fmt.Printf("this should never happen (out of bounds): %d, %d, %d", cursor, len(working), len(completed))
				}
				completed[cursor] = transform(working[cursor])
				mu.Unlock()

				if completionCursor == int32(cursor) {
					fmt.Printf("this should never happen")
				}

				fmt.Printf("waiting at cursor %d\n", cursor)
				for completionCursor < int32(cursor) && completionCursor != int32(cursor-1) {
				}

				if completionCursor == int32(cursor-1) {
					atomic.AddInt32(&completionCursor, 1)
					fmt.Printf("bump ->%d\n", completionCursor)
				}

				cursor += threads
			}
		}()
	}

	go func() {
		for item := range stream {

			// critical section
			mu.Lock()
			working = append(working, item)
			completed = slices.Concat(completed, make([]T, 1))
			mu.Unlock()
			// critical section

			if len(working) != len(completed) {
				fmt.Printf("this should never happen (working, compelted out of sync)")
			}
		}

		fmt.Printf("whole input was %d long\n", len(completed))

		for range threads {
			cancel <- true
		}

		fmt.Printf("done!!")
		*done = true
	}()

	var streamingCursor int32
	streamingCursor = 0

	return func(yield func(T) bool) {
		for {
			if streamingCursor == completionCursor && *done {
				yield(completed[streamingCursor])
				return
			}

			for ; streamingCursor < completionCursor; streamingCursor++ {
				if !yield(completed[streamingCursor]) {
					return
				}
			}
		}
	}
}

func MapParallelStreamingUnordered[E, T any](stream iter.Seq[E], transform func(E) T, threads int) iter.Seq[T] {
	working := []E{}

	val := false
	done := &val
	cancel := make(chan bool)
	completed := make(chan T)

	for i := range threads {
		go func() {
			cursor := i
			for {
				if cursor >= len(working) {
					select {
					case <-cancel: // check cancel channel for completion signal
						return
					default:
					}
					continue
				}

				completed <- transform(working[cursor])
				cursor += threads
			}
		}()
	}

	go func() {
		for item := range stream {
			working = append(working, item)
		}

		for range threads {
			cancel <- true
		}

		*done = true
	}()

	return func(yield func(T) bool) {
		for {
			worked := false

			select {
			case item := <-completed:
				fmt.Printf("completed: %v\n", item)
				worked = true
				if !yield(item) {
					return
				}
			default:
			}

			if !worked && *done {
				return
			}
		}
	}
}

func MapParallel[E, T any](arr []E, transform func(E) T, threads int) []T {
	return MapParallel_Strided(arr, transform, threads)
}

// version of map parallel in which chunks of paarllel processing units work on contiguous memory units, prioritizing cache linearity over even work distribution (which would imply a strided approach)
func MapParallel_CacheLinear[E, T any](arr []E, transform func(E) T, threads int) []T {
	result := make([]T, len(arr))
	latch := latch.NewCountDownLatch()
	latch.Add(threads)

	// wait := make(chan error)

	for i := range threads {
		go func() {
			batchsize := len(arr) / threads
			idx := 0
			if len(arr)%threads == 0 {
				idx = batchsize * i
			} else if i < len(arr)%threads {
				idx = (batchsize + 1) * i
				batchsize += 1
			} else {
				idx = (batchsize+1)*(len(arr)%threads) + batchsize*(i-(len(arr)%threads))
			}
			for range batchsize {
				if idx >= len(arr) {
					continue
				}
				result[idx] = transform(arr[idx])
				idx++
			}
			latch.Done()
			// wait <- nil
		}()
	}

	// for range threads {
	// 	<-wait
	// }

	timeout := latch.WaitWithTimeout(time.Second)
	if timeout {
		// TODO; figure out here
	}

	return result
}

func MapParallel_Strided[E, T any](arr []E, transform func(E) T, threads int) []T {
	result := make([]T, len(arr))
	latch := latch.NewCountDownLatch()
	latch.Add(threads)

	// wait := make(chan error)

	for i := range threads {
		go func() {
			for idx := i; idx < len(arr); idx += threads {
				result[idx] = transform(arr[idx])
			}
			latch.Done()
			// wait <- nil
		}()
	}

	// for range threads {
	// 	<-wait
	// }

	timeout := latch.WaitWithTimeout(time.Second)
	if timeout {
		// TODO; figure out here
	}

	return result
}

func MapParallel_Dgravesa[E, T any](arr []E, transform func(E) T, threads int) []T {
	result := make([]T, len(arr))

	dgravesa.SetDefaultNumGoroutines(threads)

	dgravesa.For(len(arr), func(idx, _ int) {
		result[idx] = transform(arr[idx])
	})

	// if errors != nil {
	// 	//TODO: log or force passing failure delegate [decide which]
	// }
	return result
}

func MapParallel_Gregwebs[E, T any](arr []E, transform func(E) T, threads int) []T {
	parts := make([]T, len(arr))

	gregwebs.CollectErrors(gregwebs.ArrayWorkers1(threads, arr, make(chan struct{}), func(idx int, element E) error {
		parts[idx] = transform(element)
		return nil
	}))

	// if errors != nil {
	// 	//TODO: log or force passing failure delegate [decide which]
	// }
	return parts
}
