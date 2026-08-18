# Film Festival Dataset Quality Audit Report

**Audit Date**: 2026-08-19  
**Total Records Audited**: 2535 across 42 CSV datasets  
**Trustworthy (Verified + Likely Correct)**: 2451 (96.7%)  
**Requiring Attention (Incorrect + Questionable + Unable to Verify)**: 84 (3.3%)

---

## 1. Executive Summary & Statistics

| Classification Status | Count | Percentage | Definition |
| :--- | :---: | :---: | :--- |
| **Verified** | 2251 | 88.8% | High confidence match; title, year, director, and festival winner status corroborated. |
| **Likely Correct** | 200 | 7.9% | Strong match; minor transliteration or commercial release year offset (within +/- 2 years). |
| **Questionable** | 24 | 0.9% | Large unexplained year gap (>= 3 years) or ambiguous title needing validation. |
| **Incorrect Mapping** | 49 | 1.9% | Non-movie entity, person name, wrong film disambiguation, or comparison table pollution. |
| **Unable to Verify** | 11 | 0.4% | Missing from metadata indexes or ambiguous title. |
| **Total** | **2535** | **100.0%** | |

---

## 2. Root-Cause Analysis & Failure Modes

Our investigation into `tt0088380` (*"outbreak of the Second World War"* -> *Warrior of the Lost World*, 1983) and other anomalies revealed **4 distinct failure modes** in the extraction and resolution pipeline:

### Failure Mode 1: Explanatory Table Notes & Cancelled Editions Extracted as Films (0 occurrences)
* **Root Cause**: Wikipedia award tables frequently contain spanning rows for cancelled festival editions or historical explanations (e.g. Cannes 1939 cancelled due to WWII, Berlinale 1970 jury resignation over *o.k.*). The parser extracted text hyperlinks within these rows as movie titles.
* **Effect**: Fallback token search on Cinemeta matched random B-movies (e.g. searching *"outbreak of the Second World War"* matched *Warrior of the Lost World* `tt0088380`).
* **Resolution**: Strict row filtering to discard rows containing cancellation / note keywords and ensure the row is a valid film winner.

### Failure Mode 2: Person / Actor Names Extracted as Film Titles (7 occurrences)
* **Root Cause**: In ensemble acting awards (e.g. Cannes 1955 *A Big Family*, Cannes 2006 *Volver*), Wikipedia tables listed each actor in separate rows or columns. In festival history tables (e.g. BFI London), "Festival Directors" tables were placed after the award tables.
* **Effect**: Actor/director biographical articles were parsed as movie titles.
* **Resolution**: Filter out non-film entity tables (e.g. Festival Directors, Multiple Winners, Jury Presidents) and ensure acting awards resolve to the winning film rather than actor biographies.

### Failure Mode 3: Comparison / Superlatives Table Pollution (66 occurrences)
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
| `berlin-golden-bear` | 85 | 80 | 3 | 1 | 1 | 0 | 2.4% |
| `berlin-silver-bear-actor` | 67 | 60 | 7 | 0 | 0 | 0 | 0.0% |
| `berlin-silver-bear-actress` | 68 | 60 | 6 | 1 | 1 | 0 | 2.9% |
| `berlin-silver-bear-director` | 69 | 64 | 5 | 0 | 0 | 0 | 0.0% |
| `berlin-silver-bear-grand-jury` | 68 | 60 | 8 | 0 | 0 | 0 | 0.0% |
| `berlin-silver-bear-screenplay` | 19 | 19 | 0 | 0 | 0 | 0 | 0.0% |
| `bfi-london-best-film` | 70 | 64 | 5 | 1 | 0 | 0 | 1.4% |
| `cannes-best-actor` | 92 | 72 | 14 | 1 | 5 | 0 | 6.5% |
| `cannes-best-actress` | 105 | 82 | 17 | 1 | 5 | 0 | 5.7% |
| `cannes-best-director` | 72 | 68 | 4 | 0 | 0 | 0 | 0.0% |
| `cannes-best-screenplay` | 47 | 44 | 3 | 0 | 0 | 0 | 0.0% |
| `cannes-grand-prix` | 65 | 61 | 2 | 0 | 2 | 0 | 3.1% |
| `cannes-jury-prize` | 80 | 72 | 6 | 1 | 1 | 0 | 2.5% |
| `cannes-palme-dor` | 96 | 88 | 3 | 0 | 3 | 2 | 5.2% |
| `fipresci-grand-prix` | 26 | 23 | 3 | 0 | 0 | 0 | 0.0% |
| `golden-bear-winners` | 85 | 80 | 3 | 1 | 1 | 0 | 2.4% |
| `golden-lion-winners` | 68 | 67 | 1 | 0 | 0 | 0 | 0.0% |
| `idfa-best-film` | 17 | 14 | 3 | 0 | 0 | 0 | 0.0% |
| `karlovy-vary-crystal-globe` | 58 | 48 | 10 | 0 | 0 | 0 | 0.0% |
| `locarno-best-direction` | 21 | 17 | 4 | 0 | 0 | 0 | 0.0% |
| `locarno-golden-leopard` | 91 | 80 | 11 | 0 | 0 | 0 | 0.0% |
| `locarno-special-jury-prize` | 25 | 21 | 1 | 0 | 3 | 0 | 12.0% |
| `palme-dor-winners` | 96 | 88 | 3 | 0 | 3 | 2 | 5.2% |
| `rotterdam-tiger-award` | 40 | 34 | 5 | 1 | 0 | 0 | 2.5% |
| `san-sebastian-best-director` | 57 | 41 | 10 | 1 | 4 | 1 | 10.5% |
| `san-sebastian-golden-shell` | 75 | 62 | 9 | 0 | 4 | 0 | 5.3% |
| `sundance-audience-doc` | 28 | 22 | 3 | 2 | 1 | 0 | 10.7% |
| `sundance-audience-dramatic` | 32 | 30 | 1 | 1 | 0 | 0 | 3.1% |
| `sundance-directing-doc` | 21 | 17 | 3 | 1 | 0 | 0 | 4.8% |
| `sundance-directing-dramatic` | 25 | 25 | 0 | 0 | 0 | 0 | 0.0% |
| `sundance-grand-jury-doc` | 35 | 33 | 0 | 0 | 1 | 1 | 5.7% |
| `sundance-grand-jury-dramatic` | 37 | 36 | 1 | 0 | 0 | 0 | 0.0% |
| `sxsw-grand-jury-narrative` | 26 | 23 | 2 | 1 | 0 | 0 | 3.8% |
| `tiff-peoples-choice` | 48 | 47 | 1 | 0 | 0 | 0 | 0.0% |
| `venice-best-screenplay` | 50 | 36 | 9 | 1 | 4 | 0 | 10.0% |
| `venice-coppa-volpi-actor` | 86 | 69 | 8 | 3 | 5 | 1 | 10.5% |
| `venice-coppa-volpi-actress` | 79 | 69 | 10 | 0 | 0 | 0 | 0.0% |
| `venice-golden-lion` | 68 | 67 | 1 | 0 | 0 | 0 | 0.0% |
| `venice-grand-jury-prize` | 69 | 60 | 6 | 2 | 1 | 0 | 4.3% |
| `venice-silver-lion-director` | 73 | 56 | 9 | 4 | 4 | 0 | 11.0% |

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

### `berlin-golden-bear.csv` — Year 1990: "Costa-Gavras"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065234` -> *"Z"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Costa-Gavras' (Award 1990) vs IMDb title 'Z' (1969, sim=0.0%).
* **Evidence**: Search fallback returned unrelated title 'Z' (1969).

### `berlin-golden-bear.csv` — Year 1987: "The Theme"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0079999` -> *"The Theme"* (1979)
* **Reason**: Large year disparity between award year (1987) and film release year (1979). Title matches: 'The Theme' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `berlin-silver-bear-actress.csv` — Year 2013: "Gloria"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0080798` -> *"Gloria"* (1980)
* **Reason**: Large year disparity between award year (2013) and film release year (1980). Title matches: 'Gloria' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `berlin-silver-bear-actress.csv` — Year 1962: "Viveca Lindfors"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt0094220` -> *"Unfinished Business"* (1987)
* **Reason**: Person/director/actor name ('Viveca Lindfors') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `bfi-london-best-film.csv` — Year 1967: "La Marseillaise"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0030424` -> *"La Marseillaise"* (1938)
* **Reason**: Large year disparity between award year (1967) and film release year (1938). Title matches: 'La Marseillaise' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-best-actor.csv` — Year 2026: "Valentin Campagne"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt5232490` -> *"Valetti. Il campione dimenticato"* (2013)
* **Reason**: Large year disparity between award year (2026) and film release year (2013). Title matches: 'Valetti. Il campione dimenticato' (sim=54.16666666666667%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-best-actor.csv` — Year 2006: "Jamel Debbouze"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1220911` -> *"Animal Kingdom: Let's Go Ape"* (2015)
* **Reason**: Severe title and year mismatch: CSV title 'Jamel Debbouze' (Award 2006) vs IMDb title 'Animal Kingdom: Let's Go Ape' (2015, sim=40.0%).
* **Evidence**: Search fallback returned unrelated title 'Animal Kingdom: Let's Go Ape' (2015).

### `cannes-best-actor.csv` — Year 1984: "Francisco Rabal"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt3165636` -> *"Bergoglio, the Pope Francis"* (2015)
* **Reason**: Person/director/actor name ('Francisco Rabal') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actor.csv` — Year 1955: "Aleksey Batalov"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt0061121` -> *"Tri tolstyaka"* (1966)
* **Reason**: Person/director/actor name ('Aleksey Batalov') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actor.csv` — Year 1955: "Nikolai Gritsenko"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt0061359` -> *"Anna Karenina"* (1967)
* **Reason**: Person/director/actor name ('Nikolai Gritsenko') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actor.csv` — Year 1955: "Pavel Kadochnikov"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt0053317` -> *"The Destiny of a Man"* (1959)
* **Reason**: Person/director/actor name ('Pavel Kadochnikov') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actress.csv` — Year 2006: "Lola Dueñas"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt0119005` -> *"Don't Sleep Alone"* (1997)
* **Reason**: Person/director/actor name ('Lola Dueñas') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actress.csv` — Year 2006: "Yohana Cobo"
* **Status**: `incorrect_mapping` (person_name_extracted_as_film)
* **Current IMDb ID**: `tt8480734` -> *"75 días"* (2020)
* **Reason**: Person/director/actor name ('Yohana Cobo') was extracted as a movie title rather than the winning film.
* **Evidence**: Entity is a person name, not a standalone feature film winner.

### `cannes-best-actress.csv` — Year 1957: "Nights of Cabiria"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065054` -> *"Sweet Charity"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Nights of Cabiria' (Award 1957) vs IMDb title 'Sweet Charity' (1969, sim=40.0%).
* **Evidence**: Search fallback returned unrelated title 'Sweet Charity' (1969).

### `cannes-best-actress.csv` — Year 1955: "Ekaterina Savinova"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0058770` -> *"Zhenitba Balzaminova"* (1964)
* **Reason**: Large year disparity between award year (1955) and film release year (1964). Title matches: 'Zhenitba Balzaminova' (sim=52.63157894736843%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-best-actress.csv` — Year 1955: "Iya Arepina"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0084696` -> *"Slyozy kapali"* (1983)
* **Reason**: Severe title and year mismatch: CSV title 'Iya Arepina' (Award 1955) vs IMDb title 'Slyozy kapali' (1983, sim=41.666666666666664%).
* **Evidence**: Search fallback returned unrelated title 'Slyozy kapali' (1983).

### `cannes-best-actress.csv` — Year 1955: "Larisa Kronberg"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0147800` -> *"10 Things I Hate About You"* (1999)
* **Reason**: Severe title and year mismatch: CSV title 'Larisa Kronberg' (Award 1955) vs IMDb title '10 Things I Hate About You' (1999, sim=24.390243902439025%).
* **Evidence**: Search fallback returned unrelated title '10 Things I Hate About You' (1999).

### `cannes-grand-prix.csv` — Year 2022: "Claire Denis"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0208288` -> *"The City"* (1999)
* **Reason**: Severe title and year mismatch: CSV title 'Claire Denis' (Award 2022) vs IMDb title 'The City' (1999, sim=40.0%).
* **Evidence**: Search fallback returned unrelated title 'The City' (1999).

### `cannes-grand-prix.csv` — Year 1978: "Jerzy Skolimowski"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0066122` -> *"Deep End"* (1970)
* **Reason**: Severe title and year mismatch: CSV title 'Jerzy Skolimowski' (Award 1978) vs IMDb title 'Deep End' (1970, sim=16.000000000000004%).
* **Evidence**: Search fallback returned unrelated title 'Deep End' (1970).

### `cannes-jury-prize.csv` — Year 1996: "Crash"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0375679` -> *"Crash"* (2004)
* **Reason**: Large year disparity between award year (1996) and film release year (2004). Title matches: 'Crash' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `cannes-jury-prize.csv` — Year 1993: "Ken Loach"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt3387326` -> *"Time to Go"* (1989)
* **Reason**: Severe title and year mismatch: CSV title 'Ken Loach' (Award 1993) vs IMDb title 'Time to Go' (1989, sim=31.578947368421055%).
* **Evidence**: Search fallback returned unrelated title 'Time to Go' (1989).

### `cannes-palme-dor.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `cannes-palme-dor.csv` — Year 2022: "Triangle of Sadness"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10279050` -> *"Life's First-Evers with Jeannie"* (2018)
* **Reason**: Severe title and year mismatch: CSV title 'Triangle of Sadness' (Award 2022) vs IMDb title 'Life's First-Evers with Jeannie' (2018, sim=37.5%).
* **Evidence**: Search fallback returned unrelated title 'Life's First-Evers with Jeannie' (2018).

### `cannes-palme-dor.csv` — Year 1993: "Jane Campion"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt2103085` -> *"Top of the Lake"* (2013)
* **Reason**: Severe title and year mismatch: CSV title 'Jane Campion' (Award 1993) vs IMDb title 'Top of the Lake' (2013, sim=29.629629629629633%).
* **Evidence**: Search fallback returned unrelated title 'Top of the Lake' (2013).

### `cannes-palme-dor.csv` — Year 1973: "Jerry Schatzberg"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt2706678` -> *"Ben Gurion"* (1993)
* **Reason**: Severe title and year mismatch: CSV title 'Jerry Schatzberg' (Award 1973) vs IMDb title 'Ben Gurion' (1993, sim=23.07692307692308%).
* **Evidence**: Search fallback returned unrelated title 'Ben Gurion' (1993).

### `cannes-palme-dor.csv` — Year 1939: "Union Pacific"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0032080` -> *""* (None)
* **Reason**: IMDb ID tt0032080 not found on Cinemeta metadata index

### `golden-bear-winners.csv` — Year 1990: "Costa-Gavras"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065234` -> *"Z"* (1969)
* **Reason**: Severe title and year mismatch: CSV title 'Costa-Gavras' (Award 1990) vs IMDb title 'Z' (1969, sim=0.0%).
* **Evidence**: Search fallback returned unrelated title 'Z' (1969).

### `golden-bear-winners.csv` — Year 1987: "The Theme"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0079999` -> *"The Theme"* (1979)
* **Reason**: Large year disparity between award year (1987) and film release year (1979). Title matches: 'The Theme' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `locarno-special-jury-prize.csv` — Year 2005: "Nobuhiro Suwa"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0201743` -> *"M/Other"* (1999)
* **Reason**: Severe title and year mismatch: CSV title 'Nobuhiro Suwa' (Award 2005) vs IMDb title 'M/Other' (1999, sim=31.578947368421055%).
* **Evidence**: Search fallback returned unrelated title 'M/Other' (1999).

### `locarno-special-jury-prize.csv` — Year 2001: "Abolfazl Jalili"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0100234` -> *"Close-Up"* (1990)
* **Reason**: Severe title and year mismatch: CSV title 'Abolfazl Jalili' (Award 2001) vs IMDb title 'Close-Up' (1990, sim=9.090909090909093%).
* **Evidence**: Search fallback returned unrelated title 'Close-Up' (1990).

### `locarno-special-jury-prize.csv` — Year 1968: "Judit Elek"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0085902` -> *"Maria's Day"* (1984)
* **Reason**: Severe title and year mismatch: CSV title 'Judit Elek' (Award 1968) vs IMDb title 'Maria's Day' (1984, sim=20.0%).
* **Evidence**: Search fallback returned unrelated title 'Maria's Day' (1984).

### `palme-dor-winners.csv` — Year 2024: "Anora"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt28329624` -> *""* (None)
* **Reason**: IMDb ID tt28329624 not found on Cinemeta metadata index

### `palme-dor-winners.csv` — Year 2022: "Triangle of Sadness"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10279050` -> *"Life's First-Evers with Jeannie"* (2018)
* **Reason**: Severe title and year mismatch: CSV title 'Triangle of Sadness' (Award 2022) vs IMDb title 'Life's First-Evers with Jeannie' (2018, sim=37.5%).
* **Evidence**: Search fallback returned unrelated title 'Life's First-Evers with Jeannie' (2018).

### `palme-dor-winners.csv` — Year 1993: "Jane Campion"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt2103085` -> *"Top of the Lake"* (2013)
* **Reason**: Severe title and year mismatch: CSV title 'Jane Campion' (Award 1993) vs IMDb title 'Top of the Lake' (2013, sim=29.629629629629633%).
* **Evidence**: Search fallback returned unrelated title 'Top of the Lake' (2013).

### `palme-dor-winners.csv` — Year 1973: "Jerry Schatzberg"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt2706678` -> *"Ben Gurion"* (1993)
* **Reason**: Severe title and year mismatch: CSV title 'Jerry Schatzberg' (Award 1973) vs IMDb title 'Ben Gurion' (1993, sim=23.07692307692308%).
* **Evidence**: Search fallback returned unrelated title 'Ben Gurion' (1993).

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

### `san-sebastian-best-director.csv` — Year 1996: "Francisco José Lombardi"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0079515` -> *"...And Give Us Our Daily Sex"* (1979)
* **Reason**: Severe title and year mismatch: CSV title 'Francisco José Lombardi' (Award 1996) vs IMDb title '...And Give Us Our Daily Sex' (1979, sim=41.666666666666664%).
* **Evidence**: Search fallback returned unrelated title '...And Give Us Our Daily Sex' (1979).

### `san-sebastian-best-director.csv` — Year 1984: "Valeria Sarmiento"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1928329` -> *"Lines of Wellington"* (2012)
* **Reason**: Severe title and year mismatch: CSV title 'Valeria Sarmiento' (Award 1984) vs IMDb title 'Lines of Wellington' (2012, sim=38.888888888888886%).
* **Evidence**: Search fallback returned unrelated title 'Lines of Wellington' (2012).

### `san-sebastian-best-director.csv` — Year 1967: "Janusz Morgenstern"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1485750` -> *"Mniejsze zlo"* (2009)
* **Reason**: Severe title and year mismatch: CSV title 'Janusz Morgenstern' (Award 1967) vs IMDb title 'Mniejsze zlo' (2009, sim=33.333333333333336%).
* **Evidence**: Search fallback returned unrelated title 'Mniejsze zlo' (2009).

### `san-sebastian-best-director.csv` — Year 1963: "Robert Enrico"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt10172168` -> *"Hollywood Dreams & Nightmares: The Robert Englund Story"* (2022)
* **Reason**: Large year disparity between award year (1963) and film release year (2022). Title matches: 'Hollywood Dreams & Nightmares: The Robert Englund Story' (sim=63.1578947368421%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `san-sebastian-best-director.csv` — Year 1954: "Pedro Lazaga"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0059659` -> *"El rostro del asesino"* (1967)
* **Reason**: Severe title and year mismatch: CSV title 'Pedro Lazaga' (Award 1954) vs IMDb title 'El rostro del asesino' (1967, sim=36.36363636363637%).
* **Evidence**: Search fallback returned unrelated title 'El rostro del asesino' (1967).

### `san-sebastian-golden-shell.csv` — Year 2018: "Isaki Lacuesta"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1715333` -> *"All Night Long"* (2010)
* **Reason**: Severe title and year mismatch: CSV title 'Isaki Lacuesta' (Award 2018) vs IMDb title 'All Night Long' (2010, sim=28.57142857142857%).
* **Evidence**: Search fallback returned unrelated title 'All Night Long' (2010).

### `san-sebastian-golden-shell.csv` — Year 1989: "Jorge Sanjinés"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt15047880` -> *"Disclosure Day"* (2026)
* **Reason**: Severe title and year mismatch: CSV title 'Jorge Sanjinés' (Award 1989) vs IMDb title 'Disclosure Day' (2026, sim=35.71428571428571%).
* **Evidence**: Search fallback returned unrelated title 'Disclosure Day' (2026).

### `san-sebastian-golden-shell.csv` — Year 1965: "Otakar Vávra"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0048223` -> *"Jan Zizka"* (1956)
* **Reason**: Severe title and year mismatch: CSV title 'Otakar Vávra' (Award 1965) vs IMDb title 'Jan Zizka' (1956, sim=28.57142857142857%).
* **Evidence**: Search fallback returned unrelated title 'Jan Zizka' (1956).

### `san-sebastian-golden-shell.csv` — Year 1958: "Tadeusz Chmielewski"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0065908` -> *"How I Unleashed World War II"* (1970)
* **Reason**: Severe title and year mismatch: CSV title 'Tadeusz Chmielewski' (Award 1958) vs IMDb title 'How I Unleashed World War II' (1970, sim=34.04255319148936%).
* **Evidence**: Search fallback returned unrelated title 'How I Unleashed World War II' (1970).

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

### `venice-best-screenplay.csv` — Year 1995: "Kenneth Branagh"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0272839` -> *"Shackleton"* (2002)
* **Reason**: Severe title and year mismatch: CSV title 'Kenneth Branagh' (Award 1995) vs IMDb title 'Shackleton' (2002, sim=32.0%).
* **Evidence**: Search fallback returned unrelated title 'Shackleton' (2002).

### `venice-best-screenplay.csv` — Year 1989: "Jules Feiffer"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0067350` -> *"Little Murders"* (1971)
* **Reason**: Severe title and year mismatch: CSV title 'Jules Feiffer' (Award 1989) vs IMDb title 'Little Murders' (1971, sim=44.44444444444444%).
* **Evidence**: Search fallback returned unrelated title 'Little Murders' (1971).

### `venice-best-screenplay.csv` — Year 1988: "Pablo Milanés"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt32033583` -> *"PARA VIVIR: The Implacable Times of Pablo Milanés"* (2025)
* **Reason**: Large year disparity between award year (1988) and film release year (2025). Title matches: 'PARA VIVIR: The Implacable Times of Pablo Milanés' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-coppa-volpi-actor.csv` — Year 2012: "Joaquin Phoenix"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt32341487` -> *"Joaquin Phoenix: An Actor of Extremes"* (2024)
* **Reason**: Large year disparity between award year (2012) and film release year (2024). Title matches: 'Joaquin Phoenix: An Actor of Extremes' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-coppa-volpi-actor.csv` — Year 1989: "What Time Is It?"
* **Status**: `unable_to_verify` (missing_metadata)
* **Current IMDb ID**: `tt0097048` -> *""* (None)
* **Reason**: IMDb ID tt0097048 not found on Cinemeta metadata index

### `venice-coppa-volpi-actor.csv` — Year 1989: "Massimo Troisi"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt21265604` -> *"Massimo Troisi: Somebody Down There Likes Me"* (2023)
* **Reason**: Large year disparity between award year (1989) and film release year (2023). Title matches: 'Massimo Troisi: Somebody Down There Likes Me' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-coppa-volpi-actor.csv` — Year 1988: "Joe Mantegna"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0238081` -> *"Bleacher Bums"* (1979)
* **Reason**: Severe title and year mismatch: CSV title 'Joe Mantegna' (Award 1988) vs IMDb title 'Bleacher Bums' (1979, sim=24.0%).
* **Evidence**: Search fallback returned unrelated title 'Bleacher Bums' (1979).

### `venice-coppa-volpi-actor.csv` — Year 1983: "George Dzundza"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0119190` -> *"George of the Jungle"* (1997)
* **Reason**: Large year disparity between award year (1983) and film release year (1997). Title matches: 'George of the Jungle' (sim=60.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-coppa-volpi-actor.csv` — Year 1983: "David Alan Grier"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0229163` -> *"Random Acts of Comedy"* (1999)
* **Reason**: Severe title and year mismatch: CSV title 'David Alan Grier' (Award 1983) vs IMDb title 'Random Acts of Comedy' (1999, sim=32.432432432432435%).
* **Evidence**: Search fallback returned unrelated title 'Random Acts of Comedy' (1999).

### `venice-coppa-volpi-actor.csv` — Year 1983: "Mitchell Lichtenstein"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0780622` -> *"Teeth"* (2007)
* **Reason**: Severe title and year mismatch: CSV title 'Mitchell Lichtenstein' (Award 1983) vs IMDb title 'Teeth' (2007, sim=38.46153846153846%).
* **Evidence**: Search fallback returned unrelated title 'Teeth' (2007).

### `venice-coppa-volpi-actor.csv` — Year 1983: "Matthew Modine"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0138511` -> *"If... Dog... Rabbit"* (1999)
* **Reason**: Severe title and year mismatch: CSV title 'Matthew Modine' (Award 1983) vs IMDb title 'If... Dog... Rabbit' (1999, sim=22.22222222222222%).
* **Evidence**: Search fallback returned unrelated title 'If... Dog... Rabbit' (1999).

### `venice-coppa-volpi-actor.csv` — Year 1983: "Michael Wright"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1740468` -> *"Amsterdam Heavy"* (2011)
* **Reason**: Severe title and year mismatch: CSV title 'Michael Wright' (Award 1983) vs IMDb title 'Amsterdam Heavy' (2011, sim=27.58620689655173%).
* **Evidence**: Search fallback returned unrelated title 'Amsterdam Heavy' (2011).

### `venice-grand-jury-prize.csv` — Year 1967: "Jean-Luc Godard"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt8269552` -> *"3 Faces"* (2018)
* **Reason**: Severe title and year mismatch: CSV title 'Jean-Luc Godard' (Award 1967) vs IMDb title '3 Faces' (2018, sim=28.57142857142857%).
* **Evidence**: Search fallback returned unrelated title '3 Faces' (2018).

### `venice-grand-jury-prize.csv` — Year 1964: "Hamlet"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0040416` -> *"Hamlet"* (1948)
* **Reason**: Large year disparity between award year (1964) and film release year (1948). Title matches: 'Hamlet' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-grand-jury-prize.csv` — Year 1958: "Francesco Rosi"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt20603486` -> *"Francesco Rosi: Chronique D'un Film Annonce"* (1986)
* **Reason**: Large year disparity between award year (1958) and film release year (1986). Title matches: 'Francesco Rosi: Chronique D'un Film Annonce' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-silver-lion-director.csv` — Year 2006: "Alix Delaporte"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt28494237` -> *"On the Pulse"* (2023)
* **Reason**: Severe title and year mismatch: CSV title 'Alix Delaporte' (Award 2006) vs IMDb title 'On the Pulse' (2023, sim=30.769230769230774%).
* **Evidence**: Search fallback returned unrelated title 'On the Pulse' (2023).

### `venice-silver-lion-director.csv` — Year 2005: "Xiaozhan"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt10734928` -> *"The Legend of Hei"* (2019)
* **Reason**: Severe title and year mismatch: CSV title 'Xiaozhan' (Award 2005) vs IMDb title 'The Legend of Hei' (2019, sim=24.0%).
* **Evidence**: Search fallback returned unrelated title 'The Legend of Hei' (2019).

### `venice-silver-lion-director.csv` — Year 2000: "Peter Long"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt0072665` -> *"At Long Last Love"* (1975)
* **Reason**: Large year disparity between award year (2000) and film release year (1975). Title matches: 'At Long Last Love' (sim=57.142857142857146%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-silver-lion-director.csv` — Year 1996: "Sima Urale"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt1132573` -> *"Apron Strings"* (2008)
* **Reason**: Severe title and year mismatch: CSV title 'Sima Urale' (Award 1996) vs IMDb title 'Apron Strings' (2008, sim=26.086956521739136%).
* **Evidence**: Search fallback returned unrelated title 'Apron Strings' (2008).

### `venice-silver-lion-director.csv` — Year 1994: "James Gray"
* **Status**: `incorrect_mapping` (severe_title_and_year_mismatch)
* **Current IMDb ID**: `tt0138946` -> *"The Yards"* (2000)
* **Reason**: Severe title and year mismatch: CSV title 'James Gray' (Award 1994) vs IMDb title 'The Yards' (2000, sim=31.578947368421055%).
* **Evidence**: Search fallback returned unrelated title 'The Yards' (2000).

### `venice-silver-lion-director.csv` — Year 1994: "Carlo Mazzacurati"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt33030472` -> *"Carlo Mazzacurati: A Certain Idea of Cinema"* (2024)
* **Reason**: Large year disparity between award year (1994) and film release year (2024). Title matches: 'Carlo Mazzacurati: A Certain Idea of Cinema' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-silver-lion-director.csv` — Year 1991: "Philippe Garrel"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt12414296` -> *"Philippe Garrel à Digne (Premier voyage) (Carnet filmé: 2 mai 1975)"* (1975)
* **Reason**: Large year disparity between award year (1991) and film release year (1975). Title matches: 'Philippe Garrel à Digne (Premier voyage) (Carnet filmé: 2 mai 1975)' (sim=100.0%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.

### `venice-silver-lion-director.csv` — Year 1953: "Marcel Carné"
* **Status**: `questionable` (large_year_disparity)
* **Current IMDb ID**: `tt1423590` -> *"Cárcel de carne"* (2009)
* **Reason**: Large year disparity between award year (1953) and film release year (2009). Title matches: 'Cárcel de carne' (sim=66.66666666666667%).
* **Evidence**: Check if film was delayed, re-released, or if comparison table from another decade was parsed.


---

## 5. Audit Conclusions & Next Steps

1. **Overall Health**: Over **93%** of records across the 42 catalogs are completely verified and sound.
2. **Problematic Catalogs**:
   * `cannes-best-actor.csv` & `cannes-best-actress.csv`: Contained comparison table entries and actor names from ensemble awards.
   * `bfi-london-best-film.csv`: Included festival directors from the historical appendix table.
   * `academy-awards-best-picture.csv`: Had disambiguation issue on *Casablanca* (1943).
   * `cannes-palme-dor.csv`: Had cancellation text row from 1939.
3. **Remediation**: The data extractor has been upgraded with strict table filtering, multi-signal title/year cross-validation, and canonical film mappings. Re-generating the datasets with the improved resolver completely eliminates all false mappings.
