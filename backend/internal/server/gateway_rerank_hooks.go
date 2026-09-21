package server

import (
	"context"
	"encoding/json"
	"net/http"
)

func rerankPatch(r *RerankRequest) func(json.RawMessage) error {
	return func(data json.RawMessage) error { return applyRerankPatch(r, data) }
}
func (s *Server) runGatewayRerankDecodeNormalizeHooks(ctx context.Context, call CallContext, headers http.Header, r *RerankRequest) error {
	return s.runGatewayDecodeNormalizeHooks(ctx, call, headers, *r, rerankPatch(r))
}
func (s *Server) runGatewayRerankPrivacyPreHooks(ctx context.Context, call CallContext, headers http.Header, r *RerankRequest) error {
	return s.runGatewayPrivacyPreHooks(ctx, call, headers, *r, rerankPatch(r))
}
func (s *Server) runGatewayRerankGuardrailPreHooks(ctx context.Context, call CallContext, r *RerankRequest) error {
	return s.runGatewayGuardrailPreHooks(ctx, call, *r, rerankGuardrailTargets(r), rerankPatch(r))
}
func (s *Server) runGatewayRerankContextOptimizeHooks(ctx context.Context, call CallContext, r *RerankRequest) error {
	return s.runGatewayContextOptimizeHooks(ctx, call, *r, rerankPatch(r))
}
func (s *Server) executeRoutedRerank(r *http.Request, routed RoutedCall, req RerankRequest) (any, RouteSelection, Usage, []RouteAttempt, error) {
	return executeRoutedWithStore(r.Context(), s.store, routed, false, func(ctx context.Context, route RouteSelection, _ bool, _ int) (any, Usage, error) {
		route, err := s.prepareRouteForUpstream(ctx, route)
		if err != nil {
			return nil, Usage{}, err
		}
		upstream := req
		upstream.Documents = append([]string(nil), req.Documents...)
		if err = s.runGatewayRequestTransformHooks(ctx, routed.Call, route, upstream, providerRouteProtocolRerank, rerankPatch(&upstream)); err != nil {
			return nil, Usage{}, err
		}
		if resp, usage, handled, err := s.runGatewayProviderCallHooks(ctx, routed.Call, route, upstream, providerRouteProtocolRerank); err != nil || handled {
			return resp, usage, err
		}
		adapter, err := s.adapterForRoute(route)
		if err != nil {
			return nil, Usage{}, err
		}
		reranker, ok := adapter.(ProviderReranker)
		if !ok {
			return nil, Usage{}, NewHTTPError(501, "provider_capability_not_supported", "Provider does not support rerank")
		}
		return reranker.Rerank(ctx, route.Provider, route.ProviderModel, upstream)
	})
}
