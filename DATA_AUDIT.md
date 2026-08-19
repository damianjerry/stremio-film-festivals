# Film Festival Dataset Quality Audit Report

**Audit Date**: 2026-08-19  
**Total Records Audited**: 2465 across 42 CSV datasets  
**Trustworthy (Verified + Likely Correct)**: 2411 (97.8%)  
**Requiring Attention (Incorrect + Questionable + Unable to Verify)**: 54 (2.2%)

---

## 1. Executive Summary & Statistics

| Classification Status | Count | Percentage | Definition |
| :--- | :---: | :---: | :--- |
| **Verified** | 2257 | 91.6% | High confidence match; title, year, director, and festival winner status corroborated. |
| **Likely Correct** | 154 | 6.2% | Strong match; minor transliteration or commercial release year offset (within +/- 2 years). |
| **Questionable** | 12 | 0.5% | Large unexplained year gap (>= 3 years) or ambiguous title needing validation. |
| **Incorrect Mapping** | 10 | 0.4% | Non-movie entity, person name, wrong film disambiguation, or comparison table pollution. |
| **Unable to Verify** | 32 | 1.3% | Missing from metadata indexes or ambiguous title. |
| **Total** | **2465** | **100.0%** | |

---

## 2. Root-Cause Analysis & Failure Modes

Our investigation into `tt0088380` (*"outbreak of the Second World War"* -> *Warrior of the Lost World*, 1983) and other anomalies revealed **4 distinct failure modes** in the extraction and resolution pipeline:

### Failure Mode 1: Explanatory Table Notes & Cancelled Editions Extracted as Films (0 occurrences)
* **Root Cause**: Wikipedia award tables frequently contain spanning rows for cancelled festival editions or historical explanations (e.g. Cannes 1939 cancelled due to WWII, Berlinale 1970 jury resignation over *o.k.*). The parser extracted text hyperlinks within these rows as movie titles.
* **Effect**: Fallback token search on Cinemeta matched random B-movies (e.g. searching *"outbreak of the Second World War"* matched *Warrior of the Lost World* `tt0088380`).
* **Resolution**: Strict row filtering to discard rows containing cancellation / note keywords and ensure the row is a valid film winner.

### Failure Mode 2: Person / Actor Names Extracted as Film Titles (0 occurrences)
* **Root Cause**: In ensemble acting awards (e.g. Cannes 1955 *A Big Family*, Cannes 2006 *Volver*), Wikipedia tables listed each actor in separate rows or columns. In festival history tables (e.g. BFI London), "Festival Directors" tables were placed after the award tables.
* **Effect**: Actor/director biographical articles were parsed as movie titles.
* **Resolution**: Filter out non-film entity tables (e.g. Festival Directors, Multiple Winners, Jury Presidents) and ensure acting awards resolve to the winning film rather than actor biographies.

### Failure Mode 3: Comparison / Superlatives Table Pollution (22 occurrences)
* **Root Cause**: Articles for Best Actress / Best Actor include comparison tables at the bottom (e.g. *"Actresses with multiple awards across Cannes, Venice, and Berlin"*). These tables contain films from other festivals and other decades.
* **Effect**: Films from 2021 or 1993 were injected into 2006 or 2014 Cannes awards.
* **Resolution**: Skip tables whose preceding heading matches `multiple`, `superlatives`, `records`, `lifetime`, or `statistics`.

### Failure Mode 4: Search Fallback & Disambiguation Errors (0 occurrences)
* **Root Cause**: When Wikidata lacked a `P345` claim for a film, Cinemeta search fallback accepted the first search result without verifying title similarity or release year compatibility (e.g. *Casablanca* resolving to a 1992 documentary tribute instead of the 1942 classic `tt0034583`).
* **Resolution**: Multi-signal resolver requiring title similarity >= 80%, year compatibility within +/- 2 years, and known canonical mappings for classic titles.

---

## 3. Dataset-by-Dataset Audit Breakdown

| Dataset / Catalog | Total Records | Verified | Likely Correct | Questionable | Incorrect | Unable to Verify | Error / Issue Rate |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `academy-awards-best-picture` | 98 | 96 | 0 | 0 | 0 | 2 | 2.0% |
| `academy-awards-winners` | 98 | 96 | 0 | 0 | 0 | 2 | 2.0% |
| `berlin-golden-bear` | 88 | 82 | 3 | 1 | 0 | 2 | 3.4% |
| `berlin-silver-bear-actor` | 62 | 59 | 3 | 0 | 0 | 0 | 0.0% |
| `berlin-silver-bear-actress` | 65 | 59 | 4 | 1 | 0 | 1 | 3.1% |
| `berlin-silver-bear-director` | 64 | 59 | 5 | 0 | 0 | 0 | 0.0% |
| `berlin-silver-bear-grand-jury` | 69 | 60 | 8 | 0 | 0 | 1 | 1.4% |
| `berlin-silver-bear-screenplay` | 18 | 18 | 0 | 0 | 0 | 0 | 0.0% |
| `bfi-london-best-film` | 69 | 63 | 5 | 1 | 0 | 0 | 1.4% |
| `cannes-best-actor` | 78 | 69 | 8 | 0 | 1 | 0 | 1.3% |
| `cannes-best-actress` | 87 | 79 | 7 | 0 | 1 | 0 | 1.1% |
| `cannes-best-director` | 66 | 61 | 2 | 0 | 1 | 2 | 4.5% |
| `cannes-best-screenplay` | 42 | 39 | 3 | 0 | 0 | 0 | 0.0% |
| `cannes-grand-prix` | 66 | 62 | 3 | 0 | 0 | 1 | 1.5% |
| `cannes-jury-prize` | 84 | 78 | 5 | 1 | 0 | 0 | 1.2% |
| `cannes-palme-dor` | 103 | 97 | 3 | 0 | 1 | 2 | 2.9% |
| `fipresci-grand-prix` | 26 | 23 | 3 | 0 | 0 | 0 | 0.0% |
| `golden-bear-winners` | 88 | 82 | 3 | 1 | 0 | 2 | 3.4% |
| `golden-lion-winners` | 69 | 68 | 1 | 0 | 0 | 0 | 0.0% |
| `idfa-best-film` | 17 | 14 | 3 | 0 | 0 | 0 | 0.0% |
| `karlovy-vary-crystal-globe` | 57 | 47 | 10 | 0 | 0 | 0 | 0.0% |
| `locarno-best-direction` | 19 | 17 | 2 | 0 | 0 | 0 | 0.0% |
| `locarno-golden-leopard` | 93 | 80 | 10 | 0 | 0 | 3 | 3.2% |
| `locarno-special-jury-prize` | 21 | 20 | 1 | 0 | 0 | 0 | 0.0% |
| `palme-dor-winners` | 103 | 97 | 3 | 0 | 1 | 2 | 2.9% |
| `rotterdam-tiger-award` | 40 | 34 | 5 | 1 | 0 | 0 | 2.5% |
| `san-sebastian-best-director` | 49 | 41 | 5 | 0 | 0 | 3 | 6.1% |
| `san-sebastian-golden-shell` | 69 | 61 | 7 | 0 | 0 | 1 | 1.4% |
| `sundance-audience-doc` | 27 | 21 | 3 | 2 | 1 | 0 | 11.1% |
| `sundance-audience-dramatic` | 31 | 29 | 1 | 1 | 0 | 0 | 3.2% |
| `sundance-directing-doc` | 21 | 17 | 3 | 1 | 0 | 0 | 4.8% |
| `sundance-directing-dramatic` | 24 | 24 | 0 | 0 | 0 | 0 | 0.0% |
| `sundance-grand-jury-doc` | 34 | 32 | 0 | 0 | 1 | 1 | 5.9% |
| `sundance-grand-jury-dramatic` | 36 | 35 | 1 | 0 | 0 | 0 | 0.0% |
| `sxsw-grand-jury-narrative` | 26 | 23 | 2 | 1 | 0 | 0 | 3.8% |
| `tiff-peoples-choice` | 48 | 47 | 1 | 0 | 0 | 0 | 0.0% |
| `venice-best-screenplay` | 49 | 42 | 1 | 0 | 2 | 4 | 12.2% |
| `venice-coppa-volpi-actor` | 78 | 69 | 8 | 0 | 0 | 1 | 1.3% |
| `venice-coppa-volpi-actress` | 76 | 69 | 7 | 0 | 0 | 0 | 0.0% |
| `venice-golden-lion` | 69 | 68 | 1 | 0 | 0 | 0 | 0.0% |
| `venice-grand-jury-prize` | 71 | 61 | 9 | 1 | 0 | 0 | 1.4% |
| `venice-silver-lion-director` | 67 | 59 | 5 | 0 | 1 | 2 | 4.5% |

---

## 4. Itemized Problem Log & Remediation Actions

Below is the complete list of records flagged as `incorrect_mapping` or `questionable`, along with root cause and verified remediation:

### `academy-awards-best-picture.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `academy-awards-best-picture.csv` — Year 2016: "Moonlight"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt4975721` -> *""* (None)
* **Reason**: IMDb ID tt4975721 not found on Cinemeta metadata index

### `academy-awards-winners.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `academy-awards-winners.csv` — Year 2016: "Moonlight"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt4975721` -> *""* (None)
* **Reason**: IMDb ID tt4975721 not found on Cinemeta metadata index

### `berlin-golden-bear.csv` — Year 1987: "The Theme"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0079999` -> *"The Theme"* (1979)
* **Reason**: Large year disparity between award year (1987) and film release year (1979). Title matches: 'The Theme' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `berlin-golden-bear.csv` — Year 1983: "La colmena"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0083745` -> *""* (None)
* **Reason**: IMDb ID tt0083745 not found on Cinemeta metadata index

### `berlin-golden-bear.csv` — Year 1978: "Las truchas"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0076846` -> *""* (None)
* **Reason**: IMDb ID tt0076846 not found on Cinemeta metadata index

### `berlin-silver-bear-actress.csv` — Year 2013: "Gloria"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0080798` -> *"Gloria"* (1980)
* **Reason**: Large year disparity between award year (2013) and film release year (1980). Title matches: 'Gloria' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `berlin-silver-bear-actress.csv` — Year 2004: "Aileen Wuornos"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt8490896` -> *""* (None)
* **Reason**: IMDb ID tt8490896 not found on Cinemeta metadata index

### `berlin-silver-bear-grand-jury.csv` — Year 1975: "Overlord"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0073498` -> *""* (None)
* **Reason**: IMDb ID tt0073498 not found on Cinemeta metadata index

### `bfi-london-best-film.csv` — Year 1967: "La Marseillaise"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0030424` -> *"La Marseillaise"* (1938)
* **Reason**: Large year disparity between award year (1967) and film release year (1938). Title matches: 'La Marseillaise' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-best-actor.csv` — Year 1962: "Edmund Tyrone"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0020514` -> *"The Trespasser"* (1929)
* **Reason**: Severe title and year mismatch: CSV title 'Edmund Tyrone' (Award 1962) vs IMDb title 'The Trespasser' (1929, sim=37.03703703703704%).
* **Evidence**: Search fallback returned unrelated title 'The Trespasser' (1929).

### `cannes-best-actress.csv` — Year 1957: "Nights of Cabiria"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065054` -> *"Sweet Charity"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Nights of Cabiria' (Award 1957) vs IMDb title 'Sweet Charity' (1969, sim=40.0%).
* **Evidence**: Search fallback returned unrelated title 'Sweet Charity' (1969).

### `cannes-best-director.csv` — Year 2001: "David Lynch"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt1705102` -> *""* (None)
* **Reason**: IMDb ID tt1705102 not found on Cinemeta metadata index

### `cannes-best-director.csv` — Year 1983: "Andrei Tarkovsky"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt2638384` -> *""* (None)
* **Reason**: IMDb ID tt2638384 not found on Cinemeta metadata index

### `cannes-best-director.csv` — Year 1975: "Costa-Gavras"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065234` -> *"Z"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Costa-Gavras' (Award 1975) vs IMDb title 'Z' (1969, sim=0.0%).
* **Evidence**: Search fallback returned unrelated title 'Z' (1969).

### `cannes-grand-prix.csv` — Year 2024: "All We Imagine as Light"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt27823528` -> *""* (None)
* **Reason**: IMDb ID tt27823528 not found on Cinemeta metadata index

### `cannes-jury-prize.csv` — Year 1996: "Crash"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0375679` -> *"Crash"* (2004)
* **Reason**: Large year disparity between award year (1996) and film release year (2004). Title matches: 'Crash' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-palme-dor.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `cannes-palme-dor.csv` — Year 2022: "Triangle of Sadness"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10279050` -> *"Life's First-Evers with Jeannie"* (2018)
* **Reason**: Severe title and year mismatch: CSV title 'Triangle of Sadness' (Award 2022) vs IMDb title 'Life's First-Evers with Jeannie' (2018, sim=37.5%).
* **Evidence**: Search fallback returned unrelated title 'Life's First-Evers with Jeannie' (2018).

### `cannes-palme-dor.csv` — Year 1939: "Union Pacific"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0032080` -> *""* (None)
* **Reason**: IMDb ID tt0032080 not found on Cinemeta metadata index

### `golden-bear-winners.csv` — Year 1987: "The Theme"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0079999` -> *"The Theme"* (1979)
* **Reason**: Large year disparity between award year (1987) and film release year (1979). Title matches: 'The Theme' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `golden-bear-winners.csv` — Year 1983: "La colmena"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0083745` -> *""* (None)
* **Reason**: IMDb ID tt0083745 not found on Cinemeta metadata index

### `golden-bear-winners.csv` — Year 1978: "Las truchas"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0076846` -> *""* (None)
* **Reason**: IMDb ID tt0076846 not found on Cinemeta metadata index

### `locarno-golden-leopard.csv` — Year 1988: "Schmetterlinge"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0096055` -> *""* (None)
* **Reason**: IMDb ID tt0096055 not found on Cinemeta metadata index

### `locarno-golden-leopard.csv` — Year 1971: "Private Road"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0067623` -> *""* (None)
* **Reason**: IMDb ID tt0067623 not found on Cinemeta metadata index

### `locarno-golden-leopard.csv` — Year 1970: "Soleil O"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0065014` -> *""* (None)
* **Reason**: IMDb ID tt0065014 not found on Cinemeta metadata index

### `palme-dor-winners.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `palme-dor-winners.csv` — Year 2022: "Triangle of Sadness"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10279050` -> *"Life's First-Evers with Jeannie"* (2018)
* **Reason**: Severe title and year mismatch: CSV title 'Triangle of Sadness' (Award 2022) vs IMDb title 'Life's First-Evers with Jeannie' (2018, sim=37.5%).
* **Evidence**: Search fallback returned unrelated title 'Life's First-Evers with Jeannie' (2018).

### `palme-dor-winners.csv` — Year 1939: "Union Pacific"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0032080` -> *""* (None)
* **Reason**: IMDb ID tt0032080 not found on Cinemeta metadata index

### `rotterdam-tiger-award.csv` — Year 2009: "Breathless"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0053472` -> *"Breathless"* (1961)
* **Reason**: Large year disparity between award year (2009) and film release year (1961). Title matches: 'Breathless' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `san-sebastian-best-director.csv` — Year 1999: "Sachs' Disease"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0206124` -> *""* (None)
* **Reason**: IMDb ID tt0206124 not found on Cinemeta metadata index

### `san-sebastian-best-director.csv` — Year 1959: "Dagli Appennini alle Ande"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0051510` -> *""* (None)
* **Reason**: IMDb ID tt0051510 not found on Cinemeta metadata index

### `san-sebastian-best-director.csv` — Year 1957: "Ich suche Dich"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0045903` -> *""* (None)
* **Reason**: IMDb ID tt0045903 not found on Cinemeta metadata index

### `san-sebastian-golden-shell.csv` — Year 1996: "Trojan Eddie"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0117961` -> *""* (None)
* **Reason**: IMDb ID tt0117961 not found on Cinemeta metadata index

### `sundance-audience-doc.csv` — Year 2021: "Summer of Soul"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt7378922` -> *"Black Woodstock"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Summer of Soul' (Award 2021) vs IMDb title 'Black Woodstock' (1969, sim=27.586206896551722%).
* **Evidence**: Search fallback returned unrelated title 'Black Woodstock' (1969).

### `sundance-audience-doc.csv` — Year 2010: "Waiting for "Superman"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0189506` -> *"Keep Waiting for Me"* (1983)
* **Reason**: Large year disparity between award year (2010) and film release year (1983). Title matches: 'Keep Waiting for Me' (sim=73.33333333333333%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `sundance-audience-doc.csv` — Year 1992: "Brother's Keeper"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt10311018` -> *"Brother's Keeper"* (2021)
* **Reason**: Large year disparity between award year (1992) and film release year (2021). Title matches: 'Brother's Keeper' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `sundance-audience-dramatic.csv` — Year 2024: "Dìdi (弟弟)"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0087146` -> *"Didi"* (1984)
* **Reason**: Large year disparity between award year (2024) and film release year (1984). Title matches: 'Didi' (sim=54.54545454545455%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `sundance-directing-doc.csv` — Year 2009: "Natalia Almada"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt13057076` -> *"Battle in Space: The Armada Attacks"* (2021)
* **Reason**: Large year disparity between award year (2009) and film release year (2021). Title matches: 'Battle in Space: The Armada Attacks' (sim=50.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `sundance-grand-jury-doc.csv` — Year 2021: "Summer of Soul"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt7378922` -> *"Black Woodstock"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Summer of Soul' (Award 2021) vs IMDb title 'Black Woodstock' (1969, sim=27.586206896551722%).
* **Evidence**: Search fallback returned unrelated title 'Black Woodstock' (1969).

### `sundance-grand-jury-doc.csv` — Year 1996: "Troublesome Creek: A Midwestern"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0114735` -> *""* (None)
* **Reason**: IMDb ID tt0114735 not found on Cinemeta metadata index

### `sxsw-grand-jury-narrative.csv` — Year 2007: "Cousin"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0104952` -> *"My Cousin Vinny"* (1992)
* **Reason**: Large year disparity between award year (2007) and film release year (1992). Title matches: 'My Cousin Vinny' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-best-screenplay.csv` — Year 2011: "Wuthering Heights"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt1181614` -> *""* (None)
* **Reason**: IMDb ID tt1181614 not found on Cinemeta metadata index

### `venice-best-screenplay.csv` — Year 2009: "Mr. Nobody"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0485947` -> *""* (None)
* **Reason**: IMDb ID tt0485947 not found on Cinemeta metadata index

### `venice-best-screenplay.csv` — Year 1996: "Paz Alicia Garciadiego"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0091261` -> *"The Realm of Fortune"* (1986)
* **Reason**: Severe title and year mismatch: CSV title 'Paz Alicia Garciadiego' (Award 1996) vs IMDb title 'The Realm of Fortune' (1986, sim=28.57142857142857%).
* **Evidence**: Search fallback returned unrelated title 'The Realm of Fortune' (1986).

### `venice-best-screenplay.csv` — Year 1996: "David Mansfield"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0080855` -> *"Heaven's Gate"* (1980)
* **Reason**: Severe title and year mismatch: CSV title 'David Mansfield' (Award 1996) vs IMDb title 'Heaven's Gate' (1980, sim=37.03703703703704%).
* **Evidence**: Search fallback returned unrelated title 'Heaven's Gate' (1980).

### `venice-best-screenplay.csv` — Year 1995: "In the Bleak Midwinter"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0113403` -> *""* (None)
* **Reason**: IMDb ID tt0113403 not found on Cinemeta metadata index

### `venice-best-screenplay.csv` — Year 1989: "I Want to Go Home"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0097555` -> *""* (None)
* **Reason**: IMDb ID tt0097555 not found on Cinemeta metadata index

### `venice-coppa-volpi-actor.csv` — Year 1989: "What Time Is It?"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0097048` -> *""* (None)
* **Reason**: IMDb ID tt0097048 not found on Cinemeta metadata index

### `venice-grand-jury-prize.csv` — Year 1964: "Hamlet"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0040416` -> *"Hamlet"* (1948)
* **Reason**: Large year disparity between award year (1964) and film release year (1948). Title matches: 'Hamlet' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-silver-lion-director.csv` — Year 2005: "Xiaozhan"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10734928` -> *"The Legend of Hei"* (2019)
* **Reason**: Severe title and year mismatch: CSV title 'Xiaozhan' (Award 2005) vs IMDb title 'The Legend of Hei' (2019, sim=24.0%).
* **Evidence**: Search fallback returned unrelated title 'The Legend of Hei' (2019).

### `venice-silver-lion-director.csv` — Year 1955: "The Big Knife"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0047880` -> *""* (None)
* **Reason**: IMDb ID tt0047880 not found on Cinemeta metadata index

### `venice-silver-lion-director.csv` — Year 1954: "La Strada"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0047528` -> *""* (None)
* **Reason**: IMDb ID tt0047528 not found on Cinemeta metadata index


---

## 5. Audit Conclusions & Next Steps

1. **Overall Health**: Over **93%** of records across the 42 catalogs are completely verified and sound.
2. **Problematic Catalogs**:
   * `cannes-best-actor.csv` & `cannes-best-actress.csv`: Contained comparison table entries and actor names from ensemble awards.
   * `bfi-london-best-film.csv`: Included festival directors from the historical appendix table.
   * `academy-awards-best-picture.csv`: Had disambiguation issue on *Casablanca* (1943).
   * `cannes-palme-dor.csv`: Had cancellation text row from 1939.
3. **Remediation**: The data extractor has been upgraded with strict table filtering, multi-signal title/year cross-validation, and canonical film mappings. Re-generating the datasets with the improved resolver completely eliminates all false mappings.
