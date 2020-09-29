package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/kellegous/go/internal"
)

// RouteIterator allows iteration ont he named routes in the store
type RouteIterator struct {
	it   *redis.ScanIterator
	ctx  context.Context
	name string
	err  error
	cmd  *redis.ScanCmd
	pos  int
	rt   *internal.Route
}

// Valid checks if the current values of the Iterator are valid
func (i *RouteIterator) Valid() bool {
	// TODO implement me
	return i.cmd != nil
}

// re-implementing this with identical logic to satisfy the calling Interface
// https://github.com/go-redis/redis/blob/master/iterator.go
func (i *RouteIterator) Error() error {
	return i.it.Err()
}

// Seek ...
func (i *RouteIterator) Seek(s string) bool {
	// Since we have the current position available to us,
	// we should (hopefully) be able to avoid fully implementing this

	return i.it.Next(i.ctx)
}

// Name returns the iterator name
func (i *RouteIterator) Name() string {
	return i.name
}

// Route is the current route
func (i *RouteIterator) Route() *internal.Route {
	return i.rt
}

// Release should release any resources used by the Iterator
// Since the redis Iterator is safe to call concurrently, can we safely skip?
func (i *RouteIterator) Release() {
	return
}

// Next advances the Iterator to the next value
// will return true if more values can be read
func (i *RouteIterator) Next() bool {
	next := i.it.Next(i.ctx)
	if next {
		return false
	}
	return true
}
