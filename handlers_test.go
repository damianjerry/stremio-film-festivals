package main

import (
	"testing"

	"github.com/deflix-tv/go-stremio"
	"go.uber.org/zap"
)

func TestMovieHandler(t *testing.T) {
	logger := zap.NewNop()
	store, err := LoadCatalogs("data", logger)
	if err != nil {
		t.Fatalf("LoadCatalogs failed: %v", err)
	}

	randomizer := NewRandomizer(store, OrderModeDailyRandom, logger)
	handler := createMovieHandler(randomizer)

	// Test all primary catalogs
	for _, cat := range FestivalCatalogs {
		items, err := handler(cat.ID, nil)
		if err != nil {
			t.Fatalf("Handler failed for catalog %s: %v", cat.ID, err)
		}
		if len(items) == 0 {
			t.Fatalf("Handler returned 0 items for catalog %s", cat.ID)
		}
	}

	// Test legacy aliases
	aliases := []string{
		"palme-dor-winners",
		"golden-lion-winners",
		"golden-bear-winners",
		"academy-awards-winners",
	}

	for _, alias := range aliases {
		items, err := handler(alias, nil)
		if err != nil {
			t.Fatalf("Handler failed for alias %s: %v", alias, err)
		}
		if len(items) == 0 {
			t.Fatalf("Handler returned 0 items for alias %s", alias)
		}
	}

	// Test unknown catalog ID
	_, err = handler("unknown-festival-catalog", nil)
	if err != stremio.NotFound {
		t.Fatalf("Expected stremio.NotFound for unknown catalog, got %v", err)
	}
}
