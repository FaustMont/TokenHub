package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pluginmeta "tokenhub/backend/internal/plugin"
)

func TestAdminPluginLifecycleRejectsUnsatisfiedDependencies(t *testing.T) {
	t.Run("enable", func(t *testing.T) {
		pluginDir := t.TempDir()
		consumerDir := filepath.Join(pluginDir, "consumer")
		writeServerPluginManifest(t, consumerDir, serverDependencyManifest("tokenhub.consumer", "1.0.0", "tokenhub.missing", "^1.0.0"))
		if err := os.WriteFile(filepath.Join(consumerDir, "plugin.state.json"), []byte(`{"status":"disabled"}`), 0o644); err != nil {
			t.Fatal(err)
		}
		server := NewWithConfig(NewMemoryStore(), Config{AdminToken: "dev_admin_token", PluginDir: pluginDir})

		response := doJSON(t, server.Handler(), http.MethodPatch, "/api/admin/plugins/tokenhub.consumer/state", map[string]any{"status": "enabled"}, "dev_admin_token")
		assertResponseBodyJSONError(t, response, http.StatusConflict, "plugin_dependency_unsatisfied")
		consumer, ok := server.pluginRegistry.Describe("tokenhub.consumer")
		if !ok || consumer.Status != pluginmeta.StatusDisabled {
			t.Fatalf("consumer descriptor = %+v, %t; want disabled", consumer, ok)
		}
	})

	t.Run("disable and uninstall", func(t *testing.T) {
		pluginDir := t.TempDir()
		writeServerPluginManifest(t, filepath.Join(pluginDir, "core"), serverDependencyManifest("tokenhub.core", "1.4.0", "", ""))
		writeServerPluginManifest(t, filepath.Join(pluginDir, "consumer"), serverDependencyManifest("tokenhub.consumer", "1.0.0", "tokenhub.core", "^1.0.0"))
		server := NewWithConfig(NewMemoryStore(), Config{AdminToken: "dev_admin_token", PluginDir: pluginDir})

		disable := doJSON(t, server.Handler(), http.MethodPatch, "/api/admin/plugins/tokenhub.core/state", map[string]any{"status": "disabled"}, "dev_admin_token")
		assertResponseBodyJSONError(t, disable, http.StatusConflict, "plugin_dependency_in_use")
		uninstall := doJSON(t, server.Handler(), http.MethodDelete, "/api/admin/plugin-packages/tokenhub.core", nil, "dev_admin_token")
		assertResponseBodyJSONError(t, uninstall, http.StatusConflict, "plugin_dependency_in_use")
		if _, err := os.Stat(filepath.Join(pluginDir, "core", "plugin.yaml")); err != nil {
			t.Fatalf("required dependency was removed: %v", err)
		}
	})
}

func TestAdminPluginUpdateRejectsMismatchedPluginID(t *testing.T) {
	pluginDir := t.TempDir()
	currentDir := filepath.Join(pluginDir, "current")
	archive := adminPluginZip(t, map[string]string{
		"plugin.yaml": adminPluginManifest("tokenhub.other", "Other Plugin", "2.0.0"),
	})
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive)
	}))
	defer upstream.Close()
	writeServerPluginManifest(t, currentDir, adminPluginManifestWithTestDistribution("tokenhub.current", "Current Plugin", "1.0.0", upstream.URL+"/other.zip", adminSHA256Hex(archive), ""))
	server := NewWithConfig(NewMemoryStore(), Config{AdminToken: "dev_admin_token", PluginDir: pluginDir})
	server.pluginInstallClient = upstream.Client()

	response := doJSON(t, server.Handler(), http.MethodPost, "/api/admin/plugins/tokenhub.current/update", map[string]any{}, "dev_admin_token")
	assertResponseBodyJSONError(t, response, http.StatusBadRequest, "plugin_id_mismatch")
	if _, err := os.Stat(filepath.Join(currentDir, "plugin.yaml")); err != nil {
		t.Fatalf("target package changed after ID mismatch: %v", err)
	}
	if _, err := os.Stat(filepath.Join(pluginDir, "tokenhub.other")); !os.IsNotExist(err) {
		t.Fatalf("mismatched package was installed: %v", err)
	}
}

func TestAdminPluginUpdatePublishesNewActions(t *testing.T) {
	pluginDir := t.TempDir()
	archive := adminPluginZip(t, map[string]string{
		"plugin.yaml": adminPluginActionManifest("tokenhub.hot-reload", "1.1.0", "sync.new", "", ""),
	})
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive)
	}))
	defer upstream.Close()
	writeServerPluginManifest(t, filepath.Join(pluginDir, "hot-reload"), adminPluginActionManifest("tokenhub.hot-reload", "1.0.0", "sync.old", upstream.URL+"/hot-reload.zip", adminSHA256Hex(archive)))
	server := NewWithConfig(NewMemoryStore(), Config{AdminToken: "dev_admin_token", PluginDir: pluginDir})
	server.pluginInstallClient = upstream.Client()

	before := doJSON(t, server.Handler(), http.MethodGet, "/api/admin/plugin-actions", nil, "dev_admin_token")
	if !strings.Contains(before.Body, `"action_id":"sync.old"`) {
		t.Fatalf("old action missing before update: %s", before.Body)
	}
	update := doJSON(t, server.Handler(), http.MethodPost, "/api/admin/plugins/tokenhub.hot-reload/update", map[string]any{}, "dev_admin_token")
	if update.Code != http.StatusOK {
		t.Fatalf("update plugin: expected 200, got %d: %s", update.Code, update.Body)
	}
	after := doJSON(t, server.Handler(), http.MethodGet, "/api/admin/plugin-actions", nil, "dev_admin_token")
	if !strings.Contains(after.Body, `"action_id":"sync.new"`) || strings.Contains(after.Body, `"action_id":"sync.old"`) {
		t.Fatalf("runtime actions were not replaced after update: %s", after.Body)
	}
}

func TestAdminPluginUpdateRejectsDependencyBreakingVersion(t *testing.T) {
	pluginDir := t.TempDir()
	archive := adminPluginZip(t, map[string]string{
		"plugin.yaml": serverDependencyManifest("tokenhub.core", "2.0.0", "", ""),
	})
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive)
	}))
	defer upstream.Close()
	writeServerPluginManifest(t, filepath.Join(pluginDir, "core"), adminPluginManifestWithTestDistribution("tokenhub.core", "Core Plugin", "1.4.0", upstream.URL+"/core.zip", adminSHA256Hex(archive), "automation"))
	writeServerPluginManifest(t, filepath.Join(pluginDir, "consumer"), serverDependencyManifest("tokenhub.consumer", "1.0.0", "tokenhub.core", "^1.0.0"))
	server := NewWithConfig(NewMemoryStore(), Config{AdminToken: "dev_admin_token", PluginDir: pluginDir})
	server.pluginInstallClient = upstream.Client()

	response := doJSON(t, server.Handler(), http.MethodPost, "/api/admin/plugins/tokenhub.core/update", map[string]any{}, "dev_admin_token")
	assertResponseBodyJSONError(t, response, http.StatusConflict, "plugin_dependency_unsatisfied")
	core, ok := server.pluginRegistry.Describe("tokenhub.core")
	if !ok || core.Version != "1.4.0" {
		t.Fatalf("core descriptor = %+v, %t; want version 1.4.0", core, ok)
	}
}

func TestReloadPluginRuntimeWaitsForRequestSnapshot(t *testing.T) {
	server := NewWithConfig(NewMemoryStore(), Config{PluginDir: t.TempDir()})
	requestEntered := make(chan struct{})
	releaseRequest := make(chan struct{})
	requestDone := make(chan struct{})
	originalRegistry := server.pluginRegistry
	server.mux.HandleFunc("GET /test/plugin-runtime-snapshot", func(w http.ResponseWriter, _ *http.Request) {
		close(requestEntered)
		<-releaseRequest
		if server.pluginRegistry != originalRegistry {
			t.Error("request observed a different plugin registry within one runtime snapshot")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	go func() {
		defer close(requestDone)
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test/plugin-runtime-snapshot", nil))
		if response.Code != http.StatusNoContent {
			t.Errorf("snapshot request status = %d", response.Code)
		}
	}()
	<-requestEntered

	reloadDone := make(chan error, 1)
	go func() { reloadDone <- server.reloadPluginRuntime(context.Background()) }()
	select {
	case err := <-reloadDone:
		close(releaseRequest)
		<-requestDone
		t.Fatalf("reload completed while request held a runtime snapshot: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseRequest)
	<-requestDone
	select {
	case err := <-reloadDone:
		if err != nil {
			t.Fatalf("reload plugin runtime: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("reload did not complete after request released its snapshot")
	}
	if server.pluginRegistry == originalRegistry {
		t.Fatal("reload did not publish a new plugin registry")
	}
}

func TestReloadPluginRuntimeWaitsForScheduledJobSnapshot(t *testing.T) {
	server := NewWithConfig(NewMemoryStore(), Config{PluginDir: t.TempDir()})
	jobEntered := make(chan struct{})
	releaseJob := make(chan struct{})
	jobDone := make(chan struct{})
	originalRegistry := server.pluginRegistry
	broker := pluginmeta.NewBackgroundJobBroker()
	if err := broker.Register(pluginmeta.BackgroundJobDescriptor{
		PluginID:       "tokenhub.snapshot",
		JobID:          "snapshot.check",
		Schedule:       "1m",
		MaxConcurrency: 1,
	}, pluginmeta.BackgroundJobHandlerFunc(func(context.Context, pluginmeta.BackgroundJobInvocation) (pluginmeta.BackgroundJobResult, error) {
		close(jobEntered)
		<-releaseJob
		if server.pluginRegistry != originalRegistry {
			t.Error("scheduled job observed a different plugin registry within one runtime snapshot")
		}
		return pluginmeta.BackgroundJobResult{}, nil
	})); err != nil {
		t.Fatal(err)
	}
	server.pluginBackgroundRunner.SetBroker(broker)
	go func() {
		defer close(jobDone)
		server.pluginBackgroundRunner.RunDue(context.Background(), time.Now().UTC(), "schedule")
	}()
	<-jobEntered

	reloadDone := make(chan error, 1)
	go func() { reloadDone <- server.reloadPluginRuntime(context.Background()) }()
	select {
	case err := <-reloadDone:
		close(releaseJob)
		<-jobDone
		t.Fatalf("reload completed while scheduled job held a runtime snapshot: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseJob)
	<-jobDone
	select {
	case err := <-reloadDone:
		if err != nil {
			t.Fatalf("reload plugin runtime: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("reload did not complete after scheduled job released its snapshot")
	}
}

func TestAdminPluginActionFailurePersistsSafeMessage(t *testing.T) {
	store := NewMemoryStore()
	server := NewWithConfig(store, Config{AdminToken: "dev_admin_token"})
	if err := server.pluginRegistry.Register(pluginmeta.Descriptor{ID: "tokenhub.failure", Name: "Failure", Version: "1.0.0"}); err != nil {
		t.Fatal(err)
	}
	if err := server.pluginActions.Register(pluginmeta.ActionDescriptor{PluginID: "tokenhub.failure", ActionID: "fail", Kind: pluginmeta.ActionKindRead}, pluginmeta.ActionHandlerFunc(func(context.Context, pluginmeta.ActionInvocation) (pluginmeta.ActionResult, error) {
		return pluginmeta.ActionResult{}, errors.New("plugin stderr leaked-secret")
	})); err != nil {
		t.Fatal(err)
	}

	response := doJSON(t, server.Handler(), http.MethodPost, "/api/admin/plugins/tokenhub.failure/actions/fail", map[string]any{"access_token": "request-secret"}, "dev_admin_token")
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("plugin action failure status = %d: %s", response.Code, response.Body)
	}
	for _, secret := range []string{"leaked-secret", "request-secret"} {
		if strings.Contains(response.Body, secret) {
			t.Fatalf("plugin action response leaked %q: %s", secret, response.Body)
		}
	}
	events := store.ListAuditEvents()
	if len(events) == 0 || events[0].Message != "Plugin action failed" {
		t.Fatalf("plugin action audit event = %+v, want safe failure message", events)
	}
	if strings.Contains(events[0].Message, "leaked-secret") || strings.Contains(events[0].Message, "request-secret") {
		t.Fatalf("plugin action audit message leaked secrets: %+v", events[0])
	}
}

func TestPluginDescriptorHasSettingsRequiresRenderedCapability(t *testing.T) {
	descriptor := pluginmeta.Descriptor{Settings: pluginmeta.ManifestSettings{Scopes: []string{"administrator", "project"}}}
	if pluginDescriptorHasSettings(descriptor) {
		t.Fatal("settings scopes without a rendered capability exposed an empty Settings page")
	}
	descriptor.Capabilities = []pluginmeta.CapabilityDescriptor{{Kind: pluginmeta.CapabilityKindSIM, Name: pluginmeta.SIMCapabilityThemeTokens}}
	if !pluginDescriptorHasSettings(descriptor) {
		t.Fatal("SIM theme settings capability did not expose Settings")
	}
}

func serverDependencyManifest(id string, version string, dependencyID string, constraint string) string {
	dependencies := ""
	if dependencyID != "" {
		dependencies = "dependencies:\n  - id: " + dependencyID + "\n    version: '" + constraint + "'\n"
	}
	return `
schema_version: 2
id: ` + id + `
name: Dependency Test Plugin
version: ` + version + `
summary: Exercises plugin dependencies.
category: automation
tokenhub:
  plugin_api: v2
kinds: [extension]
placement: []
` + dependencies + `permissions:
  data:
    read: []
    write: []
`
}

func adminPluginManifestWithTestDistribution(id string, name string, version string, downloadURL string, checksum string, category string) string {
	schemaVersion := "1"
	pluginAPI := "v1"
	summary := ""
	categoryLine := ""
	if category != "" {
		schemaVersion = "2"
		pluginAPI = "v2"
		summary = "summary: Exercises plugin updates.\n"
		categoryLine = "category: " + category + "\n"
	}
	return "schema_version: " + schemaVersion + "\n" +
		"id: " + id + "\n" +
		"name: " + name + "\n" +
		"version: " + version + "\n" +
		summary + categoryLine +
		"distribution:\n  download_url: " + downloadURL + "\n  checksum_sha256: " + checksum + "\n" +
		"tokenhub:\n  plugin_api: " + pluginAPI + "\n" +
		"kinds: [extension]\nplacement: []\npermissions:\n  data:\n    read: []\n    write: []\n"
}

func adminPluginActionManifest(id string, version string, actionID string, downloadURL string, checksum string) string {
	distribution := ""
	if downloadURL != "" {
		distribution = "distribution:\n  download_url: " + downloadURL + "\n  checksum_sha256: " + checksum + "\n"
	}
	return `
schema_version: 1
id: ` + id + `
name: Hot Reload Plugin
version: ` + version + `
` + distribution + `tokenhub:
  plugin_api: v1
kinds: [extension]
placement: [management_action]
capabilities:
  actions:
    - id: ` + actionID + `
      kind: read
      title: Synchronize
`
}
