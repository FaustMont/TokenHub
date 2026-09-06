package main

import (
	"fmt"
	"strings"
)

type manifest struct {
	SchemaVersion int    `yaml:"schema_version"`
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	Summary       string `yaml:"summary"`
	Description   string `yaml:"description"`
	Category      string `yaml:"category"`
	HostAdapter   string `yaml:"host_adapter"`
	Dependencies  []struct {
		ID      string `yaml:"id"`
		Version string `yaml:"version"`
	} `yaml:"dependencies"`
	Settings struct {
		Scopes []string `yaml:"scopes"`
	} `yaml:"settings"`
	TokenHub struct {
		PluginAPI string `yaml:"plugin_api"`
		MinCore   string `yaml:"min_core"`
		MaxCore   string `yaml:"max_core"`
	} `yaml:"tokenhub"`
	Kinds     []string `yaml:"kinds"`
	Placement []string `yaml:"placement"`
	Entry     struct {
		Backend *struct {
			Protocol string `yaml:"protocol"`
			Command  string `yaml:"command"`
		} `yaml:"backend"`
	} `yaml:"entry"`
	Capabilities struct {
		ProviderTypes         []string                `yaml:"provider_types"`
		ProviderResourceTypes []string                `yaml:"provider_resource_types"`
		Provider              map[string]any          `yaml:"provider"`
		Actions               []manifestAction        `yaml:"actions"`
		Hooks                 []manifestHook          `yaml:"hooks"`
		Background            []manifestBackgroundJob `yaml:"background_jobs"`
		Gateway               []string                `yaml:"gateway"`
	} `yaml:"capabilities"`
	Permissions struct {
		Data struct {
			Read []string `yaml:"read"`
		} `yaml:"data"`
	} `yaml:"permissions"`
	Distribution map[string]any `yaml:"distribution"`
}

type manifestAction struct {
	ID           string            `yaml:"id"`
	Kind         string            `yaml:"kind"`
	Title        string            `yaml:"title"`
	Capability   string            `yaml:"capability"`
	Subject      string            `yaml:"subject"`
	Metadata     map[string]string `yaml:"metadata"`
	InputSchema  map[string]any    `yaml:"input_schema"`
	OutputSchema map[string]any    `yaml:"output_schema"`
}

type manifestHook struct {
	ID            string   `yaml:"id"`
	Stage         string   `yaml:"stage"`
	Priority      int      `yaml:"priority"`
	Before        []string `yaml:"before"`
	After         []string `yaml:"after"`
	FailurePolicy string   `yaml:"failure_policy"`
	Reads         []string `yaml:"reads"`
	Writes        []string `yaml:"writes"`
}

func validateManifestAPI(value manifest) error {
	switch {
	case value.SchemaVersion == 1 && value.TokenHub.PluginAPI == "v1":
		return nil
	case value.SchemaVersion == 2 && value.TokenHub.PluginAPI == "v2":
		if strings.TrimSpace(value.Summary) == "" {
			return fmt.Errorf("plugin API v2 manifest must declare summary")
		}
		switch value.Category {
		case "provider_integration", "request_pipeline", "ui_template", "automation":
			return nil
		default:
			return fmt.Errorf("plugin API v2 manifest category %q is unsupported", value.Category)
		}
	default:
		return fmt.Errorf("unsupported manifest schema/plugin API pair %d/%q", value.SchemaVersion, value.TokenHub.PluginAPI)
	}
}

type manifestBackgroundJob struct {
	ID             string                  `yaml:"id"`
	Title          string                  `yaml:"title"`
	Capability     string                  `yaml:"capability"`
	Subject        string                  `yaml:"subject"`
	Schedule       string                  `yaml:"schedule"`
	TimeoutMillis  int                     `yaml:"timeout_millis"`
	MaxConcurrency int                     `yaml:"max_concurrency"`
	Retry          manifestBackgroundRetry `yaml:"retry"`
	InputSchema    map[string]any          `yaml:"input_schema"`
	OutputSchema   map[string]any          `yaml:"output_schema"`
}

type manifestBackgroundRetry struct {
	MaxAttempts   int `yaml:"max_attempts"`
	BackoffMillis int `yaml:"backoff_millis"`
}
