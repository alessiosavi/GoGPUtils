package cache

import "sync/atomic"

// Stats is a snapshot of cache counters. Counters are summed across shards
// and are not an atomic snapshot of the whole cache.
type Stats struct {
	// Hits and Misses count Get and GetOrLoad lookups.
	Hits, Misses uint64
	// Loads counts loader executions started; LoadErrors counts loads that
	// returned an error, panicked, or called runtime.Goexit.
	Loads, LoadErrors uint64
	// Evictions counts capacity evictions (ReasonEvicted only).
	Evictions uint64
	// Expirations counts removals with ReasonExpired.
	Expirations uint64
}

// HitRatio returns Hits/(Hits+Misses), or 0 when both are 0.
func (s Stats) HitRatio() float64 {
	total := float64(s.Hits) + float64(s.Misses)
	if total == 0 {
		return 0
	}

	return float64(s.Hits) / total
}

// shardStats holds per-shard counters, summed by Stats.
type shardStats struct {
	hits, misses, loads, loadErrors, evictions, expirations atomic.Uint64
}
