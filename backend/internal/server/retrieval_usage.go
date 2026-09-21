package server

import (
	"encoding/json"
	"math"
)

// RetrievalUsageEvidence stores only metering facts, never request documents.
// A nil quantity means absent; a pointer to zero is an observed zero.
type RetrievalUsageEvidence struct {
	Unit     string `json:"unit"`
	Quantity *int64 `json:"quantity,omitempty"`
	Source   string `json:"source"`
}

func retrievalUsage(body map[string]any, searchUnits bool) Usage {
	usage := usageFromMap(body)
	evidence := &RetrievalUsageEvidence{Unit: "token", Source: "unreported"}
	var value any
	if searchUnits {
		evidence.Unit = "search_unit"
		meta, _ := body["meta"].(map[string]any)
		billed, _ := meta["billed_units"].(map[string]any)
		value = billed["search_units"]
	} else {
		raw, _ := body["usage"].(map[string]any)
		value = firstNonNil(raw["total_tokens"], raw["prompt_tokens"], raw["input_tokens"])
	}
	if value != nil {
		encoded, err := json.Marshal(value)
		var count float64
		if err != nil || json.Unmarshal(encoded, &count) != nil || count < 0 || math.Trunc(count) != count || count >= float64(math.MaxInt64) {
			usage.MeteringInvalid = true
			evidence.Source = "invalid"
		} else {
			n := int64(count)
			evidence.Quantity = &n
			evidence.Source = "upstream"
			if !searchUnits && usage.PromptTokens == 0 {
				usage.PromptTokens = n
				usage.TotalTokens = n
			}
		}
	}
	usage.RetrievalEvidence = evidence
	return usage
}
