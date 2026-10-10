// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package wconfig

import (
	"fmt"
	"regexp"
	"strings"
)

const MaxAIHeaders = 16

// HTTP header field names are RFC 7230 tokens. Values are validated separately
// because they must not contain CR/LF (header injection).
var httpTokenRe = regexp.MustCompile("^[!#$%&'*+\\-.^_\\x60|~0-9A-Za-z]+$")

// Credential-bearing headers are reserved for Wave's own token mechanism so a
// user-configured ai:header can never hijack authentication.
var protectedAIHeaders = map[string]bool{
	"authorization": true,
	"cookie":        true,
	"set-cookie":    true,
}

func IsProtectedAIHeader(name string) bool {
	return protectedAIHeaders[strings.ToLower(name)]
}

func ValidateAIHeaders(headers map[string]string) error {
	if len(headers) > MaxAIHeaders {
		return fmt.Errorf("too many ai:headers (%d), maximum is %d", len(headers), MaxAIHeaders)
	}
	for name, value := range headers {
		if !httpTokenRe.MatchString(name) {
			return fmt.Errorf("invalid ai:header name %q", name)
		}
		if IsProtectedAIHeader(name) {
			return fmt.Errorf("ai:header %q is reserved and cannot be set", name)
		}
		if value == "" {
			return fmt.Errorf("ai:header %q has an empty value", name)
		}
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("ai:header %q value must not contain CR or LF", name)
		}
	}
	return nil
}
