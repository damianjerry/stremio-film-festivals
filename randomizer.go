package main

import (
	"hash/fnv"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/deflix-tv/go-stremio"
	"go.uber.org/zap"
)

const (
	// OrderModeDailyRandom generates a deterministic randomized discovery permutation rotated daily.
	OrderModeDailyRandom = "daily-random"
	// OrderModeChronologicalDesc orders films from newest year to oldest year.
	OrderModeChronologicalDesc = "chronological-desc"
	// OrderModeChronologicalAsc orders films from oldest year to newest year.
	OrderModeChronologicalAsc = "chronological-asc"
)

// Randomizer provides discovery-oriented ordering and deterministic daily randomization for catalogs.
type Randomizer struct {
	store     *CatalogStore
	orderMode string
	logger    *zap.Logger

	mu         sync.RWMutex
	cacheDate  string
	dailyCache map[string][]stremio.MetaPreviewItem
}

// NewRandomizer creates a new Randomizer with the given CatalogStore and order mode.
func NewRandomizer(store *CatalogStore, orderMode string, logger *zap.Logger) *Randomizer {
	if orderMode == "" {
		orderMode = OrderModeDailyRandom
	}
	return &Randomizer{
		store:      store,
		orderMode:  orderMode,
		logger:     logger,
		dailyCache: make(map[string][]stremio.MetaPreviewItem),
	}
}

// GetCatalogResponse returns the MetaPreviewItems for the requested catalog ID and date.
func (r *Randomizer) GetCatalogResponse(id string, now time.Time) ([]stremio.MetaPreviewItem, error) {
	items, ok := r.store.GetItems(id)
	if !ok {
		return nil, stremio.NotFound
	}
	if len(items) == 0 {
		return []stremio.MetaPreviewItem{}, nil
	}

	switch r.orderMode {
	case OrderModeChronologicalDesc:
		return r.getChronological(items, true), nil
	case OrderModeChronologicalAsc:
		return r.getChronological(items, false), nil
	case OrderModeDailyRandom:
		fallthrough
	default:
		dateStr := now.UTC().Format("2006-01-02")
		return r.getDailyRandom(id, dateStr, items), nil
	}
}

func (r *Randomizer) getDailyRandom(catalogID string, dateStr string, source []stremio.MetaPreviewItem) []stremio.MetaPreviewItem {
	r.mu.RLock()
	if r.cacheDate == dateStr {
		if cached, ok := r.dailyCache[catalogID]; ok {
			r.mu.RUnlock()
			return cached
		}
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Re-check after acquiring write lock
	if r.cacheDate != dateStr {
		r.cacheDate = dateStr
		r.dailyCache = make(map[string][]stremio.MetaPreviewItem)
	} else if cached, ok := r.dailyCache[catalogID]; ok {
		return cached
	}

	// Deterministic permutation using FNV-1a hash of (catalogID + ":" + dateStr)
	seed := hashCatalogDate(catalogID, dateStr)
	rng := rand.New(rand.NewSource(seed))

	shuffled := make([]stremio.MetaPreviewItem, len(source))
	copy(shuffled, source)

	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	r.dailyCache[catalogID] = shuffled
	return shuffled
}

func (r *Randomizer) getChronological(source []stremio.MetaPreviewItem, desc bool) []stremio.MetaPreviewItem {
	result := make([]stremio.MetaPreviewItem, len(source))
	copy(result, source)

	sort.SliceStable(result, func(i, j int) bool {
		yearI, _ := strconv.Atoi(result[i].ReleaseInfo)
		yearJ, _ := strconv.Atoi(result[j].ReleaseInfo)
		if desc {
			return yearI > yearJ
		}
		return yearI < yearJ
	})

	return result
}

func hashCatalogDate(catalogID, dateStr string) int64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(catalogID))
	hasher.Write([]byte(":"))
	hasher.Write([]byte(dateStr))
	return int64(hasher.Sum64())
}
