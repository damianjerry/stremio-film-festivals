package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var imdbIDRegex = regexp.MustCompile(`^tt\d{7,8}$`)

func TestAllCSVDataIntegrity(t *testing.T) {
	files, err := filepath.Glob("data/*.csv")
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}

	if len(files) < 35 {
		t.Fatalf("Expected at least 35 CSV files in data/, found %d", len(files))
	}

	totalRecords := 0

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("Could not open file: %v", err)
			}
			defer f.Close()

			reader := csv.NewReader(f)
			rows, err := reader.ReadAll()
			if err != nil {
				t.Fatalf("CSV read error: %v", err)
			}

			if len(rows) < 2 {
				t.Fatalf("CSV %s must have at least 1 data row, has %d rows", file, len(rows))
			}

			header := rows[0]
			if len(header) != 3 || header[0] != "year" || header[1] != "title" || header[2] != "IMDb ID" {
				t.Fatalf("Invalid header in %s: expected [year, title, IMDb ID], got %v", file, header)
			}

			seenKeys := make(map[string]bool)
			for rowIdx, row := range rows[1:] {
				if len(row) != 3 {
					t.Fatalf("%s row %d: expected 3 columns, got %d", file, rowIdx+2, len(row))
				}

				yearStr := strings.TrimSpace(row[0])
				titleStr := strings.TrimSpace(row[1])
				imdbStr := strings.TrimSpace(row[2])

				// Validate year
				year, err := strconv.Atoi(yearStr)
				if err != nil || year < 1920 || year > 2030 {
					t.Fatalf("%s row %d: invalid year %s", file, rowIdx+2, yearStr)
				}

				// Validate title
				if titleStr == "" {
					t.Fatalf("%s row %d: empty title", file, rowIdx+2)
				}

				// Validate IMDb ID
				if !imdbIDRegex.MatchString(imdbStr) {
					t.Fatalf("%s row %d: invalid IMDb ID format '%s'", file, rowIdx+2, imdbStr)
				}

				// Validate no duplicate (year, IMDb ID) in same catalog
				key := yearStr + ":" + imdbStr
				if seenKeys[key] {
					t.Fatalf("%s row %d: duplicate key '%s' (%s)", file, rowIdx+2, key, titleStr)
				}
				seenKeys[key] = true
			}

			totalRecords += len(rows) - 1
		})
	}

	t.Logf("Validated %d total film festival winner records across %d CSV files", totalRecords, len(files))
}
