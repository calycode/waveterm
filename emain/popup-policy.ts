// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

export type PopupDecision = "allow" | "prompt" | "deny";

export type PopupPolicy = {
    enabled: boolean;
    allowedOrigins: string[];
    askForUnknown: boolean;
    allowLocalhostHttp?: boolean;
};

function isLoopbackHost(hostname: string): boolean {
    return hostname === "localhost" || hostname === "127.0.0.1" || hostname === "[::1]";
}

// normalizeOrigin canonicalizes a URL down to scheme://host[:port]. Only https
// is accepted; http is allowed only for loopback when explicitly enabled. The
// WHATWG URL parser already lowercases the host and drops default ports.
export function normalizeOrigin(url: string, allowLocalhostHttp = false): string | null {
    let parsed: URL;
    try {
        parsed = new URL(url);
    } catch {
        return null;
    }
    const scheme = parsed.protocol.replace(/:$/, "").toLowerCase();
    const host = parsed.hostname.toLowerCase();
    const port = parsed.port;
    const withPort = port ? `${scheme}://${host}:${port}` : `${scheme}://${host}`;
    if (scheme === "https") {
        return withPort;
    }
    if (scheme === "http" && allowLocalhostHttp && isLoopbackHost(host)) {
        return withPort;
    }
    return null;
}

// isOriginAllowed compares the normalized origin of targetUrl against each
// allowlist entry using strict equality (no wildcards, no subdomain matching).
export function isOriginAllowed(targetUrl: string, allowedOrigins: string[], allowLocalhostHttp = false): boolean {
    const target = normalizeOrigin(targetUrl, allowLocalhostHttp);
    if (target == null) {
        return false;
    }
    for (const entry of allowedOrigins) {
        if (normalizeOrigin(entry, allowLocalhostHttp) === target) {
            return true;
        }
    }
    return false;
}

export function evaluatePopupPolicy(targetUrl: string, policy: PopupPolicy): PopupDecision {
    if (!policy.enabled) {
        return "deny";
    }
    const allowLocalhostHttp = policy.allowLocalhostHttp ?? false;
    if (normalizeOrigin(targetUrl, allowLocalhostHttp) == null) {
        return "deny";
    }
    if (isOriginAllowed(targetUrl, policy.allowedOrigins ?? [], allowLocalhostHttp)) {
        return "allow";
    }
    if (policy.askForUnknown) {
        return "prompt";
    }
    return "deny";
}
