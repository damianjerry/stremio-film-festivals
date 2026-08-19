# Changelog

## [1.1.0](https://github.com/damianjerry/stremio-film-festivals/compare/v1.0.0...v1.1.0) (2026-08-19)


### Features

* add embedded branding logo, static configuration UI, and dynamic manifest filtering ([6a9dd69](https://github.com/damianjerry/stremio-film-festivals/commit/6a9dd6923ca54faf84f301b41ce682129fc5ba54))
* **branding:** update addon name to Film Festivals 2 | ElfHosted, shorten catalog names, and add ElfHosted/GitHub links ([f5ebaf2](https://github.com/damianjerry/stremio-film-festivals/commit/f5ebaf21b31a389b8eeecaea73af3e2fcbaeca99))
* **core:** implement dynamic catalog registry, fallback preview loader, and daily discovery randomizer ([b500c6e](https://github.com/damianjerry/stremio-film-festivals/commit/b500c6e405d007e5dd0f9924a4569c185bd3878a))
* **data:** add comprehensive festival winner datasets and extractor ([480bc36](https://github.com/damianjerry/stremio-film-festivals/commit/480bc3631b8873a4da46b3ca07db5529f4b10f2f))
* **ui:** add standard ElfHosted sponsor card and addons guide link ([d3ee2a3](https://github.com/damianjerry/stremio-film-festivals/commit/d3ee2a3fb15f299247ea3c0e8cec2f11448ffcb2))


### Bug Fixes

* **branding:** update logo filename reference to 'film festivals 2 logo.png' ([a6ed703](https://github.com/damianjerry/stremio-film-festivals/commit/a6ed7037faf54a8297111beebc9248e451ad580d))
* **data:** sanitize IMDb IDs, enforce Unix line endings, and correct canonical IDs for Anora, Triangle of Sadness, and Linha de Passe ([f508e81](https://github.com/damianjerry/stremio-film-festivals/commit/f508e81e0abffba285a710844fc941d38c11c653))
* **scraper:** implement Global Rules A, B, C for person rejection, movie type enforcement, and canonical IDs ([af26b46](https://github.com/damianjerry/stremio-film-festivals/commit/af26b4678af8041505ab69048d22d4f26752173a))
* **scraper:** resolve continuation-row director swaps, filter non-film entities, and verify modern Palme d'Or winners ([1a214ef](https://github.com/damianjerry/stremio-film-festivals/commit/1a214efb3ace2809f1b05cc3cca95020e0c662a6))
* **ui:** add installation step guide for local HTTP addons in Stremio Desktop ([d328f9e](https://github.com/damianjerry/stremio-film-festivals/commit/d328f9e2148affa1ddd18644ee9bf23ed2b02748))
* **ui:** dynamic client-side manifest URL generation for ElfHosted reverse proxy and local dev ([7273975](https://github.com/damianjerry/stremio-film-festivals/commit/727397581fdef9963dae0a5eb8fe56015e3640ff))
* **ui:** rename to Film Festivals 2, center logo inside box, and fix client-side port in manifest URLs ([9d0b1e5](https://github.com/damianjerry/stremio-film-festivals/commit/9d0b1e550fe86f83f2115c03ef8ed7d4d1d705a0))
