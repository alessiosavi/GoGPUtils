package cache

import (
	"iter"
	"sync"
	"sync/atomic"
	"time"
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
	inflight map[K]*call[V]
	stats    shardStats
}

func newShard[K comparable, V any](capacity int) *shard[K, V] {
	return &shard[K, V]{
		items:    make(map[K]*entry[K, V]),
		capacity: capacity,
		inflight: make(map[K]*call[V]),
	}
}

// pushHead links e as the newest entry. The caller holds s.mu.
func (s *shard[K, V]) pushHead(e *entry[K, V]) {
	e.prev, e.next = s.head, nil
	if s.head != nil {
		s.head.next = e
	} else {
		s.tail = e
	}
	s.head = e
}

// unlink removes e from the list. The caller holds s.mu.
func (s *shard[K, V]) unlink(e *entry[K, V]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		s.tail = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		s.head = e.prev
	}
}

// remove deletes e, moves the hand off it, updates counters, records the
// removal, and clears e's references so its value can be collected. The
// caller holds s.mu.
func (s *shard[K, V]) remove(e *entry[K, V], reason EvictionReason, rec *[]removal[K, V]) {
	if s.hand == e {
		s.hand = e.next
	}
	s.unlink(e)
	delete(s.items, e.key)

	switch reason {
	case ReasonEvicted:
		s.stats.evictions.Add(1)
	case ReasonExpired:
		s.stats.expirations.Add(1)
	}

	if rec != nil {
		*rec = append(*rec, removal[K, V]{key: e.key, value: e.value, reason: reason})
	}

	var zero V
	e.value, e.prev, e.next = zero, nil, nil
}

// evict removes one entry with SIEVE: the hand walks from the tail toward the
// head (wrapping), clearing visited bits, and removes the first expired or
// unvisited entry. The caller holds s.mu.
func (s *shard[K, V]) evict(now int64, rec *[]removal[K, V]) {
	e := s.hand
	if e == nil {
		e = s.tail
	}

	for e != nil {
		expired := e.expired(now)
		if !expired && e.visited.Load() {
			e.visited.Store(false)
			if e = e.next; e == nil {
				e = s.tail
			}

			continue
		}

		next := e.next // saved before remove clears it
		reason := ReasonEvicted
		if expired {
			reason = ReasonExpired
		}
		s.remove(e, reason, rec)
		s.hand = next

		return
	}
}

// set inserts or replaces key. Replacing keeps the entry's position and marks
// it visited. The caller holds s.mu.
func (s *shard[K, V]) set(key K, value V, expiresAt, now int64, rec *[]removal[K, V]) {
	if e, ok := s.items[key]; ok {
		reason := ReasonReplaced
		if e.expired(now) {
			reason = ReasonExpired
			s.stats.expirations.Add(1)
		}
		if rec != nil {
			*rec = append(*rec, removal[K, V]{key: key, value: e.value, reason: reason})
		}
		e.value, e.expiresAt = value, expiresAt
		e.visited.Store(true)

		return
	}

	if s.capacity > 0 && len(s.items) >= s.capacity {
		s.evict(now, rec)
	}

	e := &entry[K, V]{key: key, value: value, expiresAt: expiresAt}
	s.items[key] = e
	s.pushHead(e)
}

func (c *inner[K, V]) get(key K) (V, bool) {
	var zero V
	if !reflexive(key) {
		c.shards[0].stats.misses.Add(1)

		return zero, false
	}

	s := c.shardFor(key)
	s.mu.RLock()
	e, ok := s.items[key]
	if ok && c.live(e) {
		e.visited.Store(true)
		v := e.value
		s.mu.RUnlock()
		s.stats.hits.Add(1)

		return v, true
	}
	s.mu.RUnlock()

	if !ok {
		s.stats.misses.Add(1)

		return zero, false
	}

	// Expired: re-fetch under the write lock; it may have been replaced.
	if c.hooks.afterReadUnlock != nil {
		c.hooks.afterReadUnlock()
	}

	var rec []removal[K, V]
	s.mu.Lock()
	now := c.now()
	if e, ok = s.items[key]; ok && !e.expired(now) {
		e.visited.Store(true)
		v := e.value
		s.mu.Unlock()
		s.stats.hits.Add(1)

		return v, true
	}
	if ok {
		s.remove(e, ReasonExpired, c.recorder(&rec))
	}
	s.mu.Unlock()
	s.stats.misses.Add(1)
	c.dispatch(rec)

	return zero, false
}

func (c *inner[K, V]) peek(key K) (V, bool) {
	var zero V
	if !reflexive(key) {
		return zero, false
	}

	s := c.shardFor(key)
	s.mu.RLock()
	defer s.mu.RUnlock()

	if e, ok := s.items[key]; ok && c.live(e) {
		return e.value, true
	}

	return zero, false
}

func (c *inner[K, V]) setWithTTL(key K, value V, ttl time.Duration) {
	if !reflexive(key) {
		return
	}

	s := c.shardFor(key)
	var rec []removal[K, V]
	s.mu.Lock()
	s.invalidate(key)
	now := c.now()
	s.set(key, value, c.expiresAt(ttl, now), now, c.recorder(&rec))
	s.mu.Unlock()
	c.dispatch(rec)
}

func (c *inner[K, V]) delete(key K) bool {
	if !reflexive(key) {
		return false
	}

	s := c.shardFor(key)
	var rec []removal[K, V]
	live := false
	s.mu.Lock()
	s.invalidate(key)
	if e, ok := s.items[key]; ok {
		if e.expired(c.now()) {
			s.remove(e, ReasonExpired, c.recorder(&rec))
		} else {
			s.remove(e, ReasonDeleted, c.recorder(&rec))
			live = true
		}
	}
	s.mu.Unlock()
	c.dispatch(rec)

	return live
}

func (c *inner[K, V]) clear() {
	processed := 0
	defer func() {
		// A callback panic skips callbacks, but not the remaining shards.
		for _, s := range c.shards[processed:] {
			c.clearShard(s, nil)
		}
	}()

	for _, s := range c.shards {
		var rec []removal[K, V]
		c.clearShard(s, c.recorder(&rec))
		processed++ // The shard is complete before its callbacks run.
		c.dispatch(rec)
	}
}

func (c *inner[K, V]) clearShard(s *shard[K, V], rec *[]removal[K, V]) {
	s.mu.Lock()
	for key, cl := range s.inflight {
		cl.invalidated = true
		delete(s.inflight, key)
	}
	for e := s.tail; e != nil; {
		next := e.next
		s.remove(e, ReasonDeleted, rec)
		e = next
	}
	s.mu.Unlock()
}

func (c *inner[K, V]) length() int {
	n := 0
	for _, s := range c.shards {
		s.mu.RLock()
		n += len(s.items)
		s.mu.RUnlock()
	}

	return n
}

func (c *inner[K, V]) ttlOf(key K) (time.Duration, bool) {
	if !reflexive(key) {
		return 0, false
	}

	s := c.shardFor(key)
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.items[key]
	if !ok {
		return 0, false
	}
	if e.expiresAt == 0 {
		return NoExpiration, true
	}

	now := c.now()
	if e.expired(now) {
		return 0, false
	}

	return time.Duration(e.expiresAt - now), true
}

func (c *inner[K, V]) deleteExpired() int {
	processed := 0
	defer func() {
		// A callback panic skips callbacks, but not the remaining shards.
		for _, s := range c.shards[processed:] {
			c.deleteExpiredShard(s, nil)
		}
	}()

	removed := 0
	for _, s := range c.shards {
		var rec []removal[K, V]
		removed += c.deleteExpiredShard(s, c.recorder(&rec))
		processed++ // The shard is complete before its callbacks run.
		c.dispatch(rec)
	}

	return removed
}

func (c *inner[K, V]) deleteExpiredShard(s *shard[K, V], rec *[]removal[K, V]) int {
	removed := 0
	s.mu.Lock()
	now := c.now()
	for e := s.tail; e != nil; {
		next := e.next
		if e.expired(now) {
			s.remove(e, ReasonExpired, rec)
			removed++
		}
		e = next
	}
	s.mu.Unlock()

	return removed
}

func (c *inner[K, V]) all() iter.Seq2[K, V] {
	type pair struct {
		key   K
		value V
	}

	return func(yield func(K, V) bool) {
		var buf []pair
		for _, s := range c.shards {
			buf = buf[:0]
			s.mu.RLock()
			now := c.now()
			for e := s.tail; e != nil; e = e.next {
				if !e.expired(now) {
					buf = append(buf, pair{e.key, e.value})
				}
			}
			s.mu.RUnlock()

			for _, p := range buf {
				if !yield(p.key, p.value) {
					return
				}
			}
		}
	}
}
