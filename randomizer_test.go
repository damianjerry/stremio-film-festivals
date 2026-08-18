package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/deflix-tv/go-stremio"
	"go.uber.org/zap"
)

func TestRandomizerDeterminism(t *testing.T) {
	logger := zap.NewNop()
	store, err := LoadCatalogs("data", logger)
	if err != nil {
		t.Fatalf("LoadCatalogs failed: %v", err)
	}

	randomizer := NewRandomizer(store, OrderModeDailyRandom, logger)
	date1 := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	date1Later := time.Date(2026, 8, 19, 22, 30, 0, 0, time.UTC)

	catalogID := "cannes-palme-dor"

	resp1, err := randomizer.GetCatalogResponse(catalogID, date1)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	resp2, err := randomizer.GetCatalogResponse(catalogID, date1Later)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	if len(resp1) != len(resp2) {
		t.Fatalf("Expected same length on same date, got %d vs %d", len(resp1), len(resp2))
	}

	for i := range resp1 {
		if resp1[i].ID != resp2[i].ID {
			t.Fatalf("Discrepancy at index %d on same date: %s vs %s", i, resp1[i].ID, resp2[i].ID)
		}
	}
}

func TestRandomizerDailyRotation(t *testing.T) {
	logger := zap.NewNop()
	store, err := LoadCatalogs("data", logger)
	if err != nil {
		t.Fatalf("LoadCatalogs failed: %v", err)
	}

	randomizer := NewRandomizer(store, OrderModeDailyRandom, logger)
	dateDay1 := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	dateDay2 := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)

	catalogID := "cannes-palme-dor"

	respDay1, err := randomizer.GetCatalogResponse(catalogID, dateDay1)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	respDay2, err := randomizer.GetCatalogResponse(catalogID, dateDay2)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	if len(respDay1) != len(respDay2) {
		t.Fatalf("Expected same length across days, got %d vs %d", len(respDay1), len(respDay2))
	}

	// Verify all items are preserved
	set1 := make(map[string]bool)
	for _, item := range respDay1 {
		set1[item.ID] = true
	}

	for _, item := range respDay2 {
		if !set1[item.ID] {
			t.Fatalf("Item %s missing in Day 2 permutation", item.ID)
		}
	}

	// Verify the permutation is different (not identical ordering)
	identicalOrder := true
	for i := range respDay1 {
		if respDay1[i].ID != respDay2[i].ID {
			identicalOrder = false
			break
		}
	}

	if identicalOrder {
		t.Fatal("Expected different permutations across consecutive days, but ordering was identical")
	}
}

func TestRandomizerChronologicalModes(t *testing.T) {
	logger := zap.NewNop()
	store, err := LoadCatalogs("data", logger)
	if err != nil {
		t.Fatalf("LoadCatalogs failed: %v", err)
	}

	now := time.Now().UTC()
	catalogID := "academy-awards-best-picture"

	// Test Descending
	randomizerDesc := NewRandomizer(store, OrderModeChronologicalDesc, logger)
	respDesc, err := randomizerDesc.GetCatalogResponse(catalogID, now)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	for i := 0; i < len(respDesc)-1; i++ {
		yrCurr, _ := strconv.Atoi(respDesc[i].ReleaseInfo)
		yrNext, _ := strconv.Atoi(respDesc[i+1].ReleaseInfo)
		if yrCurr < yrNext {
			t.Fatalf("Descending order violated at index %d: %d < %d", i, yrCurr, yrNext)
		}
	}

	// Test Ascending
	randomizerAsc := NewRandomizer(store, OrderModeChronologicalAsc, logger)
	respAsc, err := randomizerAsc.GetCatalogResponse(catalogID, now)
	if err != nil {
		t.Fatalf("GetCatalogResponse failed: %v", err)
	}

	for i := 0; i < len(respAsc)-1; i++ {
		yrCurr, _ := strconv.Atoi(respAsc[i].ReleaseInfo)
		yrNext, _ := strconv.Atoi(respAsc[i+1].ReleaseInfo)
		if yrCurr > yrNext {
			t.Fatalf("Ascending order violated at index %d: %d > %d", i, yrCurr, yrNext)
		}
	}
}

func TestRandomizerNotFound(t *testing.T) {
	logger := zap.NewNop()
	store, _ := LoadCatalogs("data", logger)
	randomizer := NewRandomizer(store, OrderModeDailyRandom, logger)

	_, err := randomizer.GetCatalogResponse("invalid-catalog-id", time.Now().UTC())
	if err != stremio.NotFound {
		t.Fatalf("Expected stremio.NotFound, got %v", err)
	}
}
