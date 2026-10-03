package cache

import (
	"runtime"
	"time"
)

// requestStop signals the janitor to exit. It never waits and is safe to call
// any number of times; it is the only code that closes stop.
func (c *inner[K, V]) requestStop() {
	c.stopOnce.Do(func() { close(c.stop) })
}

// Close stops the janitor, if any, and waits for it to exit. It does not
// cancel in-flight loads. Close is idempotent and safe for concurrent use;
// the cache remains usable afterwards. It always returns nil.
func (c *Cache[K, V]) Close() error {
	c.in.requestStop()
	if c.in.janitorDone != nil {
		<-c.in.janitorDone
	}

	if c.hasCleanup {
		c.cleanupOnce.Do(c.cleanup.Stop)
	}

	runtime.KeepAlive(c)

	return nil
}

// startJanitor starts the goroutine that removes expired entries every
// interval until stop is closed. It references only the inner state.
func (c *inner[K, V]) startJanitor(interval time.Duration) {
	c.janitorDone = make(chan struct{})

	go func() {
		defer close(c.janitorDone)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-c.stop:
				return
			case <-ticker.C:
				c.deleteExpired()
			}
		}
	}()
}
