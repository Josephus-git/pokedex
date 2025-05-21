package cache

import (
	"sync"
	"time"
)

// Cache holds a map consisting of the url keyed to its value
type Cache struct {
	MapCache         map[string]cacheEntry
	mu               sync.RWMutex
	timeToExpiration time.Duration
}

// cacheEntry holds the cache value and its created time
type cacheEntry struct {
	createdAt time.Time
	Val       []byte
}

// NewCache creates a new Cache with the given time to expiration
func NewCache(interval time.Duration) *Cache {
	newC := &Cache{
		MapCache:         make(map[string]cacheEntry),
		timeToExpiration: interval,
	}
	go newC.reapLoop()
	return newC
}

// adds a key, value pair into a cache map and also the time the cache was created
func (C *Cache) Add(key string, val []byte) {
	C.mu.Lock()
	defer C.mu.Unlock()
	newCVal := cacheEntry{
		createdAt: time.Now(),
		Val:       val,
	}
	C.MapCache[key] = newCVal
}

// gets the value with a key from the cache map and returns nil if not passed
func (C *Cache) Get(key string) ([]byte, bool) {
	C.mu.RLock()
	defer C.mu.RUnlock()
	item, found := C.MapCache[key]
	if !found {
		return []byte{}, false
	}
	return item.Val, true
}

func (C *Cache) reapLoop() {
	ticker := time.NewTicker(C.timeToExpiration)
	defer ticker.Stop()
	// the below ticker.C indicates that we wait for a tick to be sent when the interval has passed before carrying out the function
	for range ticker.C {
		C.deleteExpired()
	}
}

func (C *Cache) deleteExpired() {
	for key, Cvalue := range C.MapCache {
		expiresAt := Cvalue.createdAt.Add(C.timeToExpiration)
		if time.Now().After(expiresAt) {
			//item has expired, remove it
			C.mu.Lock()
			defer C.mu.Unlock()
			delete(C.MapCache, key)
		}
	}
}
