package server

import "testing"

func TestRetrievalUsagePreservesPresenceBeforePricing(t *testing.T) {
	for _, tc := range []struct {
		name  string
		body  map[string]any
		known bool
	}{
		{"missing", map[string]any{}, false}, {"zero", map[string]any{"usage": map[string]any{"total_tokens": 0}}, true}, {"reported", map[string]any{"usage": map[string]any{"total_tokens": 5}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := retrievalUsage(tc.body, false)
			if (u.RetrievalEvidence.Quantity != nil) != tc.known {
				t.Fatal("lost presence")
			}
			charge := shadowPrice(nil, u, 0)
			if charge.Evidence == nil || (charge.Evidence.Quantity != nil) != tc.known {
				t.Fatal("missing price erased evidence")
			}
		})
	}
	u := retrievalUsage(map[string]any{"meta": map[string]any{"billed_units": map[string]any{"search_units": 2}}}, true)
	if u.TotalTokens != 0 || u.RetrievalEvidence.Unit != "search_unit" || *u.RetrievalEvidence.Quantity != 2 {
		t.Fatal("search units were converted to tokens")
	}
}
