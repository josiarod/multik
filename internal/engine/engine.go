package engine

import (
	"context"
	"sync"
	"time"
)

type Result[T any] struct {
	Key  string
	Item T
	Err  error
}

func FanOut[T any](
	ctx context.Context,
	keys []string,
	maxParallel int,
	timeoutPerKey time.Duration,
	fn func(context.Context, string) (T, error),
) []Result[T] {
	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup
	out := make(chan Result[T], len(keys))

	for _, k := range keys {
		k := k
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			ctx2, cancel := context.WithTimeout(ctx, timeoutPerKey)
			defer cancel()

			item, err := fn(ctx2, k)
			out <- Result[T]{Key: k, Item: item, Err: err}
		}()
	}

	go func() { wg.Wait(); close(out) }()

	var results []Result[T]
	for r := range out {
		results = append(results, r)
	}
	return results
}
