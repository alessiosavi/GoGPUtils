package cache

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNew_RejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name  string
		cfg   Config[string, int]
		field string
	}{
		{"negative MaxEntries", Config[string, int]{MaxEntries: -1}, "MaxEntries"},
		{"negative TTL", Config[string, int]{TTL: -time.Second}, "TTL"},
		{"negative Shards", Config[string, int]{Shards: -1}, "Shards"},
		{"too many Shards", Config[string, int]{Shards: MaxShards + 1}, "Shards"},
		{"negative JanitorInterval", Config[string, int]{JanitorInterval: -time.Second}, "JanitorInterval"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.cfg)
			if c != nil {
				t.Fatalf("New returned a cache for an invalid config")
			}
			if !errors.Is(err, ErrInvalidConfig) || !strings.Contains(err.Error(), tt.field) {
				t.Fatalf("New error = %v, want ErrInvalidConfig naming %s", err, tt.field)
			}
		})
	}
}

func TestNew_AcceptsBoundaryConfig(t *testing.T) {
	tests := []struct {
		name       string
		cfg        Config[string, int]
		wantShards int
	}{
		{"single shard", Config[string, int]{Shards: 1}, 1},
		{"MaxShards", Config[string, int]{Shards: MaxShards}, MaxShards},
		{"fewer entries than shards", Config[string, int]{MaxEntries: 1, Shards: 64}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.cfg)
			if err != nil || c == nil {
				t.Fatalf("New = %v, %v", c, err)
			}
			if got := len(c.in.shards); got != tt.wantShards {
				t.Fatalf("shards = %d, want %d", got, tt.wantShards)
			}
			if c.in.janitorDone != nil {
				t.Fatal("New started a janitor without JanitorInterval")
			}
		})
	}
}

func TestShardCount(t *testing.T) {
	tests := []struct {
		maxEntries, shards, procs, want int
	}{
		{0, 0, 4, 16},
		{0, 0, 1, 4},
		{0, 0, 1 << 20, MaxShards},
		{2047, 0, 8, 1},
		{2048, 0, 8, 2},
		{1 << 20, 0, 8, 32},
		{0, 3, 8, 4},
		{10, 16, 8, 8},
		{1, 64, 8, 1},
		{5, 1, 8, 1},
		{0, MaxShards, 8, MaxShards},
	}
	for _, tt := range tests {
		if got := shardCount(tt.maxEntries, tt.shards, tt.procs); got != tt.want {
			t.Errorf("shardCount(%d, %d, %d) = %d, want %d", tt.maxEntries, tt.shards, tt.procs, got, tt.want)
		}
	}
}

func TestSplitCapacity(t *testing.T) {
	tests := []struct {
		maxEntries, n int
		want          []int
	}{
		{10, 4, []int{3, 3, 2, 2}},
		{0, 4, []int{0, 0, 0, 0}},
		{5, 1, []int{5}},
		{8, 8, []int{1, 1, 1, 1, 1, 1, 1, 1}},
	}
	for _, tt := range tests {
		if got := splitCapacity(tt.maxEntries, tt.n); !slices.Equal(got, tt.want) {
			t.Errorf("splitCapacity(%d, %d) = %v, want %v", tt.maxEntries, tt.n, got, tt.want)
		}
	}
}

func TestEvictionReason_String(t *testing.T) {
	tests := map[EvictionReason]string{
		ReasonEvicted:     "evicted",
		ReasonExpired:     "expired",
		ReasonDeleted:     "deleted",
		ReasonReplaced:    "replaced",
		EvictionReason(0): "EvictionReason(0)",
	}
	for reason, want := range tests {
		if got := reason.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", int(reason), got, want)
		}
	}
}

func TestStats_HitRatio(t *testing.T) {
	if got := (Stats{Hits: 1 << 63, Misses: 1 << 63}).HitRatio(); got != 0.5 {
		t.Fatalf("large-counter HitRatio = %v, want 0.5", got)
	}
	if got := (Stats{}).HitRatio(); got != 0 {
		t.Errorf("zero Stats HitRatio = %v, want 0", got)
	}
	if got := (Stats{Hits: 3, Misses: 1}).HitRatio(); got != 0.75 {
		t.Errorf("HitRatio = %v, want 0.75", got)
	}
}

func TestPanicError(t *testing.T) {
	boom := errors.New("boom")
	pe := &PanicError{Value: boom}
	if !errors.Is(pe, boom) {
		t.Error("PanicError does not unwrap an error panic value")
	}
	if got, want := pe.Error(), "cache: loader panicked: boom"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if (&PanicError{Value: 42}).Unwrap() != nil {
		t.Error("Unwrap of a non-error panic value should be nil")
	}
}
