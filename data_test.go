package main

import (
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestLoadCatalogs(t *testing.T) {
	logger := zap.NewNop()
	store, err := LoadCatalogs("data", logger)
	if err != nil {
		t.Fatalf("LoadCatalogs failed: %v", err)
	}

	if store.TotalCatalogs() < 30 {
		t.Fatalf("Expected at least 30 catalogs loaded, got %d", store.TotalCatalogs())
	}

	for _, cat := range FestivalCatalogs {
		items, ok := store.GetItems(cat.ID)
		if !ok {
			t.Fatalf("Catalog %s was not loaded", cat.ID)
		}
		if len(items) == 0 {
			t.Fatalf("Catalog %s has 0 items loaded", cat.ID)
		}

		for idx, item := range items {
			if item.ID == "" || !strings.HasPrefix(item.ID, "tt") {
				t.Fatalf("Catalog %s item %d has invalid IMDb ID: %s", cat.ID, idx, item.ID)
			}
			if item.Name == "" {
				t.Fatalf("Catalog %s item %d has empty Name", cat.ID, idx)
			}
			if item.Type != "movie" {
				t.Fatalf("Catalog %s item %d has invalid Type: %s", cat.ID, idx, item.Type)
			}
			if item.Poster == "" {
				t.Fatalf("Catalog %s item %d has empty Poster URL", cat.ID, idx)
			}
			if item.ReleaseInfo == "" {
				t.Fatalf("Catalog %s item %d has empty ReleaseInfo", cat.ID, idx)
			}
		}
	}

	// Test aliases
	legacyAliases := []string{
		"palme-dor-winners",
		"golden-lion-winners",
		"golden-bear-winners",
		"academy-awards-winners",
	}

	for _, alias := range legacyAliases {
		items, ok := store.GetItems(alias)
		if !ok {
			t.Fatalf("Expected alias %s to return items", alias)
		}
		if len(items) == 0 {
			t.Fatalf("Alias %s returned empty items", alias)
		}
	}

	// Test non-existent
	if store.HasCatalog("non-existent-catalog") {
		t.Fatal("Expected HasCatalog to return false for non-existent catalog")
	}
}
