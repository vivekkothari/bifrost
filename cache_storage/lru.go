package cache_storage

import (
	"container/list"
)

// LRUCache implements CacheStorageInterface using a least-recently-used strategy
type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	lruList  *list.List
}

type cacheEntry struct {
	queryIndex int
	response   string
}

// NewLRUCache initializes an LRU cache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		lruList:  list.New(),
	}
}

// SetResponse sets the response for a given query index in the LRU cache.
func (l *LRUCache) SetResponse(queryIndex int, response string) {
	if element, exists := l.cache[queryIndex]; exists {
		l.lruList.MoveToBack(element)
		element.Value.(*cacheEntry).response = response
	} else {
		if len(l.cache) >= l.capacity {
			// Remove the least recently used item
			evictElement := l.lruList.Front()
			if evictElement != nil {
				evictEntry := evictElement.Value.(*cacheEntry)
				delete(l.cache, evictEntry.queryIndex)
				l.lruList.Remove(evictElement)
			}
		}
		// Add a new entry
		entry := &cacheEntry{queryIndex: queryIndex, response: response}
		element := l.lruList.PushBack(entry)
		l.cache[queryIndex] = element
	}
}

// GetResponse retrieves the response for a given query index from the LRU cache.
func (l *LRUCache) GetResponse(queryIndex int) *string {
	if element, exists := l.cache[queryIndex]; exists {
		l.lruList.MoveToBack(element)
		return &element.Value.(*cacheEntry).response
	}
	return nil
}
