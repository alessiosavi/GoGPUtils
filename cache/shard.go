package cache

import (
	"sync"
	"sync/atomic"
)

// entry is a cached item. prev points toward the tail (older entries) and
// next toward the head (newer entries).
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt int64 // monotonic ns since the cache epoch; 0 = never
	visited   atomic.Bool
	prev      *entry[K, V]
	next      *entry[K, V]
}

func (e *entry[K, V]) expired(now int64) bool {
	return e.expiresAt != 0 && now >= e.expiresAt
}

// shard is one independent SIEVE cache guarded by mu.
type shard[K comparable, V any] struct {
	mu       sync.RWMutex
	items    map[K]*entry[K, V]
	head     *entry[K, V] // newest
	tail     *entry[K, V] // oldest
	hand     *entry[K, V] // SIEVE hand; nil = start at tail
	capacity int          // 0 = unbounded
	stats    shardStats
}

func newShard[K comparable, V any](capacity int) *shard[K, V] {
	return &shard[K, V]{
		items:    make(map[K]*entry[K, V]),
		capacity: capacity,
	}
}
