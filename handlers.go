package main

import (
	"time"

	"github.com/deflix-tv/go-stremio"
)

// createMovieHandler creates a Stremio CatalogHandler that routes catalog requests through the Randomizer.
func createMovieHandler(randomizer *Randomizer) stremio.CatalogHandler {
	return func(id string, userData interface{}) ([]stremio.MetaPreviewItem, error) {
		return randomizer.GetCatalogResponse(id, time.Now().UTC())
	}
}
