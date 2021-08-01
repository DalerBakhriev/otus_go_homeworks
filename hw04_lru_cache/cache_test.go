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

	t.Run("removing elements", func(t *testing.T) {
		c := NewCache(3)
		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)
		c.Set("d", 4)

		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		val, ok = c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		c.Set("e", 5)
		val, ok = c.Get("e")
		require.True(t, ok)
		require.Equal(t, 5, val)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)

		c.Set("f", 6)
		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 6, val)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("removing element after update", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 1)
		c.Set("b", 2)

		c.Set("b", 1)
		c.Set("a", 2)

		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 1, val)

		c.Set("c", 3)
		// Элемент "a" обновлялся последним поэтому должен вытолкнуться элемент "b"

		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("Wrong cache initialization (zero)", func(t *testing.T) {
		require.PanicsWithError(t, ErrCacheWithZeroOrNegativeCapacity.Error(), func() { NewCache(0) })
	})

	t.Run("Wrong cache initialization (negative)", func(t *testing.T) {
		require.PanicsWithError(t, ErrCacheWithZeroOrNegativeCapacity.Error(), func() { NewCache(-1) })
	})

	t.Run("Clearing cache", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 1)
		c.Set("b", 2)

		c.Clear()
		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
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
