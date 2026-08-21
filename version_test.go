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

// TestConfigurePageAvoidsBlockedClassNames keeps the hosting-attribution card
// out of the way of content blockers. EasyList -- the default list in uBlock
// Origin, AdGuard, Brave and others -- ships a generic cosmetic rule for the
// class "sponsor-text" that hides any element carrying it, on every site. The
// card used exactly that class, so the attribution paragraph was invisible to
// every visitor running a blocker, while the logo beside it rendered normally.
//
// This reads the source rather than a rendered page because the markup is built
// inline in HandleConfigureEndpoint with no separately callable renderer. That
// is a little unusual for a test, but it is the honest way to lint a string
// literal, and it covers the CSS block as well as the markup.
func TestConfigurePageAvoidsBlockedClassNames(t *testing.T) {
	source, err := os.ReadFile("configure.go")
	if err != nil {
		t.Fatalf("Couldn't read configure.go: %v", err)
	}
	page := string(source)

	// "sponsor-text" is the class EasyList actually blocks today. The others are
	// close enough neighbours to be worth staying clear of, since generic rules
	// accrete over time and this failure mode is silent.
	for _, blocked := range []string{"sponsor-text", "sponsor-content", "sponsor-logo", "sponsored"} {
		// Skip this test file's own mention of the names in the loop above.
		if strings.Contains(page, `class="`+blocked) || strings.Contains(page, "."+blocked+" {") {
			t.Fatalf("the configure page uses the class %q, which content blockers hide with a generic cosmetic rule -- the attribution silently disappears for anyone running uBlock Origin, AdGuard or Brave", blocked)
		}
	}

	// ...and the attribution itself has to survive the renaming.
	if !strings.Contains(page, "instance-note-text") {
		t.Fatal("the hosting attribution card is missing from the configure page")
	}
}
