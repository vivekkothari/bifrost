package embedding

import "context"

// BaseEmbedding defines the interface for embedding models.
type BaseEmbedding interface {
	// GetEmbeddings generates embeddings for the given text.
	// text can be a single string or a slice of strings.
	// kwargs are optional additional parameters.
	GetEmbeddings(ctx context.Context, text interface{}, kwargs ...interface{}) ([]float64, error)

	// Dimension returns the size (number of dimensions) of the embeddings.
	Dimension(ctx context.Context) int
}
