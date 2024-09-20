package embedding

import (
	"context"
	"fmt"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// OpenAIEmbeddings implements BaseEmbedding for OpenAI's embeddings API using openai-go.
type OpenAIEmbeddings struct {
	client                  *openai.Client
	modelName               openai.EmbeddingNewParamsModel
	dimension               int
	modelToDimensionMapping map[openai.EmbeddingNewParamsModel]int
}

// NewOpenAIEmbeddings creates a new instance of OpenAIEmbeddings.
func NewOpenAIEmbeddings(apiKey string, modelName openai.EmbeddingNewParamsModel) *OpenAIEmbeddings {
	return &OpenAIEmbeddings{
		client:    openai.NewClient(option.WithAPIKey(apiKey)),
		modelName: modelName,
		modelToDimensionMapping: map[openai.EmbeddingNewParamsModel]int{
			openai.EmbeddingNewParamsModelTextEmbedding3Large: 3072,
			openai.EmbeddingNewParamsModelTextEmbedding3Small: 1536,
			openai.EmbeddingNewParamsModelTextEmbeddingAda002: 1536,
		},
	}
}

// GetEmbeddings generates embeddings for the given text.
func (o *OpenAIEmbeddings) GetEmbeddings(ctx context.Context, text interface{}, kwargs ...interface{}) ([]float64, error) {
	var input string
	switch v := text.(type) {
	case string:
		input = v
	default:
		return nil, fmt.Errorf("invalid type for text parameter: %T", v)
	}

	resp, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.F(o.modelName),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](shared.UnionString(input)),
	})
	if err != nil {
		return nil, err
	}
	var embeddings []float64
	for _, item := range resp.Data {
		embeddings = append(embeddings, item.Embedding...)
	}
	return embeddings, nil
}

// Dimension returns the size (number of dimensions) of the embeddings.
func (o *OpenAIEmbeddings) Dimension(ctx context.Context) int {
	if o.dimension == 0 {
		if dim, ok := o.modelToDimensionMapping[o.modelName]; ok {
			o.dimension = dim
		} else {
			// Dynamically determine the dimension using a sample
			sampleEmbedding, err := o.GetEmbeddings(ctx, "sample")
			if err != nil {
				return 0
			}
			o.dimension = len(sampleEmbedding)
			o.modelToDimensionMapping[o.modelName] = o.dimension
		}
	}
	return o.dimension
}

// SetModel changes the model used for generating embeddings.
func (o *OpenAIEmbeddings) SetModel(modelName openai.EmbeddingNewParamsModel) {
	o.modelName = modelName
	o.dimension = 0 // Invalidate the cached dimension
}
