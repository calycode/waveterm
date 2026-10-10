// Copyright 2026, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

import { describe, expect, it } from "vitest";
import { evaluatePopupPolicy, isOriginAllowed, normalizeOrigin } from "./popup-policy";

describe("normalizeOrigin", () => {
    it("lowercases the host and strips the default https port", () => {
        expect(normalizeOrigin("https://Example.COM:443/path?q=1")).toBe("https://example.com");
    });

    it("keeps a non-default port", () => {
        expect(normalizeOrigin("https://example.com:8443/app")).toBe("https://example.com:8443");
    });

    it("accepts a bare origin", () => {
        expect(normalizeOrigin("https://docs.calycode.dev")).toBe("https://docs.calycode.dev");
    });

    it("rejects non-https schemes", () => {
        expect(normalizeOrigin("http://example.com")).toBeNull();
        expect(normalizeOrigin("file:///etc/passwd")).toBeNull();
        expect(normalizeOrigin("javascript:alert(1)")).toBeNull();
    });

    it("allows http for loopback only when explicitly enabled", () => {
        expect(normalizeOrigin("http://localhost:3000", true)).toBe("http://localhost:3000");
        expect(normalizeOrigin("http://127.0.0.1:3000", true)).toBe("http://127.0.0.1:3000");
        expect(normalizeOrigin("http://localhost:3000", false)).toBeNull();
        expect(normalizeOrigin("http://example.com", true)).toBeNull();
    });

    it("returns null for garbage input", () => {
        expect(normalizeOrigin("not a url")).toBeNull();
        expect(normalizeOrigin("")).toBeNull();
    });
});

describe("isOriginAllowed", () => {
    const allowed = ["https://app.example.com", "https://docs.calycode.dev:8443"];

    it("matches an exact origin regardless of path", () => {
        expect(isOriginAllowed("https://app.example.com/deep/link", allowed)).toBe(true);
    });

    it("does not match a subdomain", () => {
        expect(isOriginAllowed("https://evil.app.example.com", allowed)).toBe(false);
    });

    it("does not match a parent domain", () => {
        expect(isOriginAllowed("https://example.com", allowed)).toBe(false);
    });

    it("respects ports", () => {
        expect(isOriginAllowed("https://docs.calycode.dev:8443/x", allowed)).toBe(true);
        expect(isOriginAllowed("https://docs.calycode.dev/x", allowed)).toBe(false);
    });

    it("ignores malformed allowlist entries", () => {
        expect(isOriginAllowed("https://app.example.com", ["not a url", "https://app.example.com"])).toBe(true);
    });
});

describe("evaluatePopupPolicy", () => {
    const base = {
        enabled: true,
        allowedOrigins: ["https://app.example.com"],
        askForUnknown: false,
    };

    it("allows an allowlisted origin", () => {
        expect(evaluatePopupPolicy("https://app.example.com/login", base)).toBe("allow");
    });

    it("denies an unknown origin when askForUnknown is false", () => {
        expect(evaluatePopupPolicy("https://evil.example.com", base)).toBe("deny");
    });

    it("prompts for an unknown origin when askForUnknown is true", () => {
        expect(evaluatePopupPolicy("https://evil.example.com", { ...base, askForUnknown: true })).toBe("prompt");
    });

    it("denies everything when disabled", () => {
        expect(evaluatePopupPolicy("https://app.example.com", { ...base, enabled: false })).toBe("deny");
    });

    it("denies non-https targets", () => {
        expect(evaluatePopupPolicy("http://app.example.com", base)).toBe("deny");
    });
});
