// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package aiusechat

import (
	"testing"

	"github.com/wavetermdev/waveterm/pkg/aiusechat/uctypes"
	"github.com/wavetermdev/waveterm/pkg/wconfig"
)

func TestApplyProviderDefaultsGroq(t *testing.T) {
	config := wconfig.AIModeConfigType{
		Provider: uctypes.AIProvider_Groq,
	}
	applyProviderDefaults(&config)
	if config.APIType != uctypes.APIType_OpenAIChat {
		t.Fatalf("expected API type %q, got %q", uctypes.APIType_OpenAIChat, config.APIType)
	}
	if config.Endpoint != GroqChatEndpoint {
		t.Fatalf("expected endpoint %q, got %q", GroqChatEndpoint, config.Endpoint)
	}
	if config.APITokenSecretName != GroqAPITokenSecretName {
		t.Fatalf("expected API token secret name %q, got %q", GroqAPITokenSecretName, config.APITokenSecretName)
	}
}

func TestApplyProviderDefaultsKeepsProxyURL(t *testing.T) {
	config := wconfig.AIModeConfigType{
		Provider: uctypes.AIProvider_OpenAI,
		Model:    "gpt-5-mini",
		ProxyURL: "http://localhost:8080",
	}
	applyProviderDefaults(&config)
	if config.ProxyURL != "http://localhost:8080" {
		t.Fatalf("expected proxy URL to be preserved, got %q", config.ProxyURL)
	}
}

func TestComputeResolvedAIModeConfigsStripsHeaders(t *testing.T) {
	fullConfig := wconfig.FullConfigType{
		WaveAIModes: map[string]wconfig.AIModeConfigType{
			"opencode-deepseek-flash": {
				DisplayName: "DeepSeek V4.1 Flash",
				Provider:    uctypes.AIProvider_Custom,
				APIType:     uctypes.APIType_OpenAIChat,
				Model:       "deepseek-v4.1-flash",
				Headers:     map[string]string{"x-opencode-session": "secret-value"},
			},
		},
	}

	resolved := ComputeResolvedAIModeConfigs(fullConfig)

	cfg, ok := resolved["opencode-deepseek-flash"]
	if !ok {
		t.Fatal("expected mode to be resolved")
	}
	if len(cfg.Headers) != 0 {
		t.Errorf("expected headers to be stripped from broadcast config, got %v", cfg.Headers)
	}
	if cfg.Model != "deepseek-v4.1-flash" {
		t.Errorf("expected other fields preserved, got model %q", cfg.Model)
	}
}
