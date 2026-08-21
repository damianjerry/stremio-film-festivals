package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestVersionMatchesReleaseManifest is the guard that keeps the const in
// config.go honest. release-please rewrites it via the "extra-files" entry in
// release-please-config.json, keyed off the `x-release-please-version`
// annotation comment -- drop the annotation, drop the extra-files entry, or
// reformat the line so the updater no longer matches it, and the const silently
// stops tracking releases. That failure is invisible in review and only shows
// up as a stale `version` in /manifest.json, which is what Stremio clients read
// to decide whether an installed addon has an update.
func TestVersionMatchesReleaseManifest(t *testing.T) {
	raw, err := os.ReadFile(".release-please-manifest.json")
	if err != nil {
		t.Fatalf("Couldn't read .release-please-manifest.json: %v", err)
	}

	var manifest map[string]string
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("Couldn't parse .release-please-manifest.json: %v", err)
	}

	released, ok := manifest["."]
	if !ok {
		t.Fatal(`.release-please-manifest.json has no "." entry`)
	}

	if version != released {
		t.Fatalf("version const is %q but the last release was %q -- release-please is no longer rewriting config.go; check the x-release-please-version annotation and the extra-files entry in release-please-config.json", version, released)
	}
}

// TestDefaultRedirectURLIsRelative pins the root redirect to a same-host target.
// It used to point at an unrelated external site inherited from the project this
// was forked from, which sent anyone typing the bare hostname off-site instead of
// to the addon.
func TestDefaultRedirectURLIsRelative(t *testing.T) {
	if !strings.HasPrefix(defaultRedirectURL, "/") {
		t.Fatalf("defaultRedirectURL should be a relative, same-host path so it works for every deployment, got %q", defaultRedirectURL)
	}
	if defaultRedirectURL != "/configure" {
		t.Fatalf("expected the root to land on the configuration UI, got %q", defaultRedirectURL)
	}
}

// TestGetEnvAllowEmptyKeepsEmptyValues covers the reason REDIRECT_URL does not
// read through getEnv: an explicitly empty value is how the root redirect is
// turned off, and getEnv would collapse it back into the default, so the option
// would work with the flag and silently do nothing via the environment.
func TestGetEnvAllowEmptyKeepsEmptyValues(t *testing.T) {
	const key = "REDIRECT_URL_TEST_ONLY"

	if got := getEnvAllowEmpty(key, "/configure"); got != "/configure" {
		t.Fatalf("unset should fall back to the default, got %q", got)
	}

	t.Setenv(key, "")
	if got := getEnvAllowEmpty(key, "/configure"); got != "" {
		t.Fatalf("an explicitly empty value should survive, got %q", got)
	}
	if got := getEnv(key, "/configure"); got != "/configure" {
		t.Fatalf("getEnv is expected to collapse empty into the default -- if this changed, getEnvAllowEmpty may no longer be needed, got %q", got)
	}

	t.Setenv(key, "https://example.org")
	if got := getEnvAllowEmpty(key, "/configure"); got != "https://example.org" {
		t.Fatalf("a set value should win over the default, got %q", got)
	}
}
