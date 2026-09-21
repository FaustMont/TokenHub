package server

import (
	"encoding/json"
	"strings"
)

// Unconfigured routes only fail over within the same provider and model. A
// shared explicit space is an administrator assertion, not a name heuristic.
func compatibleEmbeddingRoutes(routes []RouteSelection) []RouteSelection {
	if len(routes) < 2 {
		return routes
	}
	anchor := routes[0]
	space := embeddingRouteSpace(anchor)
	result := make([]RouteSelection, 0, len(routes))
	for _, route := range routes {
		candidate := embeddingRouteSpace(route)
		if (space != "" && candidate == space) || (space == "" && candidate == "" && route.Provider.ID == anchor.Provider.ID && route.ProviderModel == anchor.ProviderModel) {
			result = append(result, route)
		}
	}
	return result
}

func embeddingRouteSpace(route RouteSelection) string {
	var spaces map[string]string
	if json.Unmarshal([]byte(route.Provider.Options["embedding_spaces"]), &spaces) != nil {
		return ""
	}
	return strings.TrimSpace(spaces[route.ProviderModel])
}
