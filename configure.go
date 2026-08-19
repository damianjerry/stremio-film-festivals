package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/deflix-tv/go-stremio"
	"github.com/gofiber/fiber"
)

//go:embed "film festivals logo.png"
var logoBytes []byte

// CustomManifest represents the manifest JSON with configurable behavior hints
type CustomManifest struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Version       string                 `json:"version"`
	ResourceItems []stremio.ResourceItem `json:"resources,omitempty"`
	Types         []string               `json:"types"`
	Catalogs      []stremio.CatalogItem  `json:"catalogs"`
	IDprefixes    []string               `json:"idPrefixes,omitempty"`
	Background    string                 `json:"background,omitempty"`
	Logo          string                 `json:"logo,omitempty"`
	ContactEmail  string                 `json:"contactEmail,omitempty"`
	BehaviorHints struct {
		Configurable          bool `json:"configurable"`
		ConfigurationRequired bool `json:"configurationRequired"`
	} `json:"behaviorHints"`
}

// BuildCustomManifest constructs a CustomManifest, optionally filtered to specified catalog IDs.
func BuildCustomManifest(selectedCatalogIDs []string, baseURL string) CustomManifest {
	var selectedMap map[string]bool
	if len(selectedCatalogIDs) > 0 {
		selectedMap = make(map[string]bool, len(selectedCatalogIDs))
		for _, id := range selectedCatalogIDs {
			cleanID := strings.TrimSpace(id)
			if cleanID != "" {
				selectedMap[cleanID] = true
			}
		}
	}

	catalogs := make([]stremio.CatalogItem, 0, len(FestivalCatalogs))
	for _, cat := range FestivalCatalogs {
		if selectedMap == nil || selectedMap[cat.ID] {
			catalogs = append(catalogs, stremio.CatalogItem{
				Type: "movie",
				ID:   cat.ID,
				Name: cat.Name,
			})
		}
	}

	logoURL := "/logo.png"
	if baseURL != "" {
		logoURL = strings.TrimRight(baseURL, "/") + "/logo.png"
	}

	cm := CustomManifest{
		ID:          "tv.deflix.stremio-film-festivals",
		Name:        "Film Festivals",
		Description: "Discover arthouse, auteur, and award-winning cinema from Cannes, Venice, Berlinale, Locarno, Sundance, TIFF, Rotterdam, San Sebastián, Karlovy Vary, BFI London, IDFA, SXSW, FIPRESCI & the Oscars.",
		Version:     version,
		ResourceItems: []stremio.ResourceItem{
			{
				Name: "catalog",
			},
		},
		Types:      []string{"movie"},
		Catalogs:   catalogs,
		IDprefixes: []string{"tt"},
		Background: "https://images.metahub.space/background/medium/tt0078788/img",
		Logo:       logoURL,
	}

	cm.BehaviorHints.Configurable = true
	cm.BehaviorHints.ConfigurationRequired = false

	return cm
}

// ParseSelectedCatalogs extracts catalog IDs from query param or path segment.
func ParseSelectedCatalogs(input string) []string {
	if input == "" || input == "all" {
		return nil
	}

	// Remove leading "festivals=" if present
	if strings.HasPrefix(input, "festivals=") {
		input = strings.TrimPrefix(input, "festivals=")
	}

	// URL decode
	decoded, err := url.QueryUnescape(input)
	if err == nil {
		input = decoded
	}

	parts := strings.Split(input, ",")
	var result []string
	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if clean != "" {
			result = append(result, clean)
		}
	}
	return result
}

// HandleLogoEndpoint serves the embedded PNG logo.
func HandleLogoEndpoint(c *fiber.Ctx) {
	c.Set(fiber.HeaderContentType, "image/png")
	c.Set(fiber.HeaderCacheControl, "public, max-age=604800") // 7 days cache
	c.Send(logoBytes)
}

// HandleManifestMiddleware intercepts manifest requests to allow dynamic filtering and configurable behavior hints.
func HandleManifestMiddleware(c *fiber.Ctx) {
	path := strings.Trim(c.Path(), "/")
	if !strings.HasSuffix(path, "manifest.json") {
		c.Next()
		return
	}

	var selectedIDs []string
	// 1. Check query param: ?festivals=cannes-palme-dor,venice-golden-lion
	if q := c.Query("festivals"); q != "" {
		selectedIDs = ParseSelectedCatalogs(q)
	} else if path != "manifest.json" {
		// 2. Check path parameter: /:userData/manifest.json
		prefix := strings.TrimSuffix(path, "/manifest.json")
		if prefix != "" && prefix != "manifest.json" {
			selectedIDs = ParseSelectedCatalogs(prefix)
		}
	}

	// Detect Base URL
	protocol := "http"
	if c.Secure() || c.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	host := c.Hostname()
	baseURL := fmt.Sprintf("%s://%s", protocol, host)

	manifest := BuildCustomManifest(selectedIDs, baseURL)
	c.Set(fiber.HeaderContentType, "application/json; charset=utf-8")
	c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
	c.Set(fiber.HeaderCacheControl, "no-cache")

	resBytes, err := json.Marshal(manifest)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error marshaling manifest")
		return
	}
	c.Send(resBytes)
}

// FestivalGroup represents a grouping of catalogs by festival for the HTML UI.
type FestivalGroup struct {
	Festival string
	Flag     string
	Catalogs []CatalogConfig
}

// GroupCatalogsByFestival groups all 38 catalogs by festival for display.
func GroupCatalogsByFestival() []FestivalGroup {
	festMap := make(map[string][]CatalogConfig)
	var order []string

	for _, cat := range FestivalCatalogs {
		if _, exists := festMap[cat.Festival]; !exists {
			order = append(order, cat.Festival)
		}
		festMap[cat.Festival] = append(festMap[cat.Festival], cat)
	}

	flags := map[string]string{
		"Cannes Film Festival":                              "🇫🇷",
		"Venice Film Festival":                              "🇮🇹",
		"Berlin International Film Festival":                "🇩🇪",
		"Locarno Film Festival":                             "🇨🇭",
		"Sundance Film Festival":                            "🇺🇸",
		"Toronto International Film Festival":               "🇨🇦",
		"International Film Festival Rotterdam":             "🇳🇱",
		"San Sebastián International Film Festival":         "🇪🇸",
		"Karlovy Vary International Film Festival":          "🇨🇿",
		"BFI London Film Festival":                          "🇬🇧",
		"International Documentary Film Festival Amsterdam": "🇳🇱",
		"SXSW Film & TV Festival":                           "🇺🇸",
		"FIPRESCI":                                          "🌍",
		"Academy Awards":                                    "🇺🇸",
	}

	groups := make([]FestivalGroup, 0, len(order))
	for _, fest := range order {
		flag := flags[fest]
		if flag == "" {
			flag = "🎬"
		}
		groups = append(groups, FestivalGroup{
			Festival: fest,
			Flag:     flag,
			Catalogs: festMap[fest],
		})
	}
	return groups
}

// HandleConfigureEndpoint serves the HTML configuration UI.
func HandleConfigureEndpoint(c *fiber.Ctx) {
	groups := GroupCatalogsByFestival()

	protocol := "http"
	if c.Secure() || c.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	host := c.Hostname()
	currentHost := fmt.Sprintf("%s://%s", protocol, host)

	var builder strings.Builder
	for _, g := range groups {
		builder.WriteString(fmt.Sprintf(`<div class="group-card">
			<div class="group-header">
				<span class="group-flag">%s</span>
				<h3>%s</h3>
				<span class="group-count">%d catalogs</span>
			</div>
			<div class="checkbox-grid">`, g.Flag, g.Festival, len(g.Catalogs)))

		for _, cat := range g.Catalogs {
			builder.WriteString(fmt.Sprintf(`
				<label class="catalog-item" title="%s">
					<input type="checkbox" name="festival" value="%s" checked onchange="updateManifestURL()">
					<div class="catalog-info">
						<span class="cat-award">%s</span>
						<span class="cat-desc">%s</span>
					</div>
				</label>`, cat.Description, cat.ID, cat.Award, cat.Description))
		}

		builder.WriteString(`</div></div>`)
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Film Festivals — Stremio Addon Configuration</title>
	<link rel="icon" type="image/png" href="/logo.png">
	<style>
		:root {
			--bg-base: #0f1117;
			--bg-surface: #1a1d26;
			--bg-card: #232734;
			--bg-card-hover: #2b3040;
			--accent: #8e44ad;
			--accent-hover: #9b59b6;
			--accent-gold: #f39c12;
			--text-main: #f0f2f5;
			--text-muted: #9ba1b0;
			--border: rgba(255, 255, 255, 0.08);
			--radius: 12px;
		}

		* {
			box-sizing: border-box;
			margin: 0;
			padding: 0;
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
		}

		body {
			background: var(--bg-base);
			color: var(--text-main);
			min-height: 100vh;
			display: flex;
			flex-direction: column;
			align-items: center;
			padding: 40px 20px;
		}

		.container {
			max-width: 900px;
			width: 100%%;
		}

		header {
			text-align: center;
			margin-bottom: 32px;
		}

		.logo {
			width: 96px;
			height: 96px;
			border-radius: 20px;
			box-shadow: 0 10px 30px rgba(142, 68, 173, 0.35);
			margin-bottom: 16px;
		}

		h1 {
			font-size: 2.2rem;
			font-weight: 800;
			letter-spacing: -0.5px;
			margin-bottom: 8px;
		}

		.subtitle {
			color: var(--text-muted);
			font-size: 1.05rem;
			max-width: 600px;
			margin: 0 auto 20px auto;
			line-height: 1.5;
		}

		.toolbar {
			background: var(--bg-surface);
			border: 1px solid var(--border);
			border-radius: var(--radius);
			padding: 16px 20px;
			display: flex;
			flex-wrap: wrap;
			gap: 12px;
			align-items: center;
			justify-content: space-between;
			margin-bottom: 24px;
			position: sticky;
			top: 16px;
			z-index: 100;
			backdrop-filter: blur(12px);
			box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
		}

		.quick-filters {
			display: flex;
			gap: 8px;
			flex-wrap: wrap;
		}

		.btn-filter {
			background: var(--bg-card);
			border: 1px solid var(--border);
			color: var(--text-main);
			padding: 8px 14px;
			border-radius: 8px;
			font-size: 0.85rem;
			font-weight: 600;
			cursor: pointer;
			transition: all 0.2s ease;
		}

		.btn-filter:hover {
			background: var(--bg-card-hover);
			border-color: rgba(255, 255, 255, 0.2);
		}

		.btn-filter.active {
			background: var(--accent);
			border-color: var(--accent-hover);
		}

		.selection-badge {
			font-size: 0.9rem;
			font-weight: 700;
			color: var(--accent-gold);
			background: rgba(243, 156, 18, 0.12);
			padding: 6px 14px;
			border-radius: 20px;
			border: 1px solid rgba(243, 156, 18, 0.25);
		}

		.groups {
			display: flex;
			flex-direction: column;
			gap: 20px;
			margin-bottom: 32px;
		}

		.group-card {
			background: var(--bg-surface);
			border: 1px solid var(--border);
			border-radius: var(--radius);
			padding: 20px;
		}

		.group-header {
			display: flex;
			align-items: center;
			gap: 10px;
			margin-bottom: 16px;
			border-bottom: 1px solid var(--border);
			padding-bottom: 12px;
		}

		.group-flag {
			font-size: 1.4rem;
		}

		.group-header h3 {
			font-size: 1.15rem;
			font-weight: 700;
			flex-grow: 1;
		}

		.group-count {
			font-size: 0.8rem;
			color: var(--text-muted);
			background: var(--bg-card);
			padding: 4px 10px;
			border-radius: 12px;
		}

		.checkbox-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
			gap: 10px;
		}

		.catalog-item {
			display: flex;
			align-items: flex-start;
			gap: 12px;
			background: var(--bg-card);
			border: 1px solid transparent;
			padding: 12px 14px;
			border-radius: 10px;
			cursor: pointer;
			transition: all 0.2s ease;
			user-select: none;
		}

		.catalog-item:hover {
			background: var(--bg-card-hover);
			border-color: rgba(255, 255, 255, 0.15);
		}

		.catalog-item input[type="checkbox"] {
			appearance: none;
			-webkit-appearance: none;
			width: 20px;
			height: 20px;
			background: rgba(255, 255, 255, 0.08);
			border: 2px solid rgba(255, 255, 255, 0.25);
			border-radius: 6px;
			cursor: pointer;
			position: relative;
			flex-shrink: 0;
			margin-top: 2px;
			transition: all 0.2s ease;
		}

		.catalog-item input[type="checkbox"]:checked {
			background: var(--accent);
			border-color: var(--accent-hover);
		}

		.catalog-item input[type="checkbox"]:checked::after {
			content: "✓";
			position: absolute;
			color: white;
			font-size: 13px;
			font-weight: bold;
			top: 50%%;
			left: 50%%;
			transform: translate(-50%%, -50%%);
		}

		.catalog-info {
			display: flex;
			flex-direction: column;
			gap: 3px;
		}

		.cat-award {
			font-size: 0.92rem;
			font-weight: 600;
			color: var(--text-main);
		}

		.cat-desc {
			font-size: 0.76rem;
			color: var(--text-muted);
			line-height: 1.3;
			display: -webkit-box;
			-webkit-line-clamp: 2;
			-webkit-box-orient: vertical;
			overflow: hidden;
		}

		.action-footer {
			background: var(--bg-surface);
			border: 1px solid var(--border);
			border-radius: var(--radius);
			padding: 24px;
			text-align: center;
			display: flex;
			flex-direction: column;
			gap: 16px;
			box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
		}

		.actions-row {
			display: flex;
			gap: 12px;
			justify-content: center;
			flex-wrap: wrap;
		}

		.btn-primary {
			background: linear-gradient(135deg, #8e44ad, #9b59b6);
			color: white;
			text-decoration: none;
			font-weight: 700;
			font-size: 1.05rem;
			padding: 14px 32px;
			border-radius: 10px;
			border: none;
			cursor: pointer;
			box-shadow: 0 6px 20px rgba(142, 68, 173, 0.4);
			transition: all 0.2s ease;
			display: inline-flex;
			align-items: center;
			gap: 8px;
		}

		.btn-primary:hover {
			transform: translateY(-2px);
			box-shadow: 0 10px 25px rgba(142, 68, 173, 0.6);
		}

		.btn-secondary {
			background: var(--bg-card);
			border: 1px solid var(--border);
			color: var(--text-main);
			font-weight: 600;
			font-size: 0.95rem;
			padding: 14px 24px;
			border-radius: 10px;
			cursor: pointer;
			transition: all 0.2s ease;
		}

		.btn-secondary:hover {
			background: var(--bg-card-hover);
			border-color: rgba(255, 255, 255, 0.2);
		}

		.url-display {
			background: var(--bg-base);
			border: 1px solid var(--border);
			border-radius: 8px;
			padding: 10px 14px;
			font-family: monospace;
			font-size: 0.82rem;
			color: var(--text-muted);
			word-break: break-all;
		}

		.toast {
			position: fixed;
			bottom: 30px;
			background: #27ae60;
			color: white;
			padding: 12px 24px;
			border-radius: 30px;
			font-weight: 600;
			box-shadow: 0 6px 20px rgba(0,0,0,0.4);
			opacity: 0;
			transform: translateY(20px);
			transition: all 0.3s ease;
			pointer-events: none;
			z-index: 999;
		}

		.toast.show {
			opacity: 1;
			transform: translateY(0);
		}
	</style>
</head>
<body>
	<div class="container">
		<header>
			<img src="/logo.png" alt="Film Festivals Logo" class="logo">
			<h1>Film Festivals</h1>
			<p class="subtitle">Customize your festival discovery catalogs in Stremio. Select the festival sections you want to appear in your Discover movies menu.</p>
		</header>

		<div class="toolbar">
			<div class="quick-filters">
				<button class="btn-filter" onclick="selectAll(true)">Select All</button>
				<button class="btn-filter" onclick="selectAll(false)">Deselect All</button>
				<button class="btn-filter" onclick="selectPreset('big-three')">Big Three (Cannes, Venice, Berlin)</button>
				<button class="btn-filter" onclick="selectPreset('auteurs')">Auteur & Indie Gems</button>
			</div>
			<div id="selection-badge" class="selection-badge">38 of 38 selected</div>
		</div>

		<div class="groups">
			%s
		</div>

		<div class="action-footer">
			<div class="actions-row">
				<a id="btn-install" href="#" class="btn-primary">
					<span>Install in Stremio</span>
				</a>
				<button class="btn-secondary" onclick="copyManifestURL()">
					Copy Manifest URL
				</button>
			</div>
			<div id="manifest-url-display" class="url-display"></div>
		</div>
	</div>

	<div id="toast" class="toast">Manifest URL copied to clipboard!</div>

	<script>
		const totalCatalogs = 38;
		const host = "%s";

		const presetBigThree = [
			'cannes-palme-dor', 'cannes-grand-prix', 'cannes-jury-prize', 'cannes-best-director', 'cannes-best-screenplay', 'cannes-best-actress', 'cannes-best-actor',
			'venice-golden-lion', 'venice-grand-jury-prize', 'venice-silver-lion-director', 'venice-best-screenplay', 'venice-coppa-volpi-actress', 'venice-coppa-volpi-actor',
			'berlin-golden-bear', 'berlin-silver-bear-grand-jury', 'berlin-silver-bear-director', 'berlin-silver-bear-screenplay', 'berlin-silver-bear-actress', 'berlin-silver-bear-actor'
		];

		const presetAuteurs = [
			'cannes-palme-dor', 'cannes-grand-prix', 'cannes-jury-prize', 'cannes-best-director',
			'venice-golden-lion', 'venice-grand-jury-prize', 'venice-silver-lion-director',
			'berlin-golden-bear', 'berlin-silver-bear-grand-jury', 'berlin-silver-bear-director',
			'locarno-golden-leopard', 'locarno-special-jury-prize', 'locarno-best-direction',
			'sundance-grand-jury-dramatic', 'sundance-audience-dramatic', 'sundance-directing-dramatic',
			'rotterdam-tiger-award', 'san-sebastian-golden-shell', 'san-sebastian-best-director',
			'fipresci-grand-prix', 'tiff-peoples-choice'
		];

		function getSelectedCatalogs() {
			const checkboxes = document.querySelectorAll('input[name="festival"]:checked');
			return Array.from(checkboxes).map(cb => cb.value);
		}

		function updateManifestURL() {
			const selected = getSelectedCatalogs();
			const count = selected.length;
			
			document.getElementById('selection-badge').innerText = count + " of " + totalCatalogs + " selected";

			let httpUrl = host + "/manifest.json";
			let stremioUrl = host.replace(/^https?:\/\//, "stremio://") + "/manifest.json";

			if (count < totalCatalogs && count > 0) {
				const query = "?festivals=" + selected.join(",");
				httpUrl += query;
				stremioUrl += query;
			} else if (count === 0) {
				httpUrl += "?festivals=none";
				stremioUrl += "?festivals=none";
			}

			document.getElementById('manifest-url-display').innerText = httpUrl;
			document.getElementById('btn-install').href = stremioUrl;
		}

		function selectAll(checked) {
			document.querySelectorAll('input[name="festival"]').forEach(cb => cb.checked = checked);
			updateManifestURL();
		}

		function selectPreset(name) {
			const targetList = name === 'big-three' ? presetBigThree : presetAuteurs;
			document.querySelectorAll('input[name="festival"]').forEach(cb => {
				cb.checked = targetList.includes(cb.value);
			});
			updateManifestURL();
		}

		function copyManifestURL() {
			const url = document.getElementById('manifest-url-display').innerText;
			navigator.clipboard.writeText(url).then(() => {
				const toast = document.getElementById('toast');
				toast.classList.add('show');
				setTimeout(() => toast.classList.remove('show'), 2500);
			});
		}

		// Initialize on load
		updateManifestURL();
	</script>
</body>
</html>`, builder.String(), currentHost)

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	c.Set(fiber.HeaderCacheControl, "no-cache")
	c.SendString(htmlContent)
}
