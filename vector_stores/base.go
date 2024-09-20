package vector_stores

// VectorStoreInterface defines the interface for vector store operations.
type VectorStoreInterface interface {
	Add(embedding []float64, options map[string]interface{}) string
	Search(embedding []float64, topN int, includeDistances bool, options map[string]interface{}) ([]string, []float64)
}
