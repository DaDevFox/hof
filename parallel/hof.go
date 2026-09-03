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
