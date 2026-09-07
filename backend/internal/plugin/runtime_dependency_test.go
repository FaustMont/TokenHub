package plugin

import (
	"path/filepath"
	"testing"
)

func TestRuntimeLoadsDependenciesIndependentOfDirectoryOrder(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, filepath.Join(root, "a-consumer"), dependencyManifestForTest("tokenhub.consumer", "1.0.0", "tokenhub.core", "^1.0.0"))
	writeManifest(t, filepath.Join(root, "z-core"), dependencyManifestForTest("tokenhub.core", "1.4.0", "", ""))
	registry := NewRegistry()

	packages, err := NewRuntime(root).LoadInto(registry, NewGatewayChainRegistry())
	if err != nil {
		t.Fatalf("load runtime: %v", err)
	}
	if len(packages) != 2 {
		t.Fatalf("packages = %d, want 2", len(packages))
	}
	for _, pluginID := range []string{"tokenhub.consumer", "tokenhub.core"} {
		descriptor, ok := registry.Describe(pluginID)
		if !ok || descriptor.Status != StatusEnabled {
			t.Fatalf("descriptor %s = %+v, %t; want enabled", pluginID, descriptor, ok)
		}
	}
}

func TestRuntimeRejectsUnsatisfiedDependencies(t *testing.T) {
	for _, tc := range []struct {
		name        string
		coreVersion string
	}{
		{name: "missing"},
		{name: "incompatible", coreVersion: "2.0.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeManifest(t, filepath.Join(root, "consumer"), dependencyManifestForTest("tokenhub.consumer", "1.0.0", "tokenhub.core", "^1.0.0"))
			if tc.coreVersion != "" {
				writeManifest(t, filepath.Join(root, "core"), dependencyManifestForTest("tokenhub.core", tc.coreVersion, "", ""))
			}
			registry := NewRegistry()

			packages, err := NewRuntime(root).LoadInto(registry, NewGatewayChainRegistry())
			if err != nil {
				t.Fatalf("load runtime: %v", err)
			}
			if _, ok := registry.Describe("tokenhub.consumer"); ok {
				t.Fatal("consumer with unsatisfied dependency was registered")
			}
			var consumer Package
			for _, pkg := range packages {
				if pkg.Manifest.ID == "tokenhub.consumer" {
					consumer = pkg
					break
				}
			}
			if consumer.State.Status != StatusFailedStartup || consumer.State.LastErrorCode != "plugin_startup_failed" {
				t.Fatalf("consumer state = %+v, want failed startup", consumer.State)
			}
		})
	}
}

func dependencyManifestForTest(id string, version string, dependencyID string, constraint string) string {
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
