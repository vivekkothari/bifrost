package cache_storage

import (
	"testing"
)

// Test basic set and get functionality of LFUCache
func TestLFUCache_SetGetResponse(t *testing.T) {
	cache := NewLFUCache(3)

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

// Test LFU eviction behavior
func TestLFUCache_Eviction(t *testing.T) {
	cache := NewLFUCache(3)

	// Add items to fill the cache
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Access Query 1 multiple times to increase its frequency
	cache.GetResponse(1) // Frequency of Query 1 is now 2
	cache.GetResponse(1) // Frequency of Query 1 is now 3

	// At this point:
	// Query 1 has frequency 3
	// Query 2 has frequency 1
	// Query 3 has frequency 1

	// Add two more items, causing the least frequently used (Query 2 and 3) to be evicted
	cache.SetResponse(4, "Response for Query 4")
	cache.SetResponse(5, "Response for Query 5")

	// Verify eviction of least frequently used query
	// Query 1 should still be in the cache (it was accessed the most)
	if response := cache.GetResponse(1); response == nil || *response != "Response for Query 1" {
		t.Errorf("Expected 'Response for Query 1' to still be present, got %v", response)
	}

	// One of Query 2 or Query 3 should have been evicted
	// Since both have the same frequency, LRU eviction should occur.
	if response := cache.GetResponse(2); response != nil {
		t.Errorf("Expected Query 2 to be evicted, but got %v", *response)
	}

	if response := cache.GetResponse(3); response != nil {
		t.Errorf("Expected Query 3 to be evicted, but got %v", *response)
	}

	// Query 4 should have been added successfully
	if response := cache.GetResponse(4); response == nil || *response != "Response for Query 4" {
		t.Errorf("Expected 'Response for Query 4', got %v", response)
	}
	// Query 5 should have been added successfully
	if response := cache.GetResponse(5); response == nil || *response != "Response for Query 5" {
		t.Errorf("Expected 'Response for Query 5', got %v", response)
	}
}

// Test that items with the lowest frequency are evicted first
func TestLFUCache_FrequencyEviction(t *testing.T) {
	cache := NewLFUCache(3)

	// Add items to fill the cache
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Access Query 1 twice and Query 2 once to adjust frequencies
	cache.GetResponse(1)
	cache.GetResponse(1)
	cache.GetResponse(2)

	// Add another item, causing Query 3 (least frequent) to be evicted
	cache.SetResponse(4, "Response for Query 4")

	// Check that Query 3 has been evicted
	if response := cache.GetResponse(3); response != nil {
		t.Errorf("Expected Query 3 to be evicted, but got %v", *response)
	}

	// Check that Query 1, Query 2, and Query 4 are still in the cache
	if response := cache.GetResponse(1); response == nil || *response != "Response for Query 1" {
		t.Errorf("Expected 'Response for Query 1' to still be present, got %v", response)
	}

	if response := cache.GetResponse(2); response == nil || *response != "Response for Query 2" {
		t.Errorf("Expected 'Response for Query 2' to still be present, got %v", response)
	}

	if response := cache.GetResponse(4); response == nil || *response != "Response for Query 4" {
		t.Errorf("Expected 'Response for Query 4', got %v", response)
	}
}

// Test that items with equal frequency follow least recently used eviction
func TestLFUCache_LeastRecentlyUsedEviction(t *testing.T) {
	cache := NewLFUCache(3)

	// Add items to fill the cache
	cache.SetResponse(1, "Response for Query 1")
	cache.SetResponse(2, "Response for Query 2")
	cache.SetResponse(3, "Response for Query 3")

	// Access Query 1 and Query 2 to increase their frequency
	cache.GetResponse(1)
	cache.GetResponse(2)

	// Add another item, causing Query 3 (least recently used and least frequent) to be evicted
	cache.SetResponse(4, "Response for Query 4")

	// Check that Query 3 has been evicted
	if response := cache.GetResponse(3); response != nil {
		t.Errorf("Expected Query 3 to be evicted, but got %v", *response)
	}

	// Check that Query 1, Query 2, and Query 4 are still in the cache
	if response := cache.GetResponse(1); response == nil || *response != "Response for Query 1" {
		t.Errorf("Expected 'Response for Query 1' to still be present, got %v", response)
	}

	if response := cache.GetResponse(2); response == nil || *response != "Response for Query 2" {
		t.Errorf("Expected 'Response for Query 2' to still be present, got %v", response)
	}

	if response := cache.GetResponse(4); response == nil || *response != "Response for Query 4" {
		t.Errorf("Expected 'Response for Query 4', got %v", response)
	}
}
