package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic capacity", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		// "aaa" should not exist, first added
		_, ok := c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("purge logic lru", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		_, ok := c.Get("bbb")
		require.True(t, ok)
		wasInCache = c.Set("bbb", 1000)
		require.True(t, wasInCache)
		for i := 0; i < 10; i++ {
			_, ok = c.Get("aaa")
			require.True(t, ok)
			if i%2 == 0 {
				_, ok = c.Get("ccc")
				require.True(t, ok)
			}
		}
		wasInCache = c.Set("aaa", 500)
		require.True(t, wasInCache)
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		// "ccc" should not exist, least recently used
		_, ok = c.Get("ccc")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
