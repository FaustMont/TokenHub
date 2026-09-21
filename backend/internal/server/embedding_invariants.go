package server

import (
	"encoding/json"
	"reflect"
)

// Only text may be rewritten by retrieval hooks. Vector shape and task
// semantics belong to the client and must survive every gateway stage.
func validateEmbeddingPatch(original, next EmbeddingsRequest) error {
	if err := validateEmbeddingRequest(next); err != nil {
		return NewHTTPError(502, "gateway_hook_patch_invalid", "Gateway plugin returned an invalid embedding request")
	}
	before, _, err := embeddingInputCount(original.Input)
	if err != nil {
		return err
	}
	after, _, err := embeddingInputCount(next.Input)
	if err != nil {
		return err
	}
	if original.Model != next.Model || before != after || !reflect.DeepEqual(original.Dimensions, next.Dimensions) ||
		original.EncodingFormat != next.EncodingFormat || original.InputType != next.InputType || original.Task != next.Task ||
		!reflect.DeepEqual(original.Normalized, next.Normalized) || !reflect.DeepEqual(original.Truncation, next.Truncation) || !reflect.DeepEqual(original.LateChunking, next.LateChunking) {
		return NewHTTPError(502, "gateway_hook_patch_invalid", "Gateway plugin cannot change embedding model, cardinality, dimensions, encoding or task semantics")
	}
	return nil
}

func normalizeEmbeddingResult(response any, request EmbeddingsRequest) (any, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return nil, NewHTTPError(502, "invalid_embedding_response", "Embedding response cannot be encoded")
	}
	var body map[string]any
	if json.Unmarshal(data, &body) != nil || body == nil {
		return nil, NewHTTPError(502, "invalid_embedding_response", "Expected an embedding response object")
	}
	return normalizeEmbeddingResponse(body, request)
}
