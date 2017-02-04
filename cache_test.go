package internalCache

import (
	"testing"
	"time"
)

type testStruct struct {
	Name string
}

func TestSetAndGet(t *testing.T) {
	cache:= NewCache(1 * time.Second, 10 * time.Second)

	cache.Set("test", "TEST DATA")
	cache.Set("test", testStruct{"For example"})

	fromCache := cache.Get("test")
	if fromCache == nil {
		t.Error("Value is nil")
	}

	_, ok := fromCache.(testStruct)

	if !ok {
		t.Error("Can't convert data", fromCache)
	}

	cache.Destroy()
}

func TestGetAfterClean(t *testing.T) {
	cache:= NewCache(2 * time.Second, 10 * time.Second)

	cache.Set("test", "test")
	data := cache.Get("test")
	if data == nil {
		t.Error("Value is nil")
	}

	time.Sleep(3 * time.Second)

	data = cache.Get("test")
	if data != nil {
		t.Error("Value is not nil")
	}
}

func TestGetUnsetted(t *testing.T) {
	cache:= NewCache(1 * time.Second, 10 * time.Second)

	data := cache.Get("key")
	if data != nil {
		t.Error("Value is not nil")
	}
	cache.Destroy()
}

func TestDel(t *testing.T) {
	cache:= NewCache(1 * time.Second, 10 * time.Second)

	cache.Set("key", "DATA")
	cache.Del("key")

	data := cache.Get("key")

	if data != nil {
		t.Error("Value is not nil")
	}
	cache.Destroy()
}

func TestCleaner(t *testing.T)  {
	cache:= NewCache(2 * time.Second, 5 * time.Second)

	cache.Set("keyTest", "DATA")

	data := cache.Get("keyTest")
	if data == nil {
		t.Error("Value is nil")
	}

	timer := time.NewTimer(6 * time.Second)
	<-timer.C

	if len(cache.data) != 0 {
		t.Error("Cache was't cleaned")
	}

	timer.Stop()
	cache.Destroy()
}