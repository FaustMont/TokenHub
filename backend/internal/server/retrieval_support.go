package server

import (
	"net/http"
	"strings"
	"tokenhub/backend/internal/metering"
)

func validRerankProtocol(profile string) bool {
	switch profile {
	case "jina", "cohere", "voyage", "qwen", "dashscope", "tei":
		return true
	}
	return false
}
func (s *Server) providerRetrievalSupport(provider Provider, modality string) bool {
	descriptor, ok := s.adapterRegistry.Describe(provider.Type)
	if !ok {
		return false
	}
	switch modality {
	case "embedding":
		if !adapterSupports(descriptor, AdapterCapabilityEmbeddings) {
			return false
		}
		if provider.Type == ProviderOpenAICompatible {
			switch providerEmbeddingProtocol(provider) {
			case "openai", "cohere", "jina", "voyage", "dashscope", "tei":
				return true
			default:
				return false
			}
		}
		return true
	case "rerank":
		return adapterSupports(descriptor, AdapterCapabilityRerank) && (provider.Type != ProviderOpenAICompatible || validRerankProtocol(providerRerankProtocol(provider)))
	case "video", "audio", "ocr":
		return false
	default:
		return true
	}
}
func retrievalPriceConfigured(model Model, search bool, provider bool) bool {
	if search {
		_, err := metering.Decimal(model.Metadata[retrievalSearchUnitPriceKey])
		return err == nil
	}
	price := model.InputPriceUSDPer1M
	if !provider && model.Modality == "embedding" {
		price = model.EmbeddingPriceUSDPer1M
	}
	return price > 0 || model.Metadata["retrieval_pricing_confirmed"] == "true"
}
func (s *Server) validateRetrievalRoute(route ModelRoute, pending *Model, provider Provider) error {
	// Preserve already-published routes during unrelated edits. Their runtime
	// capability checks still apply; no upgrade silently disables old traffic.
	for _, old := range s.store.ListRoutes() {
		if old.ID == route.ID && old.ModelName == route.ModelName && old.ProviderID == route.ProviderID && old.ProviderModel == route.ProviderModel && old.Status == route.Status {
			return nil
		}
	}
	var model Model
	found := false
	if pending != nil {
		model = *pending
		found = true
	} else {
		for _, candidate := range s.store.ListModels() {
			if candidate.Name == route.ModelName {
				model = candidate
				found = true
				break
			}
		}
	}
	if !found {
		return nil
	}
	if !s.providerRetrievalSupport(provider, model.Modality) {
		return NewHTTPError(400, "model_operation_unsupported", "Provider cannot execute this model operation; retain it in the catalog without publishing a route")
	}
	if model.Modality != "embedding" && model.Modality != "rerank" {
		return nil
	}
	if provider.Type == ProviderMock {
		return nil
	}
	search := model.Modality == "rerank" && providerRerankProtocol(provider) == "cohere"
	if !retrievalPriceConfigured(model, search, false) {
		return NewHTTPError(400, "retrieval_price_required", "Configure the tenant price or explicitly confirm a free token price before publishing")
	}
	for _, upstream := range s.store.ListProviderModels() {
		if upstream.ProviderID == provider.ID && upstream.UpstreamModel == route.ProviderModel {
			if !retrievalTextInputSupported(upstream) || (upstream.Modality != "" && upstream.Modality != model.Modality) {
				return NewHTTPError(400, "model_operation_mismatch", "Upstream and public model operations differ")
			}
			if !retrievalPriceConfigured(providerModelCostModel(upstream), search, true) {
				return NewHTTPError(400, "retrieval_price_required", "Configure the provider price or explicitly confirm a free token price before publishing")
			}
			return nil
		}
	}
	return NewHTTPError(400, "provider_model_required", "Import and configure the upstream model before publishing")
}
func (s *Server) pricedRerankRoutes(model Model, routes []RouteSelection) []RouteSelection {
	result := make([]RouteSelection, 0, len(routes))
	models := s.store.ListProviderModels()
	for _, route := range routes {
		if !s.providerRetrievalSupport(route.Provider, "rerank") {
			continue
		}
		search := providerRerankProtocol(route.Provider) == "cohere"
		if !retrievalPriceConfigured(model, search, false) {
			continue
		}
		for _, upstream := range models {
			if upstream.ProviderID == route.Provider.ID && upstream.UpstreamModel == route.ProviderModel && retrievalPriceConfigured(providerModelCostModel(upstream), search, true) {
				result = append(result, route)
				break
			}
		}
	}
	return result
}
func (s *Server) handleAdminRerankTest(w http.ResponseWriter, r *http.Request) {
	user, ok := s.requireAdmin(w, r, "provider", r.Method)
	if !ok {
		return
	}
	var req struct {
		ProviderID string        `json:"provider_id"`
		ResourceID string        `json:"resource_id,omitempty"`
		Request    RerankRequest `json:"request"`
	}
	if err := s.decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if err := validateRerankRequest(req.Request); err != nil {
		writeError(w, r, err)
		return
	}
	provider, ok := s.providerByID(strings.TrimSpace(req.ProviderID))
	if !ok {
		writeError(w, r, NewHTTPError(404, "provider_not_found", "Provider not found"))
		return
	}
	selection := RouteSelection{Provider: provider, ProviderModel: req.Request.Model}
	if req.ResourceID != "" {
		resource, ok := s.providerResourceByID(req.ResourceID)
		if !ok || resource.ProviderID != provider.ID {
			writeError(w, r, NewHTTPError(400, "route_resource_mismatch", "Resource must belong to Provider"))
			return
		}
		selection.Resource = &resource
	}
	selection, err := s.prepareRouteForUpstream(r.Context(), selection)
	if err != nil {
		writeError(w, r, err)
		return
	}
	adapter, err := s.adapterForRoute(selection)
	if err != nil {
		writeError(w, r, err)
		return
	}
	reranker, ok := adapter.(ProviderReranker)
	if !ok {
		writeError(w, r, NewHTTPError(501, "provider_capability_not_supported", "Provider does not support rerank"))
		return
	}
	response, usage, err := reranker.Rerank(r.Context(), selection.Provider, selection.ProviderModel, req.Request)
	s.recordAdminAudit(r, user, "test", "provider_rerank", provider.ID, "", map[string]any{"model": req.Request.Model, "usage_evidence": usage.RetrievalEvidence, "success": err == nil})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"response": response, "usage_evidence": usage.RetrievalEvidence, "pricing_status": "unverified"})
}

func validateRetrievalPriceMetadata(metadata map[string]string) error {
	if value := metadata[retrievalSearchUnitPriceKey]; value != "" {
		if _, err := metering.Decimal(value); err != nil {
			return NewHTTPError(400, "invalid_retrieval_price", "search_unit_price_usd must be a non-negative decimal")
		}
	}
	return nil
}

func retrievalTextInputSupported(model ProviderModel) bool {
	if model.Modality != "embedding" && model.Modality != "rerank" {
		return true
	}
	if len(model.InputModalities) == 0 {
		return true
	}
	for _, modality := range model.InputModalities {
		if modality == "text" {
			return true
		}
	}
	return false
}
