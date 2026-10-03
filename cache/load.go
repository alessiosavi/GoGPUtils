package cache

import (
	"context"
	"runtime/debug"
)

// call is one load generation for a key.
type call[V any] struct {
	done        chan struct{}
	val         V     // published before done is closed
	err         error // published before done is closed
	invalidated bool  // guarded by the owning shard's mu
}

// invalidate detaches key's in-flight call so a newer write wins and later
// callers start a new generation. The caller holds s.mu.
func (s *shard[K, V]) invalidate(key K) {
	if cl, ok := s.inflight[key]; ok {
		cl.invalidated = true
		delete(s.inflight, key)
	}
}

func (c *inner[K, V]) getOrLoad(ctx context.Context, key K, load func(context.Context) (V, error)) (V, error) {
	var zero V
	switch {
	case ctx == nil:
		return zero, ErrNilContext
	case load == nil:
		return zero, ErrNilLoader
	}
	if err := ctx.Err(); err != nil {
		return zero, err
	}

	if !reflexive(key) {
		// Never stored: run an ordinary, unregistered load.
		s := c.shards[0]
		s.stats.misses.Add(1)
		s.stats.loads.Add(1)
		cl := &call[V]{done: make(chan struct{})}
		go c.runLoad(context.WithoutCancel(ctx), key, s, cl, load, false)

		return c.wait(ctx, cl)
	}

	s := c.shardFor(key)
	s.mu.RLock()
	if e, ok := s.items[key]; ok && c.live(e) {
		e.visited.Store(true)
		v := e.value
		s.mu.RUnlock()
		s.stats.hits.Add(1)

		return v, nil
	}
	s.mu.RUnlock()

	if c.hooks.afterReadUnlock != nil {
		c.hooks.afterReadUnlock()
	}

	var rec []removal[K, V]
	s.mu.Lock()
	now := c.now()
	if e, ok := s.items[key]; ok {
		if !e.expired(now) {
			e.visited.Store(true)
			v := e.value
			s.mu.Unlock()
			s.stats.hits.Add(1)

			return v, nil
		}
		s.remove(e, ReasonExpired, c.recorder(&rec))
	}
	s.stats.misses.Add(1)

	cl, ok := s.inflight[key]
	if !ok {
		cl = &call[V]{done: make(chan struct{})}
		s.inflight[key] = cl
		s.stats.loads.Add(1)
		go c.runLoad(context.WithoutCancel(ctx), key, s, cl, load, true)
	}
	s.mu.Unlock()
	c.dispatch(rec)

	return c.wait(ctx, cl)
}

// wait returns the call's result, or ctx.Err() if ctx ends first. If both are
// ready, either may be returned.
func (c *inner[K, V]) wait(ctx context.Context, cl *call[V]) (V, error) {
	select {
	case <-cl.done:
		return cl.val, cl.err
	case <-ctx.Done():
		var zero V

		return zero, ctx.Err()
	}
}

// runLoad executes load and completes cl exactly once, even if the loader
// panics or calls runtime.Goexit.
func (c *inner[K, V]) runLoad(ctx context.Context, key K, s *shard[K, V], cl *call[V], load func(context.Context) (V, error), registered bool) {
	var (
		val      V
		err      error
		returned bool
	)
	defer func() {
		if !returned {
			var zero V
			val, err = zero, ErrLoaderGoexit
		}
		c.complete(key, s, cl, val, err, registered)
	}()

	val, err = invoke(ctx, load)
	returned = true
}

// invoke calls load, converting a panic into *PanicError.
func invoke[V any](ctx context.Context, load func(context.Context) (V, error)) (v V, err error) {
	defer func() {
		if r := recover(); r != nil {
			var zero V
			v, err = zero, &PanicError{Value: r, Stack: debug.Stack()}
		}
	}()

	return load(ctx)
}

// complete publishes the load result, stores it unless the call was
// invalidated, wakes waiters, and then runs OnEvict callbacks.
func (c *inner[K, V]) complete(key K, s *shard[K, V], cl *call[V], val V, err error, registered bool) {
	var rec []removal[K, V]
	s.mu.Lock()
	if registered && s.inflight[key] == cl {
		delete(s.inflight, key)
	}
	if err != nil {
		var zero V
		val = zero
		s.stats.loadErrors.Add(1)
	} else if registered && !cl.invalidated {
		now := c.now()
		s.set(key, val, c.expiresAt(0, now), now, c.recorder(&rec))
	}
	cl.val, cl.err = val, err
	s.mu.Unlock()

	close(cl.done)
	c.dispatch(rec)
}
