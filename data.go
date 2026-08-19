package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/deflix-tv/go-stremio"
	"go.uber.org/zap"
)

var imdbIDSanitizer = regexp.MustCompile(`[^a-zA-Z0-9]`)

// FilmRecord represents a single verified festival winning film entry.
type FilmRecord struct {
	Year   string
	Title  string
	IMDbID string
}

// CatalogStore holds in-memory film records and MetaPreviewItems for each catalog.
type CatalogStore struct {
	catalogs map[string][]stremio.MetaPreviewItem
	aliases  map[string]string
	logger   *zap.Logger
}

// LoadCatalogs loads all catalog CSV files from dataDir and prepares preview items.
func LoadCatalogs(dataDir string, logger *zap.Logger) (*CatalogStore, error) {
	cleanDir := strings.TrimRight(dataDir, "/")
	store := &CatalogStore{
		catalogs: make(map[string][]stremio.MetaPreviewItem),
		aliases:  make(map[string]string),
		logger:   logger,
	}

	metasDir := filepath.Join(cleanDir, "metas")

	for _, cat := range FestivalCatalogs {
		csvPath := filepath.Join(cleanDir, cat.CSVFile)
		records, err := readCatalogCSV(csvPath)
		if err != nil {
			logger.Warn("Could not load CSV file for catalog",
				zap.String("catalogID", cat.ID),
				zap.String("csvFile", csvPath),
				zap.Error(err),
			)
			continue
		}

		items := buildMetaPreviewItems(records, metasDir, logger)
		store.catalogs[cat.ID] = items
		logger.Debug("Loaded catalog",
			zap.String("id", cat.ID),
			zap.Int("items", len(items)),
		)

		// Register aliases
		for _, alias := range cat.Aliases {
			store.aliases[alias] = cat.ID
		}
	}

	return store, nil
}

// GetItems returns the raw chronological items for a catalog ID or alias.
func (cs *CatalogStore) GetItems(id string) ([]stremio.MetaPreviewItem, bool) {
	// Check primary ID
	if items, ok := cs.catalogs[id]; ok {
		return items, true
	}
	// Check alias
	if primaryID, ok := cs.aliases[id]; ok {
		if items, ok := cs.catalogs[primaryID]; ok {
			return items, true
		}
	}
	return nil, false
}

// HasCatalog returns true if the catalog ID or alias exists.
func (cs *CatalogStore) HasCatalog(id string) bool {
	if _, ok := cs.catalogs[id]; ok {
		return true
	}
	if _, ok := cs.aliases[id]; ok {
		return true
	}
	return false
}

// TotalCatalogs returns the number of primary catalogs loaded.
func (cs *CatalogStore) TotalCatalogs() int {
	return len(cs.catalogs)
}

func readCatalogCSV(filePath string) ([]FilmRecord, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading CSV %s: %w", filePath, err)
	}

	csvReader := csv.NewReader(bytes.NewReader(fileBytes))
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV %s: %w", filePath, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty CSV %s", filePath)
	}

	header := rows[0]
	yearIdx, titleIdx, imdbIdx := -1, -1, -1
	for i, col := range header {
		cleanCol := strings.Trim(strings.ToLower(col), " \r\n\t")
		switch cleanCol {
		case "year":
			yearIdx = i
		case "title":
			titleIdx = i
		case "imdb id", "imdbid", "imdb_id", "id":
			imdbIdx = i
		}
	}

	if yearIdx == -1 || titleIdx == -1 || imdbIdx == -1 {
		return nil, fmt.Errorf("invalid header in %s: expected year,title,IMDb ID, got %v", filePath, header)
	}

	var records []FilmRecord
	for _, row := range rows[1:] {
		if len(row) <= yearIdx || len(row) <= titleIdx || len(row) <= imdbIdx {
			continue
		}
		rawIMDb := strings.Trim(row[imdbIdx], " \r\n\t")
		cleanIMDb := imdbIDSanitizer.ReplaceAllString(rawIMDb, "")
		if cleanIMDb == "" || !strings.HasPrefix(cleanIMDb, "tt") {
			continue
		}

		cleanYear := strings.Trim(row[yearIdx], " \r\n\t")
		cleanTitle := strings.Trim(row[titleIdx], " \r\n\t")

		records = append(records, FilmRecord{
			Year:   cleanYear,
			Title:  cleanTitle,
			IMDbID: cleanIMDb,
		})
	}

	return records, nil
}

func buildMetaPreviewItems(records []FilmRecord, metasDir string, logger *zap.Logger) []stremio.MetaPreviewItem {
	items := make([]stremio.MetaPreviewItem, 0, len(records))

	for _, rec := range records {
		metaFile := filepath.Join(metasDir, rec.IMDbID+".json")
		if metaBytes, err := ioutil.ReadFile(metaFile); err == nil {
			var item stremio.MetaPreviewItem
			if err := json.Unmarshal(metaBytes, &item); err == nil && item.ID != "" {
				item.ID = imdbIDSanitizer.ReplaceAllString(strings.Trim(item.ID, " \r\n\t"), "")
				if item.Name == "" {
					item.Name = rec.Title
				}
				if item.ReleaseInfo == "" {
					item.ReleaseInfo = rec.Year
				}
				if item.Poster == "" {
					item.Poster = fmt.Sprintf("https://images.metahub.space/poster/medium/%s/img.jpg", rec.IMDbID)
				}
				items = append(items, item)
				continue
			}
		}

		// Fallback: construct standard Stremio MetaPreviewItem
		items = append(items, stremio.MetaPreviewItem{
			ID:          rec.IMDbID,
			Type:        "movie",
			Name:        rec.Title,
			Poster:      fmt.Sprintf("https://images.metahub.space/poster/medium/%s/img.jpg", rec.IMDbID),
			ReleaseInfo: rec.Year,
		})
	}

	return items
}

// EnsureDataDirExists checks whether the data directory exists.
func EnsureDataDirExists(dataDir string) error {
	info, err := os.Stat(dataDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dataDir)
	}
	return nil
}
