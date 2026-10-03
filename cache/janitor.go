package cache

import "runtime"

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
