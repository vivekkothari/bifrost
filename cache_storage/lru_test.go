package cache_storage

import (
	"testing"
)

// Test that the LRUCache correctly stores and retrieves responses
func TestLRUCache_SetGetResponse(t *testing.T) {
	cache := NewLRUCache(3)

	// Test adding and retrieving items
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Get the responses and check if they are correct
	if response := cache.GetResponse(1); response == nil || *response != "Response for Query 1" {
		t.Errorf("Expected 'Response for Query 1', got %v", response)
	}

	if response := cache.GetResponse(2); response == nil || *response != "Response for Query 2" {
		t.Errorf("Expected 'Response for Query 2', got %v", response)
	}

	if response := cache.GetResponse(3); response == nil || *response != "Response for Query 3" {
		t.Errorf("Expected 'Response for Query 3', got %v", response)
	}
}

// Test LRU eviction when the cache reaches capacity
func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(3)

	// Add items to fill the cache
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Add another item, causing the least recently used (Query 1) to be evicted
	cache.SetResponse(4, "Response for Query 4")

	// Check that Query 1 has been evicted and Query 4 is added
	if response := cache.GetResponse(1); response != nil {
		t.Errorf("Expected Query 1 to be evicted, but got %v", *response)
	}

	if response := cache.GetResponse(4); response == nil || *response != "Response for Query 4" {
		t.Errorf("Expected 'Response for Query 4', got %v", response)
	}
}

// Test that recently used items are not evicted
func TestLRUCache_RecentlyUsed(t *testing.T) {
	cache := NewLRUCache(3)

	// Add items to fill the cache
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Access Query 1 to make it recently used
	cache.GetResponse(1)

	// Add another item, causing the least recently used (Query 2) to be evicted
	cache.SetResponse(4, "Response for Query 4")

	// Check that Query 2 has been evicted and Query 1 is still present
	if response := cache.GetResponse(2); response != nil {
		t.Errorf("Expected Query 2 to be evicted, but got %v", *response)
	}

	if response := cache.GetResponse(1); response == nil || *response != "Response for Query 1" {
		t.Errorf("Expected 'Response for Query 1' to be present, but got %v", response)
	}
}
