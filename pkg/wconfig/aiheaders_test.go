// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package wconfig

import "testing"

func TestValidateAIHeadersAcceptsValid(t *testing.T) {
	headers := map[string]string{
		"x-opencode-session": "wave-local-static-session",
		"X-Custom-Token":     "abc123",
		"x-request-id":       "123e4567-e89b-12d3-a456-426614174000",
	}
	if err := ValidateAIHeaders(headers); err != nil {
		t.Fatalf("expected valid headers to pass, got error: %v", err)
	}
}

func TestValidateAIHeadersAcceptsEmpty(t *testing.T) {
	if err := ValidateAIHeaders(nil); err != nil {
		t.Fatalf("expected nil headers to pass, got error: %v", err)
	}
	if err := ValidateAIHeaders(map[string]string{}); err != nil {
		t.Fatalf("expected empty headers to pass, got error: %v", err)
	}
}

func TestValidateAIHeadersRejectsInvalidName(t *testing.T) {
	cases := []string{
		"x opencode session",
		"x-opencode-session:",
		"x-opencode\nsession",
		"",
	}
	for _, name := range cases {
		if err := ValidateAIHeaders(map[string]string{name: "value"}); err == nil {
			t.Errorf("expected header name %q to be rejected", name)
		}
	}
}

func TestValidateAIHeadersRejectsCRLFValue(t *testing.T) {
	if err := ValidateAIHeaders(map[string]string{"x-custom": "value\r\nx-injected: evil"}); err == nil {
		t.Error("expected CRLF in header value to be rejected")
	}
	if err := ValidateAIHeaders(map[string]string{"x-custom": "value\nx-injected: evil"}); err == nil {
		t.Error("expected LF in header value to be rejected")
	}
}

func TestValidateAIHeadersRejectsEmptyValue(t *testing.T) {
	if err := ValidateAIHeaders(map[string]string{"x-custom": ""}); err == nil {
		t.Error("expected empty header value to be rejected")
	}
}

func TestValidateAIHeadersRejectsProtected(t *testing.T) {
	protected := []string{"Authorization", "authorization", "Cookie", "set-cookie"}
	for _, name := range protected {
		if err := ValidateAIHeaders(map[string]string{name: "value"}); err == nil {
			t.Errorf("expected protected header %q to be rejected", name)
		}
	}
}

func TestValidateAIHeadersRejectsTooMany(t *testing.T) {
	headers := make(map[string]string)
	for i := 0; i < 17; i++ {
		headers[string(rune('a'+i))+"-header"] = "value"
	}
	if err := ValidateAIHeaders(headers); err == nil {
		t.Error("expected more than 16 headers to be rejected")
	}
}
