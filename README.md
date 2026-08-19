<p align="center">
  <img src="film%20festivals%202%20logo.png" alt="Film Festivals 2 Logo" width="160" height="160" style="border-radius: 28px;">
</p>

# <p align="center">Film Festivals 2 | ElfHosted</p>

<p align="center">
  <em>A comprehensive Stremio addon for discovering arthouse, auteur, and award-winning international cinema across the world's most prestigious film festivals.</em>
</p>

<p align="center">
  <a href="https://elfhosted.com"><img src="https://img.shields.io/badge/Hosted%20by-ElfHosted-blueviolet?style=flat-square" alt="ElfHosted"></a>
  <a href="https://github.com/damianjerry/stremio-film-festivals"><img src="https://img.shields.io/badge/GitHub-Repository-blue?style=flat-square&logo=github" alt="GitHub"></a>
  <a href="https://stremio-addons-guide.elfhosted.com"><img src="https://img.shields.io/badge/Stremio-Addons%20Guide-orange?style=flat-square" alt="Stremio Addons Guide"></a>
</p>

> [!TIP]
> ⚡ **Hosted by [ElfHosted](https://elfhosted.com)** — High-performance open-source apps and Stremio addons in the cloud.  
> 🔗 **GitHub Repository**: [damianjerry/stremio-film-festivals](https://github.com/damianjerry/stremio-film-festivals)

---

## Features

* **38 Curated Festival Catalogs**: Spanning the major international festival circuits, auteur showcases, and documentary awards.
* **Deterministic Daily Discovery**: Each day, catalogs feature a deterministically randomized discovery ordering so users discover different hidden gems every day, while remaining stable throughout the day for seamless caching.
* **Interactive Configuration UI (`/configure`)**: Select specific festival sections, choose presets (Big Three, Auteurs), and generate custom manifest links.
* **Complete Historical Datasets**: Every winning feature film from inaugural editions to the present is preserved and browsable.
* **IMDb & Cinemeta Native Integration**: Standard IMDb ID mapping (`tt...`) ensures posters, metadata, synopsis, trailers, and stream resolution work natively across Stremio Web, Desktop, Android, and TV apps.
* **Zero External Dependencies at Runtime**: Fast in-memory catalog store with bundled data, instant cold starts, and resilience.
* **ElfHosted & Docker Ready**: Easy to host as a standalone containerized service with full environment variable configuration.

---

## Festival & Award Catalogs (38 Catalogs)

### 🇫🇷 Cannes Film Festival
* **Palme d'Or** (`cannes-palme-dor` / alias `palme-dor-winners`) — Highest prize of the Official Competition (1939–present).
* **Grand Prix** (`cannes-grand-prix`) — 2nd most prestigious award at Cannes.
* **Jury Prize** (`cannes-jury-prize`) — Prix du Jury celebrating distinctive auteur cinema.
* **Best Director** (`cannes-best-director`) — Prix de la mise en scène.
* **Best Screenplay** (`cannes-best-screenplay`) — Prix du scénario.
* **Best Actress** (`cannes-best-actress`) — Prix d'interprétation féminine.
* **Best Actor** (`cannes-best-actor`) — Prix d'interprétation masculine.

### 🇮🇹 Venice International Film Festival
* **Golden Lion** (`venice-golden-lion` / alias `golden-lion-winners`) — Leone d'Oro for Best Film (1949–present).
* **Grand Jury Prize** (`venice-grand-jury-prize`) — Leone d'Argento - Gran Premio della Giuria.
* **Silver Lion for Best Director** (`venice-silver-lion-director`) — Premio per la migliore regia.
* **Best Screenplay** (`venice-best-screenplay`) — Golden Osella / Premio per la migliore sceneggiatura.
* **Coppa Volpi (Best Actress)** (`venice-coppa-volpi-actress`) — Volpi Cup for Best Actress.
* **Coppa Volpi (Best Actor)** (`venice-coppa-volpi-actor`) — Volpi Cup for Best Actor.

### 🇩🇪 Berlin International Film Festival (Berlinale)
* **Golden Bear** (`berlin-golden-bear` / alias `golden-bear-winners`) — Goldener Bär for Best Film (1951–present).
* **Silver Bear Grand Jury Prize** (`berlin-silver-bear-grand-jury`) — Großer Preis der Jury.
* **Silver Bear for Best Director** (`berlin-silver-bear-director`) — Silberner Bär für die beste Regie.
* **Silver Bear for Best Screenplay** (`berlin-silver-bear-screenplay`) — Silberner Bär für das beste Drehbuch.
* **Silver Bear for Best Actress / Leading Performance** (`berlin-silver-bear-actress`) — Silberner Bär für die beste Darstellerin / Hauptrolle.
* **Silver Bear for Best Actor** (`berlin-silver-bear-actor`) — Silberner Bär für den besten Darsteller.

### 🇨🇭 Locarno Film Festival
* **Golden Leopard** (`locarno-golden-leopard`) — Pardo d'oro for Best Film (1946–present).
* **Special Jury Prize** (`locarno-special-jury-prize`) — Premio speciale della giuria.
* **Best Direction** (`locarno-best-direction`) — Pardo per la miglior regia.

### 🇺🇸 Sundance Film Festival
* **Grand Jury Prize (Dramatic)** (`sundance-grand-jury-dramatic`) — Premier U.S. Dramatic Competition award.
* **Grand Jury Prize (Documentary)** (`sundance-grand-jury-doc`) — Premier U.S. Documentary Competition award.
* **Audience Award (Dramatic)** (`sundance-audience-dramatic`) — Audience Favorite in Dramatic Competition.
* **Audience Award (Documentary)** (`sundance-audience-doc`) — Audience Favorite in Documentary Competition.
* **Directing Award (Dramatic)** (`sundance-directing-dramatic`) — Best Directing in Dramatic Feature.
* **Directing Award (Documentary)** (`sundance-directing-doc`) — Best Directing in Documentary Feature.

### 🇨🇦 Toronto International Film Festival (TIFF)
* **People's Choice Award** (`tiff-peoples-choice`) — Top award voted by TIFF festival audiences (1978–present).

### 🇳🇱 International Film Festival Rotterdam (IFFR)
* **Tiger Award** (`rotterdam-tiger-award`) — Premier award for innovative, cutting-edge arthouse and emerging auteurs.

### 🇪🇸 San Sebastián International Film Festival
* **Golden Shell** (`san-sebastian-golden-shell`) — Concha de Oro for Best Film (1953–present).
* **Silver Shell for Best Director** (`san-sebastian-best-director`) — Concha de Plata a la mejor dirección.

### 🇨🇿 Karlovy Vary International Film Festival (KVIFF)
* **Crystal Globe** (`karlovy-vary-crystal-globe`) — Křišťálový glóbus for Best Film (1946–present).

### 🇬🇧 BFI London Film Festival
* **Best Film Award** (`bfi-london-best-film`) — Official Competition Best Film & Sutherland Trophy winners.

### 🇳🇱 IDFA (International Documentary Film Festival Amsterdam)
* **Best Feature Documentary** (`idfa-best-film`) — Top prize at the world's leading documentary festival.

### 🇺🇸 SXSW Film & TV Festival
* **Grand Jury Award (Narrative)** (`sxsw-grand-jury-narrative`) — Grand Jury Award for Narrative Feature.

### 🌍 FIPRESCI (International Federation of Film Critics)
* **Grand Prix (Film of the Year)** (`fipresci-grand-prix`) — Voted by international critics across worldwide cinema.

### 🇺🇸 Academy Awards (Oscars)
* **Best Picture** (`academy-awards-best-picture` / alias `academy-awards-winners`) — Best Picture winners (1927–present).

---

## Discovery & Randomization Engine

Standard Stremio catalogs frequently show only the newest or oldest entries, leaving dozens of classic or mid-era gems buried.

This addon features **Deterministic Daily Randomization**:

1. **Daily Discovery Rotation**: Every catalog has its complete historical dataset deterministically permuted based on a daily seed (`hash(catalog_id + ":" + YYYY-MM-DD)`).
2. **Stable Within the Day**: Repeated requests from Stremio on the same day receive the exact same ordering, ensuring full compatibility with HTTP caching (`Cache-Control: max-age=24h`) and `ETag` 304 Not Modified validations.
3. **Daily Refresh**: At midnight UTC, a new seed produces a new stable discovery permutation.
4. **Complete Dataset Preserved**: All films remain accessible; users can scroll through the entire catalog history.
5. **Configurable Ordering Modes**:
   * `daily-random` (default): Deterministic daily rotating shuffle.
   * `chronological-desc`: Sorted newest to oldest.
   * `chronological-asc`: Sorted oldest to newest.

---

## Configuration & Installation

### Configuration Page
Visit `/configure` on your deployed addon to customize which catalogs appear in your Stremio Discover menu:
```text
https://<your-subdomain>.elfhosted.com/configure
```

### In Stremio
Enter the manifest URL in the Addon search box in Stremio:
```text
https://<your-subdomain>.elfhosted.com/manifest.json
```
Or for filtered catalogs:
```text
https://<your-subdomain>.elfhosted.com/manifest.json?festivals=cannes-palme-dor,venice-golden-lion,berlin-golden-bear
```

---

## Running Locally

### Option 1: Native Go Executable
```bash
# Build
go build -o stremio-film-festivals .

# Run
./stremio-film-festivals -port 8080 -dataDir data -order daily-random
```

### Option 2: Docker
```bash
# Build Docker image
docker build -t stremio-film-festivals .

# Run container
docker run -d --name stremio-film-festivals -p 8080:8080 stremio-film-festivals
```

---

## Configuration & Environment Variables

All parameters can be configured via command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
| :--- | :--- | :--- | :--- |
| `-bindAddr` | `BIND_ADDR` | `0.0.0.0` | Network interface address to bind to |
| `-port` | `PORT` | `8080` | HTTP port to listen on |
| `-dataDir` | `DATA_DIR` | `data` | Directory containing catalog CSV files |
| `-logLevel` | `LOG_LEVEL` | `info` | Log verbosity (`debug`, `info`, `warn`, `error`) |
| `-cacheAge` | `CACHE_AGE` | `24h` | HTTP Cache-Control max-age header duration |
| `-order` | `ORDER_MODE` | `daily-random` | Ordering mode (`daily-random`, `chronological-desc`, `chronological-asc`) |

---

## ElfHosted Deployment

* **Hosted on ElfHosted**: [https://elfhosted.com](https://elfhosted.com)
* The container is built using a secure multi-stage build (`gcr.io/distroless/static:nonroot`) and runs as an unprivileged user (`nonroot:nonroot`).
* Festival datasets and branding logo are pre-packaged directly inside the container (`/data`), meaning no external volume mounting is strictly required for basic operation.
* Standard healthcheck endpoint is available at `/health`.
* Reverse proxying (Traefik/Nginx) with standard HTTPS headers is fully supported.

---

## License

MIT License.
