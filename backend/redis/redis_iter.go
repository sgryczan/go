package redis

import (
	"sync"

	redis "github.com/go-redis/redis"
	"github.com/kellegous/go/internal"
)

// RouteIterator allows iteration ont he named routes in the store
type RouteIterator struct {
	redis.ScanIterator
	name string
	err  error
	mu   sync.Mutex
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
	return i.Err()
}

// Seek ...
func (i *RouteIterator) Seek(s string) bool {
	// Since we have the current position available to us,
	// we should (hopefully) be able to avoid fully implementing this

	return i.Next()
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
