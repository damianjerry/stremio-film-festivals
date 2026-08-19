package main

import (
	"strings"
	"testing"
)

func TestBuildManifest(t *testing.T) {
	manifest := BuildManifest()

	if manifest.ID == "" {
		t.Fatal("Manifest ID should not be empty")
	}
	if manifest.Name == "" {
		t.Fatal("Manifest Name should not be empty")
	}
	if manifest.Version == "" {
		t.Fatal("Manifest Version should not be empty")
	}
	if len(manifest.Types) == 0 || manifest.Types[0] != "movie" {
		t.Fatalf("Expected Types to include 'movie', got %v", manifest.Types)
	}

	// Check ID prefixes
	hasTT := false
	for _, p := range manifest.IDprefixes {
		if p == "tt" {
			hasTT = true
			break
		}
	}
	if !hasTT {
		t.Fatal("Manifest IDprefixes must include 'tt'")
	}

	// Check Catalogs
	if len(manifest.Catalogs) < 30 {
		t.Fatalf("Expected at least 30 catalogs in manifest, got %d", len(manifest.Catalogs))
	}

	seenIDs := make(map[string]bool)
	for _, cat := range manifest.Catalogs {
		if cat.ID == "" {
			t.Fatal("Catalog ID cannot be empty")
		}
		if seenIDs[cat.ID] {
			t.Fatalf("Duplicate catalog ID found in manifest: %s", cat.ID)
		}
		seenIDs[cat.ID] = true

		if cat.Name == "" {
			t.Fatalf("Catalog %s has empty Name", cat.ID)
		}
		if !strings.Contains(cat.Name, " — ") {
			t.Fatalf("Catalog %s Name does not follow '<Festival> — <Award>' naming convention: %s", cat.ID, cat.Name)
		}
		if cat.Type != "movie" {
			t.Fatalf("Catalog %s has unexpected Type: %s", cat.ID, cat.Type)
		}
	}
}

func TestFindCatalogConfig(t *testing.T) {
	// Test primary ID
	cat, found := FindCatalogConfig("cannes-palme-dor")
	if !found {
		t.Fatal("Expected to find cannes-palme-dor")
	}
	if cat.Festival != "Cannes Film Festival" {
		t.Fatalf("Unexpected festival: %s", cat.Festival)
	}

	// Test alias
	catAlias, found := FindCatalogConfig("palme-dor-winners")
	if !found {
		t.Fatal("Expected to find palme-dor-winners via alias")
	}
	if catAlias.ID != "cannes-palme-dor" {
		t.Fatalf("Expected primary ID cannes-palme-dor, got %s", catAlias.ID)
	}

	// Test non-existent
	_, found = FindCatalogConfig("non-existent-festival")
	if found {
		t.Fatal("Expected not to find non-existent-festival")
	}
}
