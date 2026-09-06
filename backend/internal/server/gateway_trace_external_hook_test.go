package server

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func TestExternalTraceHookWithoutIsolationIsSkippedByObserveOnlyPolicy(t *testing.T) {
	store := NewMemoryStore()
	app := NewWithConfig(store, Config{
		AdminToken: "external-trace-admin",
		PluginDir:  filepath.Join("..", "plugin", "testdata", "external-trace-hook"),
	})
	emitter := &recordingTraceEmitter{}
	app.traceEmitter = emitter

	app.finishCall(GatewayCallCompletion{
		Kind: CompletionKindRouted,
		Call: CallContext{
			RequestID: "req_external_trace_hook",
			Project:   Project{ID: "prj_external_trace"},
			Key:       APIKey{ID: "key_external_trace"},
			Model:     Model{Name: "gpt-trace"},
		},
		Route: RouteSelection{
			Provider: Provider{ID: "prv_external_trace", Type: ProviderMock, APIKey: "provider-secret"},
			Route:    ModelRoute{ID: "route_external_trace"},
		},
		Usage:           Usage{PromptTokens: 3, CompletionTokens: 4, TotalTokens: 7},
		StatusCode:      http.StatusOK,
		RequestPayload:  map[string]any{"input": "raw prompt sentinel"},
		ResponsePayload: map[string]any{"output": "ok"},
	})

	if completions := emitter.take(); len(completions) != 1 {
		t.Fatalf("trace emitter completions = %d, want 1", len(completions))
	}
	logs := store.ListRequestLogs()
	if len(logs) != 1 || logs[0].RequestID != "req_external_trace_hook" || logs[0].StatusCode != http.StatusOK {
		t.Fatalf("trace hook affected settlement logs: %+v", logs)
	}

	for _, event := range store.ListAuditEvents() {
		if event.Action == "plugin.gateway.trace_export" && event.ResourceID == "req_external_trace_hook" {
			t.Fatalf("external trace hook executed without isolation: %s", event.AfterSnapshot)
		}
	}
}

func TestExternalTraceHookFixtureFailureDoesNotAffectGatewayCompletion(t *testing.T) {
	store := NewMemoryStore()
	app := NewWithConfig(store, Config{
		AdminToken: "external-trace-admin",
		PluginDir:  filepath.Join("..", "plugin", "testdata", "external-trace-hook"),
	})
	emitter := &recordingTraceEmitter{}
	app.traceEmitter = emitter

	app.finishCall(GatewayCallCompletion{
		Kind: CompletionKindRouted,
		Call: CallContext{
			RequestID: "req_external_trace_fail_command",
			Project:   Project{ID: "prj_external_trace"},
			Key:       APIKey{ID: "key_external_trace"},
			Model:     Model{Name: "gpt-trace"},
		},
		Route:      RouteSelection{Provider: Provider{ID: "prv_external_trace", Type: ProviderMock}, Route: ModelRoute{ID: "route_external_trace"}},
		Usage:      Usage{PromptTokens: 3, CompletionTokens: 4, TotalTokens: 7},
		StatusCode: http.StatusOK,
	})

	if completions := emitter.take(); len(completions) != 1 {
		t.Fatalf("trace emitter completions = %d, want 1", len(completions))
	}
	logs := store.ListRequestLogs()
	if len(logs) != 1 || logs[0].RequestID != "req_external_trace_fail_command" || logs[0].StatusCode != http.StatusOK {
		t.Fatalf("external trace hook failure affected settlement logs: %+v", logs)
	}
	for _, event := range store.ListAuditEvents() {
		if event.Action == "plugin.gateway.trace_export" && strings.Contains(event.AfterSnapshot, "provider-secret") {
			t.Fatalf("failed trace hook leaked credentials: %s", event.AfterSnapshot)
		}
	}
}
