// Package pagination provides internal helpers for AWS paginated operations.
package pagination

import (
	"context"
)

// Paginator defines the interface for paginated AWS operations.
type Paginator[T any] interface {
	HasMorePages() bool
	NextPage(ctx context.Context) (T, error)
}

// Collect collects all items from a paginator using a result extractor function.
// The extractor function converts each page result to a slice of items.
//
// Example:
//
//	results, err := pagination.Collect(ctx, paginator, func(page *ListOutput) []Item {
//	    return page.Items
//	})
func Collect[T any, R any](ctx context.Context, p Paginator[T], extract func(T) []R) ([]R, error) {
	var results []R

	for p.HasMorePages() {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		page, err := p.NextPage(ctx)
		if err != nil {
			return results, err
		}

		items := extract(page)
		results = append(results, items...)
	}

	return results, nil
}

// CollectWithCallback iterates through pages and calls a callback for each.
// Useful when you want to process items as they arrive rather than collecting all.
// Return an error from the callback to stop pagination.
//
// Example:
//
//	err := pagination.CollectWithCallback(ctx, paginator, func(page *ListOutput) error {
//	    for _, item := range page.Items {
//	        if err := process(item); err != nil {
//	            return err
//	        }
//	    }
//	    return nil
//	})
func CollectWithCallback[T any](ctx context.Context, p Paginator[T], callback func(T) error) error {
	for p.HasMorePages() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		page, err := p.NextPage(ctx)
		if err != nil {
			return err
		}

		if err := callback(page); err != nil {
			return err
		}
	}

	return nil
}

// CollectN collects up to n items from a paginator.
// Returns early if n items are collected or pagination ends.
func CollectN[T any, R any](ctx context.Context, p Paginator[T], n int, extract func(T) []R) ([]R, error) {
	if n <= 0 {
		return nil, nil
	}

	results := make([]R, 0, n)

	for p.HasMorePages() && len(results) < n {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		page, err := p.NextPage(ctx)
		if err != nil {
			return results, err
		}

		items := extract(page)
		remaining := n - len(results)

		if len(items) <= remaining {
			results = append(results, items...)
		} else {
			results = append(results, items[:remaining]...)
		}
	}

	return results, nil
}
