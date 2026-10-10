// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package uctypes

import (
	"net/http"
	"testing"
)

func TestApplyHeadersSetsCustomHeader(t *testing.T) {
	opts := AIOptsType{
		Headers: map[string]string{"x-opencode-session": "wave-local-static-session"},
	}
	req, err := http.NewRequest(http.MethodPost, "https://opencode.ai/zen/go/v1/chat/completions", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	opts.ApplyHeaders(req)

	if got := req.Header.Get("x-opencode-session"); got != "wave-local-static-session" {
		t.Errorf("x-opencode-session: got %q, want %q", got, "wave-local-static-session")
	}
}

func TestApplyHeadersDoesNotOverrideAuthorization(t *testing.T) {
	opts := AIOptsType{
		Headers: map[string]string{"Authorization": "Bearer evil", "Cookie": "evil=1"},
	}
	req, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer real-token")

	opts.ApplyHeaders(req)

	if got := req.Header.Get("Authorization"); got != "Bearer real-token" {
		t.Errorf("Authorization was overridden: got %q", got)
	}
	if got := req.Header.Get("Cookie"); got != "" {
		t.Errorf("Cookie should not be set, got %q", got)
	}
}

func TestApplyHeadersNilSafe(t *testing.T) {
	opts := AIOptsType{}
	req, _ := http.NewRequest(http.MethodPost, "https://example.com", nil)
	opts.ApplyHeaders(req)
}
