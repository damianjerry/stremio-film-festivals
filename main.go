package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deflix-tv/go-stremio"
	"go.uber.org/zap"
)

var (
	bindAddr  = flag.String("bindAddr", getEnv("BIND_ADDR", "0.0.0.0"), `Local interface address to bind to. "localhost" only allows access from the local host. "0.0.0.0" binds to all network interfaces.`)
	port      = flag.Int("port", getEnvInt("PORT", 8080), "Port to listen on")
	dataDir   = flag.String("dataDir", getEnv("DATA_DIR", "data"), `Location of the data directory containing catalog CSV files and optional "metas" subdirectory`)
	logLevel  = flag.String("logLevel", getEnv("LOG_LEVEL", "info"), `Log level: "debug", "info", "warn", "error"`)
	cacheAge  = flag.String("cacheAge", getEnv("CACHE_AGE", "24h"), `Max age for client/proxy caching (e.g. "24h")`)
	orderMode = flag.String("order", getEnv("ORDER_MODE", OrderModeDailyRandom), `Catalog ordering mode: "daily-random" (default), "chronological-desc", "chronological-asc"`)
	redirect  = flag.String("redirect", getEnvAllowEmpty("REDIRECT_URL", defaultRedirectURL), `Where a request for the addon root ("/") is redirected. Relative paths are resolved against the serving host. Empty disables the redirect, and "/" returns 404.`)
)

func init() {
	http.DefaultClient.Timeout = 10 * time.Second
}

func main() {
	flag.Parse()

	// 1. Initialize Logger
	logger, err := stremio.NewLogger(*logLevel)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting Stremio Film Festivals Addon",
		zap.String("version", version),
		zap.String("orderMode", *orderMode),
		zap.String("redirectURL", *redirect),
	)

	// 2. Parse Cache Age Duration
	cacheAgeDuration, err := time.ParseDuration(*cacheAge)
	if err != nil {
		logger.Fatal("Couldn't parse cacheAge duration", zap.Error(err))
	}
	logger.Info("Cache age configured", zap.Duration("duration", cacheAgeDuration))

	// 3. Resolve Data Directory
	resolvedDataDir := resolveDataDir(*dataDir)
	logger.Info("Loading festival catalogs from data directory", zap.String("dataDir", resolvedDataDir))

	// 4. Load Catalogs into Memory
	catalogStore, err := LoadCatalogs(resolvedDataDir, logger)
	if err != nil {
		logger.Fatal("Failed to load festival catalogs", zap.Error(err))
	}
	logger.Info("Successfully loaded festival catalogs",
		zap.Int("primaryCatalogs", catalogStore.TotalCatalogs()),
	)

	// 5. Initialize Discovery / Randomization Engine
	randomizer := NewRandomizer(catalogStore, *orderMode, logger)

	// 6. Build Manifest and Catalog Handlers
	manifest := BuildManifest()
	catalogHandlers := map[string]stremio.CatalogHandler{
		"movie": createMovieHandler(randomizer),
	}

	options := stremio.Options{
		BindAddr:            *bindAddr,
		Port:                *port,
		Logger:              logger,
		RedirectURL:         *redirect,
		CacheAgeCatalogs:    cacheAgeDuration,
		CachePublicCatalogs: true,
		HandleEtagCatalogs:  true,
	}

	addon, err := stremio.NewAddon(manifest, catalogHandlers, nil, options)
	if err != nil {
		logger.Fatal("Couldn't create Stremio addon", zap.Error(err))
	}

	// 7. Register Middleware for Dynamic Manifest Filtering
	addon.AddMiddleware("", HandleManifestMiddleware)

	// 8. Register Configuration UI & Logo Endpoints
	addon.AddEndpoint("GET", "/configure", HandleConfigureEndpoint)
	addon.AddEndpoint("GET", "/:userData/configure", HandleConfigureEndpoint)
	addon.AddEndpoint("GET", "/logo.png", HandleLogoEndpoint)
	addon.AddEndpoint("GET", "/:userData/logo.png", HandleLogoEndpoint)

	logger.Info("Server running and listening for Stremio requests",
		zap.String("bindAddr", *bindAddr),
		zap.Int("port", *port),
	)

	// 9. Run Server (graceful shutdown on SIGINT / SIGTERM)
	addon.Run()
}

func resolveDataDir(dir string) string {
	clean := strings.TrimRight(dir, "/")
	if _, err := os.Stat(clean); err == nil {
		return clean
	}
	// Fallback to "data" if "." was given but files are in "data"
	if _, err := os.Stat("data"); err == nil {
		return "data"
	}
	// Fallback to "/data" in container
	if _, err := os.Stat("/data"); err == nil {
		return "/data"
	}
	return clean
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvAllowEmpty is getEnv for settings where "" is a meaningful value rather
// than shorthand for "unset". REDIRECT_URL needs it: turning the root redirect
// off is done by setting it empty, and the environment is the only way to
// configure a container, so collapsing "" into the default would make that
// documented option unreachable there.
func getEnvAllowEmpty(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
