package cache_storage

// CacheStorageInterface defines the interface for the cache storage
type CacheStorageInterface interface {
	SetResponse(queryIndex int, response string)
	GetResponse(queryIndex int) *string
}
