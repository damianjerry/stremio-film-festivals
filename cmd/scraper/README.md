# Film Festival Data Scrapers & Generator

This directory and the `scripts/` directory contain utilities for fetching, extracting, and verifying film festival winner datasets.

## Dataset Generator (`scripts/build_catalogs.py`)

The primary dataset generation and verification engine is `scripts/build_catalogs.py`.

It:
1. Queries Wikipedia APIs and Wikidata properties (`P345` IMDb ID) for authoritative winner lists.
2. Cross-references against Cinemeta search API for title, release year, and IMDb ID resolution.
3. Cleans honorary awards, non-competitive records, and ensures only genuine winning feature films are included.
4. Caches resolved IMDb IDs in `scripts/.imdb_cache.json` for fast reproducible execution.
5. Writes validated CSV files directly into `data/`.

### Running the Generator
```bash
# Using uv:
uv run --with beautifulsoup4,requests python3 scripts/build_catalogs.py

# Or using standard python3 (after pip install beautifulsoup4 requests):
python3 scripts/build_catalogs.py
```

## Go Scraper (`cmd/scraper/main.go`)

The historical Go scraper using `goquery` is retained as an alternative example for direct scraping of IMDb/Wikipedia event pages.
