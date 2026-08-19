package main

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ValidateCandidateFilm verifies that an IMDb candidate match satisfies multi-signal validation:
// 1. Title similarity >= minSim
// 2. Year compatibility within maxYearDiff
// 3. Not an explanatory note / cancellation text
// 4. Not a known person name
func ValidateCandidateFilm(expectedTitle string, expectedYear int, candidateTitle string, candidateYear int, minSim float64, maxYearDiff int) bool {
	cleanExpected := strings.ToLower(strings.TrimSpace(expectedTitle))
	cleanCandidate := strings.ToLower(strings.TrimSpace(candidateTitle))

	// Rejection rule 1: Explanatory note / cancellation text
	noteKeywords := []string{"outbreak", "second world war", "cancelled", "no festival", "not held", "no award", "timeline of"}
	for _, kw := range noteKeywords {
		if strings.Contains(cleanExpected, kw) {
			return false
		}
	}

	// Rejection rule 2: Known person names extracted as titles
	personNames := []string{"tricia tuttle", "richard roud", "michael verhoeven", "gérard depardieu", "jerry lewis"}
	for _, pn := range personNames {
		if cleanExpected == pn {
			return false
		}
	}

	// Rejection rule 3: Year disparity
	yearDiff := int(math.Abs(float64(expectedYear - candidateYear)))
	if yearDiff > maxYearDiff {
		return false
	}

	// Rejection rule 4: Title similarity
	sim := calculateJaccardSimilarity(cleanExpected, cleanCandidate)
	if sim < minSim {
		return false
	}

	return true
}

func calculateJaccardSimilarity(s1, s2 string) float64 {
	words1 := strings.Fields(regexp.MustCompile(`[^\w\s]`).ReplaceAllString(s1, ""))
	words2 := strings.Fields(regexp.MustCompile(`[^\w\s]`).ReplaceAllString(s2, ""))

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	set1 := make(map[string]bool)
	for _, w := range words1 {
		set1[w] = true
	}

	intersection := 0
	unionMap := make(map[string]bool)
	for _, w := range words1 {
		unionMap[w] = true
	}
	for _, w := range words2 {
		unionMap[w] = true
		if set1[w] {
			intersection++
		}
	}

	return float64(intersection) / float64(len(unionMap))
}

// TestRegressionDiscoveredFailureModes tests multi-signal resolution rules
func TestRegressionDiscoveredFailureModes(t *testing.T) {
	// Case 1: Known failure mode - "outbreak of the Second World War" (1939) -> "Warrior of the Lost World" (1983)
	// Must be rejected!
	if ValidateCandidateFilm("outbreak of the Second World War", 1939, "Warrior of the Lost World", 1983, 0.5, 2) {
		t.Fatal("Regression failure: 'outbreak of the Second World War' matching 'Warrior of the Lost World' should have been REJECTED")
	}

	// Case 2: Title mismatch - "Casablanca" (1943) matching "You Must Remember This: A Tribute to 'Casablanca'" (1992)
	// Must be rejected due to year mismatch (1943 vs 1992)
	if ValidateCandidateFilm("Casablanca", 1943, "You Must Remember This: A Tribute to 'Casablanca'", 1992, 0.5, 2) {
		t.Fatal("Regression failure: 'Casablanca' matching 1992 documentary tribute should have been REJECTED")
	}

	// Case 3: Correct title + correct year - "Casablanca" (1943) matching "Casablanca" (1942)
	// Must be accepted!
	if !ValidateCandidateFilm("Casablanca", 1943, "Casablanca", 1942, 0.5, 2) {
		t.Fatal("Validation failure: 'Casablanca' matching 'Casablanca' (1942) should be ACCEPTED")
	}

	// Case 4: Person name extracted as title - "Tricia Tuttle" (2018)
	// Must be rejected!
	if ValidateCandidateFilm("Tricia Tuttle", 2018, "Primitive War", 2025, 0.5, 2) {
		t.Fatal("Regression failure: 'Tricia Tuttle' person name should have been REJECTED")
	}

	// Case 5: Year mismatch > 2 years for standard release
	// "Anatomy of a Fall" (2023) matching a 2010 film with similar title should be rejected
	if ValidateCandidateFilm("Anatomy of a Fall", 2023, "Anatomy of a Fall", 2010, 0.5, 2) {
		t.Fatal("Regression failure: Large year disparity (2023 vs 2010) should have been REJECTED")
	}
}

// TestRegressionNoExplanatoryNotesInCSVs scans all CSVs in data/ to assert 0 note rows exist
func TestRegressionNoExplanatoryNotesInCSVs(t *testing.T) {
	files, err := filepath.Glob("data/*.csv")
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}

	bannedPhrases := []string{
        "outbreak of the second world war",
		"no festival held",
		"festival cancelled",
		"no award given",
		"not awarded",
		"tricia tuttle",
		"richard roud",
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("Could not open %s: %v", file, err)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		rows, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("CSV read error in %s: %v", file, err)
		}

		for rowIdx, row := range rows[1:] {
			titleLower := strings.ToLower(strings.TrimSpace(row[1]))
			for _, banned := range bannedPhrases {
				if strings.Contains(titleLower, banned) {
					t.Fatalf("Regression in %s row %d: found banned phrase/entity '%s' in title '%s'", file, rowIdx+2, banned, row[1])
				}
			}
		}
	}
}

// TestRegressionCanonicalLandmarkMappings verifies that landmark classics have exact canonical IMDb IDs
func TestRegressionCanonicalLandmarkMappings(t *testing.T) {
	landmarkMap := map[string]string{
		"cannes-palme-dor:1939":           "tt0032080", // Union Pacific
		"academy-awards-best-picture:1943": "tt0034583", // Casablanca
		"academy-awards-best-picture:1939": "tt0031381", // Gone with the Wind
		"academy-awards-best-picture:1950": "tt0042192", // All About Eve
		"academy-awards-best-picture:1972": "tt0068646", // The Godfather
		"cannes-palme-dor:1976":           "tt0075314", // Taxi Driver
		"cannes-palme-dor:1979":           "tt0078788", // Apocalypse Now
		"cannes-palme-dor:1994":           "tt0110912", // Pulp Fiction
		"cannes-palme-dor:2019":           "tt6751668", // Parasite
		"cannes-palme-dor:2022":           "tt7322224", // Triangle of Sadness
		"cannes-palme-dor:2023":           "tt17009710",// Anatomy of a Fall
		"cannes-palme-dor:2024":           "tt28607951",// Anora
		"cannes-palme-dor:1973":           "tt0070643", // Scarecrow
		"cannes-grand-prix:2022":          "tt10354106",// Stars at Noon
		"cannes-best-actress:2008":        "tt0803029", // Linha de Passe
	}

	for key, expectedID := range landmarkMap {
		parts := strings.Split(key, ":")
		catalogID := parts[0]
		targetYear := parts[1]

		csvPath := filepath.Join("data", catalogID+".csv")
		f, err := os.Open(csvPath)
		if err != nil {
			t.Fatalf("Could not open %s: %v", csvPath, err)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		rows, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("CSV read error in %s: %v", csvPath, err)
		}

		found := false
		for _, row := range rows[1:] {
			if strings.TrimSpace(row[0]) == targetYear {
				if strings.TrimSpace(row[2]) == expectedID {
					found = true
					break
				}
			}
		}

		if !found {
			t.Fatalf("Landmark entry %s (expected ID %s) not found among records for year %s in %s", key, expectedID, targetYear, csvPath)
		}
	}
}
