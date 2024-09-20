package cache_storage

import (
	"container/list"
)

type LFUCache struct {
	capacity int
	cache    map[int]*cacheItem
	freq     map[int]*list.List
	minFreq  int
}

type cacheItem struct {
	queryIndex int
	response   string
	frequency  int
}

// NewLFUCache initializes the LFU cache with a specified capacity.
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[int]*cacheItem),
		freq:     make(map[int]*list.List),
		minFreq:  0,
	}
}

// SetResponse sets the response for a query index in the LFU cache.
func (l *LFUCache) SetResponse(queryIndex int, response string) {
	if l.capacity == 0 {
		return
	}

	if item, exists := l.cache[queryIndex]; exists {
		// Update existing item
		l.removeFromFreqList(item)
		item.response = response
		item.frequency++
		l.addToFreqList(item)
	} else {
		if len(l.cache) >= l.capacity {
			// Evict the least frequently used item
			l.evict()
		}
		// Add a new item
		item := &cacheItem{queryIndex: queryIndex, response: response, frequency: 1}
		l.cache[queryIndex] = item
		l.addToFreqList(item)
		l.minFreq = 1
	}
}

// GetResponse gets the response for a query index from the LFU cache.
func (l *LFUCache) GetResponse(queryIndex int) *string {
	// Debug output
	if item, exists := l.cache[queryIndex]; exists {
		l.removeFromFreqList(item)
		item.frequency++
		l.addToFreqList(item)
		return &item.response
	}
	return nil
}

func (l *LFUCache) removeFromFreqList(item *cacheItem) {
	freqList := l.freq[item.frequency]
	for e := freqList.Front(); e != nil; e = e.Next() {
		if e.Value.(*cacheItem).queryIndex == item.queryIndex {
			freqList.Remove(e)
			break
		}
	}
	if freqList.Len() == 0 {
		delete(l.freq, item.frequency)
		if l.minFreq == item.frequency {
			l.minFreq++
		}
	}
}

func (l *LFUCache) addToFreqList(item *cacheItem) {
	if l.freq[item.frequency] == nil {
		l.freq[item.frequency] = list.New()
	}
	l.freq[item.frequency].PushBack(item)
}

// Evicts the least frequently used item
func (l *LFUCache) evict() {
	freqList := l.freq[l.minFreq]
	if freqList == nil {
		return
	}
	evictItem := freqList.Remove(freqList.Front()).(*cacheItem)
	delete(l.cache, evictItem.queryIndex)
	if freqList.Len() == 0 {
		delete(l.freq, l.minFreq)
	}
}
