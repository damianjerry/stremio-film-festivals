package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestLogoEmbedded(t *testing.T) {
	if len(logoBytes) == 0 {
		t.Fatal("Embedded logoBytes is empty")
	}

	// Verify PNG magic header: 0x89 0x50 0x4E 0x47 0x0D 0x0A 0x1A 0x0A
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.HasPrefix(logoBytes, pngHeader) {
		t.Fatal("Embedded logoBytes does not have valid PNG header")
	}

	if len(logoBytes) < 30000 {
		t.Fatalf("Embedded logoBytes unexpectedly small: %d bytes", len(logoBytes))
	}
}

func TestBuildCustomManifestAll(t *testing.T) {
	manifest := BuildCustomManifest(nil, "http://localhost:8080")

	if manifest.ID != "tv.deflix.stremio-film-festivals" {
		t.Fatalf("Unexpected manifest ID: %s", manifest.ID)
	}

	if manifest.Name != "Film Festivals 2 | ElfHosted" {
		t.Fatalf("Unexpected manifest Name: %s", manifest.Name)
	}

	if manifest.Logo != "http://localhost:8080/logo.png" {
		t.Fatalf("Unexpected manifest Logo: %s", manifest.Logo)
	}

	if len(manifest.Catalogs) != 38 {
		t.Fatalf("Expected 38 catalogs in unfiltered manifest, got %d", len(manifest.Catalogs))
	}

	if !manifest.BehaviorHints.Configurable {
		t.Fatal("Expected manifest behaviorHints.configurable to be true")
	}

	if manifest.BehaviorHints.ConfigurationRequired {
		t.Fatal("Expected manifest behaviorHints.configurationRequired to be false")
	}

	if manifest.StremioAddonsConfig == nil || manifest.StremioAddonsConfig.Issuer != "https://stremio-addons.net" {
		t.Fatalf("Expected valid stremioAddonsConfig, got %+v", manifest.StremioAddonsConfig)
	}

	rawJSON, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("Failed to marshal manifest to JSON: %v", err)
	}
	if strings.Contains(string(rawJSON), `"types":null`) || strings.Contains(string(rawJSON), `null`) {
		t.Fatalf("Manifest JSON contains illegal null values: %s", string(rawJSON))
	}
	if !strings.Contains(string(rawJSON), `"resources":["catalog"]`) {
		t.Fatalf("Manifest JSON does not contain schema-compliant \"resources\":[\"catalog\"]: %s", string(rawJSON))
	}
	if !strings.Contains(string(rawJSON), `"stremioAddonsConfig"`) || !strings.Contains(string(rawJSON), `"https://stremio-addons.net"`) {
		t.Fatalf("Manifest JSON does not contain stremioAddonsConfig: %s", string(rawJSON))
	}
}

func TestBuildCustomManifestFiltered(t *testing.T) {
	selected := []string{"cannes-palme-dor", "venice-golden-lion", "berlin-golden-bear"}
	manifest := BuildCustomManifest(selected, "https://stremio-film-festivals.deflix.tv")

	if len(manifest.Catalogs) != 3 {
		t.Fatalf("Expected 3 catalogs in filtered manifest, got %d", len(manifest.Catalogs))
	}

	expectedIDs := map[string]bool{
		"cannes-palme-dor":   true,
		"venice-golden-lion": true,
		"berlin-golden-bear": true,
	}

	for _, cat := range manifest.Catalogs {
		if !expectedIDs[cat.ID] {
			t.Fatalf("Unexpected catalog in filtered manifest: %s", cat.ID)
		}
	}
}

func TestParseSelectedCatalogs(t *testing.T) {
	// Case 1: Empty
	if res := ParseSelectedCatalogs(""); res != nil {
		t.Fatalf("Expected nil for empty input, got %v", res)
	}

	// Case 2: "all"
	if res := ParseSelectedCatalogs("all"); res != nil {
		t.Fatalf("Expected nil for 'all', got %v", res)
	}

	// Case 3: "festivals=cannes-palme-dor,venice-golden-lion"
	res := ParseSelectedCatalogs("festivals=cannes-palme-dor,venice-golden-lion")
	if len(res) != 2 || res[0] != "cannes-palme-dor" || res[1] != "venice-golden-lion" {
		t.Fatalf("Unexpected parse result: %v", res)
	}

	// Case 4: Raw comma-separated
	res = ParseSelectedCatalogs("cannes-palme-dor,berlin-golden-bear")
	if len(res) != 2 || res[0] != "cannes-palme-dor" || res[1] != "berlin-golden-bear" {
		t.Fatalf("Unexpected parse result: %v", res)
	}

	// Case 5: URL-encoded
	res = ParseSelectedCatalogs("cannes-palme-dor%2Cvenice-golden-lion")
	if len(res) != 2 || res[0] != "cannes-palme-dor" || res[1] != "venice-golden-lion" {
		t.Fatalf("Unexpected parse result for encoded string: %v", res)
	}
}

func TestGroupCatalogsByFestival(t *testing.T) {
	groups := GroupCatalogsByFestival()

	totalInGroups := 0
	for _, g := range groups {
		if g.Festival == "" || len(g.Catalogs) == 0 {
			t.Fatalf("Invalid festival group: %+v", g)
		}
		totalInGroups += len(g.Catalogs)
	}

	if totalInGroups != 38 {
		t.Fatalf("Expected 38 catalogs across all festival groups, got %d", totalInGroups)
	}
}
