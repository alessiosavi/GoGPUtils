package cache_test

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/alessiosavi/GoGPUtils/cache"
)

func Example() {
	c, err := cache.New(cache.Config[string, int]{MaxEntries: 2})
	if err != nil {
		fmt.Println(err)

		return
	}
	c.Set("a", 1)
	c.Set("b", 2)
	v, ok := c.Get("a")
	fmt.Println(v, ok, c.Len())
	// Output: 1 true 2
}

func ExampleCache_GetOrLoad() {
	c, err := cache.New(cache.Config[int, string]{TTL: time.Minute})
	if err != nil {
		fmt.Println(err)

		return
	}
	load := func(context.Context) (string, error) { return "user-42", nil }
	for range 2 {
		v, err := c.GetOrLoad(context.Background(), 42, load)
		fmt.Println(v, err)
	}
	st := c.Stats()
	fmt.Println(st.Loads, st.Hits, st.Misses)
	// Output:
	// user-42 <nil>
	// user-42 <nil>
	// 1 1 1
}

func ExampleCache_All() {
	c, err := cache.New(cache.Config[string, int]{})
	if err != nil {
		fmt.Println(err)

		return
	}
	c.Set("b", 2)
	c.Set("a", 1)
	var keys []string
	for k := range c.All() {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	fmt.Println(keys)
	// Output: [a b]
}

func ExampleConfig_onEvict() {
	c, err := cache.New(cache.Config[string, int]{
		MaxEntries: 1,
		OnEvict: func(key string, value int, reason cache.EvictionReason) {
			fmt.Println(key, value, reason)
		},
	})
	if err != nil {
		fmt.Println(err)

		return
	}
	c.Set("a", 1)
	c.Set("b", 2)
	c.Delete("b")
	// Output:
	// a 1 evicted
	// b 2 deleted
}
