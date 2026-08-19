package main

import (
	"github.com/deflix-tv/go-stremio"
)

const (
	version     = "1.0.0"
	redirectURL = "https://www.deflix.tv"
)

// CatalogConfig holds metadata and configuration for a festival catalog.
type CatalogConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Festival    string   `json:"festival"`
	Award       string   `json:"award"`
	Description string   `json:"description"`
	CSVFile     string   `json:"csvFile"`
	Aliases     []string `json:"aliases,omitempty"`
}

// FestivalCatalogs defines all supported film festival and award catalogs.
var FestivalCatalogs = []CatalogConfig{
	// Cannes Film Festival
	{
		ID:          "cannes-palme-dor",
		Name:        "Cannes — Palme d'Or",
		Festival:    "Cannes Film Festival",
		Award:       "Palme d'Or",
		Description: "Palme d'Or (and Grand Prix du Festival) winners at the Cannes Film Festival",
		CSVFile:     "cannes-palme-dor.csv",
		Aliases:     []string{"palme-dor-winners"},
	},
	{
		ID:          "cannes-grand-prix",
		Name:        "Cannes — Grand Prix",
		Festival:    "Cannes Film Festival",
		Award:       "Grand Prix",
		Description: "Grand Prix (Grand Prix Spécial du Jury) winners at the Cannes Film Festival",
		CSVFile:     "cannes-grand-prix.csv",
	},
	{
		ID:          "cannes-jury-prize",
		Name:        "Cannes — Jury Prize",
		Festival:    "Cannes Film Festival",
		Award:       "Jury Prize",
		Description: "Prix du Jury winners at the Cannes Film Festival",
		CSVFile:     "cannes-jury-prize.csv",
	},
	{
		ID:          "cannes-best-director",
		Name:        "Cannes — Best Director",
		Festival:    "Cannes Film Festival",
		Award:       "Best Director",
		Description: "Prix de la mise en scène (Best Director) winners at the Cannes Film Festival",
		CSVFile:     "cannes-best-director.csv",
	},
	{
		ID:          "cannes-best-screenplay",
		Name:        "Cannes — Best Screenplay",
		Festival:    "Cannes Film Festival",
		Award:       "Best Screenplay",
		Description: "Prix du scénario (Best Screenplay) winners at the Cannes Film Festival",
		CSVFile:     "cannes-best-screenplay.csv",
	},
	{
		ID:          "cannes-best-actress",
		Name:        "Cannes — Best Actress",
		Festival:    "Cannes Film Festival",
		Award:       "Best Actress",
		Description: "Prix d'interprétation féminine (Best Actress) winners at the Cannes Film Festival",
		CSVFile:     "cannes-best-actress.csv",
	},
	{
		ID:          "cannes-best-actor",
		Name:        "Cannes — Best Actor",
		Festival:    "Cannes Film Festival",
		Award:       "Best Actor",
		Description: "Prix d'interprétation masculine (Best Actor) winners at the Cannes Film Festival",
		CSVFile:     "cannes-best-actor.csv",
	},

	// Venice International Film Festival
	{
		ID:          "venice-golden-lion",
		Name:        "Venice — Golden Lion",
		Festival:    "Venice Film Festival",
		Award:       "Golden Lion",
		Description: "Golden Lion (Leone d'Oro) winners at the Venice International Film Festival",
		CSVFile:     "venice-golden-lion.csv",
		Aliases:     []string{"golden-lion-winners"},
	},
	{
		ID:          "venice-grand-jury-prize",
		Name:        "Venice — Grand Jury Prize",
		Festival:    "Venice Film Festival",
		Award:       "Grand Jury Prize",
		Description: "Grand Jury Prize (Leone d'Argento - Gran Premio della Giuria) winners at the Venice Film Festival",
		CSVFile:     "venice-grand-jury-prize.csv",
	},
	{
		ID:          "venice-silver-lion-director",
		Name:        "Venice — Silver Lion (Best Director)",
		Festival:    "Venice Film Festival",
		Award:       "Silver Lion for Best Director",
		Description: "Silver Lion for Best Director (Premio per la migliore regia) winners at the Venice Film Festival",
		CSVFile:     "venice-silver-lion-director.csv",
	},
	{
		ID:          "venice-best-screenplay",
		Name:        "Venice — Best Screenplay",
		Festival:    "Venice Film Festival",
		Award:       "Best Screenplay",
		Description: "Golden Osella / Best Screenplay award winners at the Venice Film Festival",
		CSVFile:     "venice-best-screenplay.csv",
	},
	{
		ID:          "venice-coppa-volpi-actress",
		Name:        "Venice — Coppa Volpi (Best Actress)",
		Festival:    "Venice Film Festival",
		Award:       "Coppa Volpi for Best Actress",
		Description: "Coppa Volpi for Best Actress winners at the Venice International Film Festival",
		CSVFile:     "venice-coppa-volpi-actress.csv",
	},
	{
		ID:          "venice-coppa-volpi-actor",
		Name:        "Venice — Coppa Volpi (Best Actor)",
		Festival:    "Venice Film Festival",
		Award:       "Coppa Volpi for Best Actor",
		Description: "Coppa Volpi for Best Actor winners at the Venice International Film Festival",
		CSVFile:     "venice-coppa-volpi-actor.csv",
	},

	// Berlin International Film Festival (Berlinale)
	{
		ID:          "berlin-golden-bear",
		Name:        "Berlin — Golden Bear",
		Festival:    "Berlin International Film Festival",
		Award:       "Golden Bear",
		Description: "Golden Bear (Goldener Bär) winners at the Berlin International Film Festival",
		CSVFile:     "berlin-golden-bear.csv",
		Aliases:     []string{"golden-bear-winners"},
	},
	{
		ID:          "berlin-silver-bear-grand-jury",
		Name:        "Berlin — Silver Bear Grand Jury Prize",
		Festival:    "Berlin International Film Festival",
		Award:       "Silver Bear Grand Jury Prize",
		Description: "Silver Bear Grand Jury Prize (Großer Preis der Jury) winners at the Berlinale",
		CSVFile:     "berlin-silver-bear-grand-jury.csv",
	},
	{
		ID:          "berlin-silver-bear-director",
		Name:        "Berlin — Silver Bear for Best Director",
		Festival:    "Berlin International Film Festival",
		Award:       "Silver Bear for Best Director",
		Description: "Silver Bear for Best Director (Beste Regie) winners at the Berlinale",
		CSVFile:     "berlin-silver-bear-director.csv",
	},
	{
		ID:          "berlin-silver-bear-screenplay",
		Name:        "Berlin — Silver Bear for Best Screenplay",
		Festival:    "Berlin International Film Festival",
		Award:       "Silver Bear for Best Screenplay",
		Description: "Silver Bear for Best Screenplay (Bestes Drehbuch) winners at the Berlinale",
		CSVFile:     "berlin-silver-bear-screenplay.csv",
	},
	{
		ID:          "berlin-silver-bear-actress",
		Name:        "Berlin — Silver Bear for Best Actress / Leading Performance",
		Festival:    "Berlin International Film Festival",
		Award:       "Silver Bear for Best Actress / Leading Performance",
		Description: "Silver Bear for Best Actress / Best Leading Performance winners at the Berlinale",
		CSVFile:     "berlin-silver-bear-actress.csv",
	},
	{
		ID:          "berlin-silver-bear-actor",
		Name:        "Berlin — Silver Bear for Best Actor",
		Festival:    "Berlin International Film Festival",
		Award:       "Silver Bear for Best Actor",
		Description: "Silver Bear for Best Actor winners at the Berlinale",
		CSVFile:     "berlin-silver-bear-actor.csv",
	},

	// Locarno Film Festival
	{
		ID:          "locarno-golden-leopard",
		Name:        "Locarno — Golden Leopard",
		Festival:    "Locarno Film Festival",
		Award:       "Golden Leopard",
		Description: "Golden Leopard (Pardo d'oro) winners at the Locarno Film Festival",
		CSVFile:     "locarno-golden-leopard.csv",
	},
	{
		ID:          "locarno-special-jury-prize",
		Name:        "Locarno — Special Jury Prize",
		Festival:    "Locarno Film Festival",
		Award:       "Special Jury Prize",
		Description: "Special Jury Prize (Premio speciale della giuria) winners at the Locarno Film Festival",
		CSVFile:     "locarno-special-jury-prize.csv",
	},
	{
		ID:          "locarno-best-direction",
		Name:        "Locarno — Best Direction",
		Festival:    "Locarno Film Festival",
		Award:       "Best Direction",
		Description: "Pardo for Best Direction winners at the Locarno Film Festival",
		CSVFile:     "locarno-best-direction.csv",
	},

	// Sundance Film Festival
	{
		ID:          "sundance-grand-jury-dramatic",
		Name:        "Sundance — Grand Jury Prize (Dramatic)",
		Festival:    "Sundance Film Festival",
		Award:       "Grand Jury Prize — Dramatic",
		Description: "Grand Jury Prize (U.S. Dramatic) winners at the Sundance Film Festival",
		CSVFile:     "sundance-grand-jury-dramatic.csv",
	},
	{
		ID:          "sundance-grand-jury-doc",
		Name:        "Sundance — Grand Jury Prize (Documentary)",
		Festival:    "Sundance Film Festival",
		Award:       "Grand Jury Prize — Documentary",
		Description: "Grand Jury Prize (U.S. Documentary) winners at the Sundance Film Festival",
		CSVFile:     "sundance-grand-jury-doc.csv",
	},
	{
		ID:          "sundance-audience-dramatic",
		Name:        "Sundance — Audience Award (Dramatic)",
		Festival:    "Sundance Film Festival",
		Award:       "Audience Award — Dramatic",
		Description: "Audience Award (U.S. Dramatic) winners at the Sundance Film Festival",
		CSVFile:     "sundance-audience-dramatic.csv",
	},
	{
		ID:          "sundance-audience-doc",
		Name:        "Sundance — Audience Award (Documentary)",
		Festival:    "Sundance Film Festival",
		Award:       "Audience Award — Documentary",
		Description: "Audience Award (U.S. Documentary) winners at the Sundance Film Festival",
		CSVFile:     "sundance-audience-doc.csv",
	},
	{
		ID:          "sundance-directing-dramatic",
		Name:        "Sundance — Directing Award (Dramatic)",
		Festival:    "Sundance Film Festival",
		Award:       "Directing Award — Dramatic",
		Description: "Directing Award (U.S. Dramatic) winners at the Sundance Film Festival",
		CSVFile:     "sundance-directing-dramatic.csv",
	},
	{
		ID:          "sundance-directing-doc",
		Name:        "Sundance — Directing Award (Documentary)",
		Festival:    "Sundance Film Festival",
		Award:       "Directing Award — Documentary",
		Description: "Directing Award (U.S. Documentary) winners at the Sundance Film Festival",
		CSVFile:     "sundance-directing-doc.csv",
	},

	// Toronto International Film Festival (TIFF)
	{
		ID:          "tiff-peoples-choice",
		Name:        "Toronto (TIFF) — People's Choice Award",
		Festival:    "Toronto International Film Festival",
		Award:       "People's Choice Award",
		Description: "People's Choice Award winners at the Toronto International Film Festival",
		CSVFile:     "tiff-peoples-choice.csv",
	},

	// International Film Festival Rotterdam (IFFR)
	{
		ID:          "rotterdam-tiger-award",
		Name:        "Rotterdam (IFFR) — Tiger Award",
		Festival:    "International Film Festival Rotterdam",
		Award:       "Tiger Award",
		Description: "Tiger Award winners for innovative and visionary cinema at IFFR",
		CSVFile:     "rotterdam-tiger-award.csv",
	},

	// San Sebastián International Film Festival
	{
		ID:          "san-sebastian-golden-shell",
		Name:        "San Sebastián — Golden Shell",
		Festival:    "San Sebastián International Film Festival",
		Award:       "Golden Shell",
		Description: "Golden Shell (Concha de Oro) winners at the San Sebastián International Film Festival",
		CSVFile:     "san-sebastian-golden-shell.csv",
	},
	{
		ID:          "san-sebastian-best-director",
		Name:        "San Sebastián — Silver Shell for Best Director",
		Festival:    "San Sebastián International Film Festival",
		Award:       "Silver Shell for Best Director",
		Description: "Silver Shell for Best Director (Concha de Plata a la mejor dirección) winners at San Sebastián",
		CSVFile:     "san-sebastian-best-director.csv",
	},

	// Karlovy Vary International Film Festival
	{
		ID:          "karlovy-vary-crystal-globe",
		Name:        "Karlovy Vary — Crystal Globe",
		Festival:    "Karlovy Vary International Film Festival",
		Award:       "Crystal Globe",
		Description: "Crystal Globe (Křišťálový glóbus) winners at the Karlovy Vary International Film Festival",
		CSVFile:     "karlovy-vary-crystal-globe.csv",
	},

	// BFI London Film Festival
	{
		ID:          "bfi-london-best-film",
		Name:        "BFI London — Best Film",
		Festival:    "BFI London Film Festival",
		Award:       "Best Film / Sutherland Trophy",
		Description: "Best Film and Sutherland Trophy winners at the BFI London Film Festival",
		CSVFile:     "bfi-london-best-film.csv",
	},

	// International Documentary Film Festival Amsterdam (IDFA)
	{
		ID:          "idfa-best-film",
		Name:        "IDFA — Best Feature Documentary",
		Festival:    "International Documentary Film Festival Amsterdam",
		Award:       "IDFA Award for Best Feature-Length Documentary",
		Description: "Best Feature-Length Documentary winners at IDFA",
		CSVFile:     "idfa-best-film.csv",
	},

	// SXSW Film Festival
	{
		ID:          "sxsw-grand-jury-narrative",
		Name:        "SXSW — Grand Jury Award (Narrative)",
		Festival:    "SXSW Film & TV Festival",
		Award:       "Grand Jury Award — Narrative Feature",
		Description: "Grand Jury Award for Narrative Feature winners at South by Southwest (SXSW)",
		CSVFile:     "sxsw-grand-jury-narrative.csv",
	},

	// FIPRESCI International Film Critics
	{
		ID:          "fipresci-grand-prix",
		Name:        "FIPRESCI — Grand Prix (Film of the Year)",
		Festival:    "FIPRESCI",
		Award:       "FIPRESCI Grand Prix",
		Description: "FIPRESCI Grand Prix for Best Film of the Year chosen by international film critics",
		CSVFile:     "fipresci-grand-prix.csv",
	},

	// Academy Awards
	{
		ID:          "academy-awards-best-picture",
		Name:        "Academy Awards — Best Picture",
		Festival:    "Academy Awards",
		Award:       "Best Picture",
		Description: "Academy Award for Best Picture winners (1927–present)",
		CSVFile:     "academy-awards-best-picture.csv",
		Aliases:     []string{"academy-awards-winners"},
	},
}

// BuildManifest constructs the Stremio Manifest from the catalog configurations.
func BuildManifest() stremio.Manifest {
	catalogs := make([]stremio.CatalogItem, 0, len(FestivalCatalogs))
	for _, cat := range FestivalCatalogs {
		catalogs = append(catalogs, stremio.CatalogItem{
			Type: "movie",
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	return stremio.Manifest{
		ID:          "tv.deflix.stremio-film-festivals",
		Name:        "Film Festivals 2 | ElfHosted",
		Description: "Comprehensive catalogs of world-renowned arthouse, auteur, and international film festival award winners: Cannes, Venice, Berlinale, Locarno, Sundance, TIFF, Rotterdam, San Sebastián, Karlovy Vary, BFI London, IDFA, SXSW, FIPRESCI & Academy Awards. Hosted on ElfHosted.",
		Version:     version,

		ResourceItems: []stremio.ResourceItem{
			{
				Name: "catalog",
			},
		},
		Types:    []string{"movie"},
		Catalogs: catalogs,

		IDprefixes: []string{"tt"},
		Background: "https://images.metahub.space/background/medium/tt0078788/img",
		Logo:       "/logo.png",
	}
}

// FindCatalogConfig finds a CatalogConfig by its primary ID or alias.
func FindCatalogConfig(id string) (CatalogConfig, bool) {
	for _, cat := range FestivalCatalogs {
		if cat.ID == id {
			return cat, true
		}
		for _, alias := range cat.Aliases {
			if alias == id {
				return cat, true
			}
		}
	}
	return CatalogConfig{}, false
}
