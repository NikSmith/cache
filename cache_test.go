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
	cache := New(2*time.Second, 10*time.Second)
	defer cache.Close()

	cache.Set("test", "test")

	if val := cache.Get("test"); val == nil {
		t.Error("Value is nil")
	}

	time.Sleep(3 * time.Second)

	if val := cache.Get("test"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestGetUnsetted(t *testing.T) {
	cache := New(1*time.Second, 10*time.Second)
	defer cache.Close()

	if val := cache.Get("key"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestDel(t *testing.T) {
	cache := New(1*time.Second, 10*time.Second)
	defer cache.Close()

	cache.Set("key", "DATA")
	cache.Del("key")

	if val := cache.Get("key"); val != nil {
		t.Error("Value is not nil")
	}
}

func TestCleaner(t *testing.T) {
	cache := New(2*time.Second, 5*time.Second)
	defer cache.Close()

	cache.Set("keyTest", "DATA")

	if val := cache.Get("keyTest"); val == nil {
		t.Error("Value is nil")
	}

	timer := time.NewTimer(6 * time.Second)
	<-timer.C

	if len(cache.data) != 0 {
		t.Error("Cache was't cleaned")
	}

	timer.Stop()
}

func BenchmarkThreads(b *testing.B) {

	cache := New(2*time.Second, 2*time.Second)
	defer cache.Close()

	add := func(threadNum int) {
		for i := 0; i < 1000; i++ {
			key := strconv.Itoa(threadNum*10000 + i)
			cache.Set(key, key)
		}
	}

	get := func(threadNum int) {
		for i := 0; i < 1000; i++ {
			key := strconv.Itoa(threadNum*10000 + i)
			cache.Get(key)
		}
	}

	del := func(threadNum int) {
		for i := 0; i < 1000; i++ {
			key := strconv.Itoa(threadNum*10000 + i)
			cache.Del(key)
		}
	}

	wg := sync.WaitGroup{}
	for thread := 0; thread < 5000; thread++ {

		go func() {
			wg.Add(1)
			defer wg.Done()

			add(thread)
		}()

		go func() {
			wg.Add(1)
			defer wg.Done()

			get(thread)
		}()

		go func() {
			wg.Add(1)
			defer wg.Done()

			del(thread)
		}()
	}

	wg.Wait()
}
