package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {

	type testStruct struct {
		Name string
	}

	cache := New(1*time.Second, 10*time.Second)
	defer cache.Close()

	cache.Set("test", "TEST DATA")
	cache.Set("test", testStruct{"For example"})

	fromCache := cache.Get("test")
	if fromCache == nil {
		t.Error("Value is nil")
	}

	if _, ok := fromCache.(testStruct); !ok {
		t.Error("Can't convert data", fromCache)
	}
}

func TestGetAfterClean(t *testing.T) {
	cache := New(time.Second/2, time.Second)
	defer cache.Close()

	cache.Set("test", "test")

	if val := cache.Get("test"); val == nil {
		t.Error("Value is nil")
	}

	time.Sleep(cache.Duration + 100)

	if val := cache.Get("test"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestGetUnsetted(t *testing.T) {
	cache := New(time.Second, 2*time.Second)
	defer cache.Close()

	if val := cache.Get("key"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestDel(t *testing.T) {
	cache := New(time.Second, 10*time.Second)
	defer cache.Close()

	cache.Set("key", "DATA")
	cache.Del("key")

	if val := cache.Get("key"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestCleaner(t *testing.T) {
	cache := New(2*time.Second, 2*time.Second)
	defer cache.Close()

	cache.Set("keyTest", "DATA")

	if val := cache.Get("keyTest"); val == nil {
		t.Error("Value is nil")
	}

	time.Sleep(cache.TickerDuration * 2)

	if len(cache.data) != 0 {
		t.Error("Cache was't cleaned")
	}
}

func BenchmarkThreads(b *testing.B) {

	cache := New(2*time.Second, 2*time.Second)
	defer cache.Close()

	runThread := func(threadNum int) {
		pref := threadNum * 10000

		for i := 0; i < 1000; i++ {
			key := strconv.Itoa(pref + i)

			operType := i % 3
			if operType == 0 {
				cache.Set(key, key)
			} else if operType == 1 {
				cache.Get(key)
			} else if operType == 2 {
				cache.Del(key)
			}
		}
	}

	wg := sync.WaitGroup{}
	for thread := 0; thread < 5000; thread++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			runThread(thread)
		}()
	}

	wg.Wait()
}
