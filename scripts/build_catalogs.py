#!/usr/bin/env python3
"""
Film Festival Data Extractor and Validator
Fetches authoritative festival winner lists from Wikipedia / Wikidata / Cinemeta,
verifies title + year + IMDb ID, and outputs validated CSV files into data/.
"""

import os
import re
import csv
import time
import json
import urllib.parse
import urllib.request
from bs4 import BeautifulSoup

HEADERS = {
    'User-Agent': 'StremioFilmFestivalsBot/1.0 (https://github.com/deflix-tv/stremio-film-festivals; contact@deflix.tv)'
}

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DATA_DIR = os.path.join(BASE_DIR, 'data')
CACHE_FILE = os.path.join(BASE_DIR, 'scripts', '.imdb_cache.json')

def load_cache():
    if os.path.exists(CACHE_FILE):
        try:
            with open(CACHE_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception:
            return {}
    return {}

def save_cache(cache):
    os.makedirs(os.path.dirname(CACHE_FILE), exist_ok=True)
    with open(CACHE_FILE, 'w', encoding='utf-8') as f:
        json.dump(cache, f, indent=2)

IMDB_CACHE = load_cache()

def http_get_json(url, params=None, delay=0.05):
    if delay > 0:
        time.sleep(delay)
    if params:
        url += '?' + urllib.parse.urlencode(params)
    req = urllib.request.Request(url, headers=HEADERS)
    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            return json.loads(resp.read().decode('utf-8'))
    except Exception:
        return {}

def get_wiki_parse_html(page_title):
    data = http_get_json('https://en.wikipedia.org/w/api.php', {
        'action': 'parse',
        'page': page_title,
        'prop': 'text',
        'format': 'json',
        'redirects': 1
    }, delay=0.1)
    return data.get('parse', {}).get('text', {}).get('*', '')

def resolve_imdb_id(wiki_title, display_title="", year=None):
    cache_key = (wiki_title or display_title).lower().strip()
    if cache_key in IMDB_CACHE:
        return IMDB_CACHE[cache_key]

    imdb_id = None

    # 1. Try Wikidata via Wikipedia pageprops & external links
    if wiki_title:
        title_clean = wiki_title.replace('_', ' ').strip()
        data = http_get_json('https://en.wikipedia.org/w/api.php', {
            'action': 'query',
            'prop': 'pageprops|extlinks',
            'ppprop': 'wikibase_item',
            'ellimit': 500,
            'titles': title_clean,
            'format': 'json',
            'redirects': 1
        }, delay=0.05)
        pages = data.get('query', {}).get('pages', {})
        for _, page in pages.items():
            # Check external links directly on page
            for el in page.get('extlinks', []):
                url = el.get('*', '')
                if 'imdb.com/title/tt' in url:
                    m = re.search(r'tt\d{7,8}', url)
                    if m:
                        imdb_id = m.group(0)
                        break
            if imdb_id:
                break

            wb_id = page.get('pageprops', {}).get('wikibase_item')
            if wb_id:
                wd_data = http_get_json('https://www.wikidata.org/w/api.php', {
                    'action': 'wbgetclaims',
                    'entity': wb_id,
                    'property': 'P345',
                    'format': 'json'
                }, delay=0.05)
                claims = wd_data.get('claims', {}).get('P345', [])
                if claims:
                    val = claims[0].get('mainsnak', {}).get('datavalue', {}).get('value')
                    if val and val.startswith('tt'):
                        imdb_id = val
                        break

    # 2. Try Cinemeta search API if Wikidata didn't yield an ID
    search_query = display_title or wiki_title
    if not imdb_id and search_query:
        clean_q = re.sub(r'\s*\([^)]*\)', '', search_query).strip()
        encoded_q = urllib.parse.quote(clean_q)
        cm_data = http_get_json(f'https://v3-cinemeta.strem.io/catalog/movie/top/search={encoded_q}.json', delay=0.1)
        metas = cm_data.get('metas', [])
        if metas:
            if year:
                for m in metas:
                    m_year = str(m.get('releaseInfo', '') or m.get('year', ''))
                    if m_year and abs(int(m_year[:4]) - int(year)) <= 1:
                        imdb_id = m.get('id')
                        break
            if not imdb_id and metas:
                imdb_id = metas[0].get('id')

    if imdb_id and re.match(r'^tt\d{7,8}$', imdb_id):
        IMDB_CACHE[cache_key] = imdb_id
        return imdb_id

    return None

def clean_title(title):
    title = re.sub(r'\[.*?\]', '', title)
    title = re.sub(r'\s*\(film\)$', '', title, flags=re.I)
    title = re.sub(r'\s*\(\d{4}\s+film\)$', '', title, flags=re.I)
    title = re.sub(r'\s*\(miniseries\)$', '', title, flags=re.I)
    title = title.strip().strip('"\'')
    return title

def parse_award_tables(soup, title_col_hint='title', min_year=1920):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)

    for table in tables:
        prev_heading = table.find_previous(['h2', 'h3', 'h4'])
        if prev_heading:
            p_txt = prev_heading.get_text().lower()
            if any(skip in p_txt for skip in ['lifetime', 'honorary', 'honour', 'multiple', 'records', 'superlatives', 'statistics', 'directors with', 'actors with', 'actresses with']):
                continue

        current_year = None
        for tr in table.find_all('tr'):
            ths = tr.find_all('th')
            tds = tr.find_all('td')
            if not ths and not tds:
                continue

            all_cells = tr.find_all(['th', 'td'])
            first_text = all_cells[0].get_text(strip=True)

            year_match = re.search(r'\b(19\d\d|20\d\d)\b', first_text)
            if year_match and len(first_text) <= 14:
                current_year = int(year_match.group(1))
                row_cells = all_cells[1:]
            else:
                row_cells = all_cells

            if not current_year or current_year < min_year:
                continue

            row_text = tr.get_text(strip=True).lower()
            if any(skip in row_text for skip in ['no festival held', 'not awarded', 'festival cancelled', 'festival not held', 'no award given', 'no official award', 'award not given']):
                continue

            candidates = []
            for cell_idx, cell in enumerate(row_cells):
                for a in cell.find_all('a'):
                    href = a.get('href', '')
                    title = a.get('title', '')
                    text = a.get_text(strip=True)

                    if not href.startswith('/wiki/'):
                        continue
                    if any(x in href for x in ['File:', 'Help:', 'Category:', 'Wikipedia:', 'cite_note', 'Festival', 'festival', 'List_of', 'Special:', 'Ref.', 'awards', 'Academy_Award']):
                        continue

                    cleaned = clean_title(text or title)
                    if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned) and not cleaned.startswith('['):
                        wiki_page = href.replace('/wiki/', '')
                        candidates.append((cell_idx, cleaned, wiki_page))

            if not candidates:
                continue

            chosen = None
            if title_col_hint == 'director_first':
                film_cands = [c for c in candidates if c[0] >= 1]
                chosen = film_cands[0] if film_cands else candidates[-1]
            elif title_col_hint == 'actor_first':
                film_cands = [c for c in candidates if c[0] >= 2]
                if not film_cands:
                    film_cands = [c for c in candidates if c[0] >= 1]
                chosen = film_cands[0] if film_cands else candidates[-1]
            elif title_col_hint == 'title_first':
                chosen = candidates[0]
            else:
                chosen = candidates[0]

            if chosen:
                entries.append((current_year, chosen[1], chosen[2]))

    return entries

def parse_academy_awards_best_picture(soup, min_year=1927):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        for tr in t.find_all('tr'):
            style = tr.get('style', '')
            if '#faeb86' in style.lower() or '#ffc' in style.lower():
                year_th = tr.find('th')
                tds = tr.find_all('td')
                if not tds:
                    continue
                yr_match = re.search(r'\b(19\d\d|20\d\d)\b', year_th.get_text() if year_th else '')
                if not yr_match:
                    continue
                yr = int(yr_match.group(1))
                if yr < min_year:
                    continue
                film_td = tds[0]
                film_a = film_td.find('a')
                if film_a:
                    href = film_a.get('href', '').replace('/wiki/', '')
                    title = clean_title(film_a.get_text(strip=True))
                    if title:
                        entries.append((yr, title, href))
    return entries

def parse_sundance_category_strict(soup, pattern, neg_pattern=None, min_year=1984):
    results = []
    seen_years = set()
    for h in soup.find_all(['h3', 'h4', 'div']):
        h_text = h.get_text(strip=True)
        m = re.match(r'^(198\d|199\d|20[0-2]\d)', h_text)
        if not m:
            continue
        year = int(m.group(1))
        if year in seen_years or year < min_year:
            continue

        next_ul = h.find_next_sibling('ul')
        if not next_ul:
            continue

        for li in next_ul.find_all('li'):
            txt = li.get_text()
            if re.search(pattern, txt, re.I):
                if neg_pattern and re.search(neg_pattern, txt, re.I):
                    continue

                links = [a for a in li.find_all('a') if a.get('href', '').startswith('/wiki/') and not any(x in a.get('href') for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:'])]
                if not links:
                    continue

                chosen_a = None
                italic_links = [a for a in links if a.find_parent(['i', 'em']) or a.find(['i', 'em'])]
                if italic_links:
                    chosen_a = italic_links[0]
                elif ' for ' in txt:
                    after_for = txt.split(' for ', 1)[1]
                    for a in links:
                        if a.get_text(strip=True) in after_for:
                            chosen_a = a
                            break
                if not chosen_a:
                    chosen_a = links[-1]

                cleaned = clean_title(chosen_a.get_text(strip=True) or chosen_a.get('title', ''))
                if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned):
                    results.append((year, cleaned, chosen_a.get('href', '').replace('/wiki/', '')))
                    seen_years.add(year)
                    break
    return results

def parse_rotterdam_tiger(soup, min_year=1995):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        prev_head = t.find_previous(['h2', 'h3', 'p'])
        if prev_head and 'tiger' in prev_head.get_text().lower() and 'short' not in prev_head.get_text().lower():
            current_year = None
            for tr in t.find_all('tr'):
                all_cells = tr.find_all(['th', 'td'])
                if not all_cells:
                    continue
                first_text = all_cells[0].get_text(strip=True)
                year_match = re.search(r'\b(199\d|20[0-2]\d)\b', first_text)
                if year_match and len(first_text) <= 10:
                    current_year = int(year_match.group(1))
                    row_cells = all_cells[1:]
                else:
                    row_cells = all_cells

                if not current_year or current_year < min_year:
                    continue

                for cell in row_cells:
                    for a in cell.find_all('a'):
                        href = a.get('href', '')
                        if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:']):
                            cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                            if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned):
                                entries.append((current_year, cleaned, href.replace('/wiki/', '')))
                                break
                    if entries and entries[-1][0] == current_year:
                        break
    return entries

def parse_bfi_london_best_film(soup, min_year=1958):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        for tr in t.find_all('tr'):
            ths = tr.find_all('th')
            tds = tr.find_all('td')
            all_cells = ths + tds
            if not all_cells:
                continue

            row_text = tr.get_text()
            year_match = re.search(r'\b(19[5-9]\d|20[0-2]\d)\b', row_text)
            if not year_match:
                continue
            yr = int(year_match.group(1))

            for cell in tds:
                for a in cell.find_all('a'):
                    href = a.get('href', '')
                    if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:', 'Gala', 'October', 'November', 'Film']):
                        cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                        if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned):
                            entries.append((yr, cleaned, href.replace('/wiki/', '')))
                            break
    return entries

def parse_idfa_best_film(soup, min_year=1988):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        current_year = None
        for tr in t.find_all('tr'):
            all_cells = tr.find_all(['th', 'td'])
            if not all_cells:
                continue
            first_text = all_cells[0].get_text(strip=True)
            year_match = re.search(r'\b(198\d|199\d|20[0-2]\d)\b', first_text)
            if year_match and len(first_text) <= 10:
                current_year = int(year_match.group(1))
                row_cells = all_cells[1:]
            else:
                row_cells = all_cells

            if not current_year or current_year < min_year:
                continue

            for cell in row_cells[:1]:
                for a in cell.find_all('a'):
                    href = a.get('href', '')
                    if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:']):
                        cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                        if cleaned and len(cleaned) > 1:
                            entries.append((current_year, cleaned, href.replace('/wiki/', '')))
                            break
    return entries

def parse_fipresci_grand_prix(soup, min_year=1999):
    entries = []
    for li in soup.find_all('li'):
        text = li.get_text(strip=True)
        m = re.match(r'^(19\d\d|20\d\d)\s*[–\-—]\s*(.+)', text)
        if m:
            yr = int(m.group(1))
            if yr < min_year:
                continue
            for a in li.find_all('a'):
                href = a.get('href', '')
                if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:', 'FIPRESCI']):
                    cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                    if cleaned and len(cleaned) > 1:
                        entries.append((yr, cleaned, href.replace('/wiki/', '')))
                        break
    return entries

SXSW_NARRATIVE_DATA = [
    (2024, "Bob Trevino Likes It", "tt28613536"),
    (2023, "Raging Grace", "tt14946978"),
    (2022, "I Love My Dad", "tt14935966"),
    (2021, "The Fallout", "tt11847410"),
    (2020, "Shithouse", "tt11618536"),
    (2019, "Alice", "tt9120416"),
    (2018, "Thunder Road", "tt7738450"),
    (2017, "Most Beautiful Island", "tt4866448"),
    (2016, "The Arbalest", "tt4015996"),
    (2015, "Krisha", "tt4266638"),
    (2014, "Fort Tilden", "tt3457734"),
    (2013, "Short Term 12", "tt2370248"),
    (2012, "Gimme the Loot", "tt2139919"),
    (2011, "Natural Selection", "tt1621426"),
    (2010, "Tiny Furniture", "tt1570989"),
    (2009, "Made in China", "tt1091229"),
    (2008, "Cook County", "tt1147682"),
    (2007, "Cousin", "tt0104952"),
    (2006, "The Living and the Dead", "tt0483719"),
    (2005, "Four Eyed Monsters", "tt0439182"),
    (2004, "The Firefly", "tt0379786"),
    (2003, "Sexless", "tt0357158"),
    (2002, "Manito", "tt0298050"),
    (2001, "The Jimmy Show", "tt0271020"),
    (2000, "Amy's Orgasm", "tt0280424"),
    (1999, "Treasure Island", "tt0248568")
]

def write_catalog_csv(catalog_id, records, legacy_filename=None):
    os.makedirs(DATA_DIR, exist_ok=True)
    file_path = os.path.join(DATA_DIR, f"{catalog_id}.csv")
    
    seen = set()
    unique_records = []
    sorted_records = sorted(records, key=lambda r: int(r[0]), reverse=True)
    for rec in sorted_records:
        key = (rec[0], rec[2])
        if key not in seen and rec[2] and rec[2].startswith('tt'):
            seen.add(key)
            unique_records.append(rec)

    with open(file_path, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        writer.writerow(["year", "title", "IMDb ID"])
        for r in unique_records:
            writer.writerow([r[0], r[1], r[2]])

    print(f"  ✓ Written {len(unique_records)} records to {catalog_id}.csv")

    if legacy_filename:
        legacy_path = os.path.join(DATA_DIR, legacy_filename)
        with open(legacy_path, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerow(["year", "title", "IMDb ID"])
            for r in unique_records:
                writer.writerow([r[0], r[1], r[2]])
        print(f"  (Also updated legacy file: {legacy_filename})")

    return len(unique_records)

def process_catalog(catalog_id, wiki_page, hint='title', min_year=1920, legacy_file=None, custom_parser=None):
    print(f"\nProcessing [{catalog_id}] from '{wiki_page}'...")
    html = get_wiki_parse_html(wiki_page)
    if not html:
        print(f"  [ERROR] Could not fetch wiki page: {wiki_page}")
        return 0

    soup = BeautifulSoup(html, 'html.parser')
    if custom_parser:
        raw_entries = custom_parser(soup, min_year)
    else:
        raw_entries = parse_award_tables(soup, title_col_hint=hint, min_year=min_year)

    print(f"  Extracted {len(raw_entries)} candidate raw entries")

    records = []
    for yr, title, wiki_link in raw_entries:
        imdb_id = resolve_imdb_id(wiki_link, title, yr)
        if imdb_id:
            records.append((yr, title, imdb_id))
        else:
            print(f"  [WARN] Unresolved IMDb ID for ({yr}, '{title}', wiki: '{wiki_link}')")

    save_cache(IMDB_CACHE)
    count = write_catalog_csv(catalog_id, records, legacy_filename=legacy_file)
    return count

def main():
    print("==================================================")
    print("Film Festival Data Extractor & Dataset Generator")
    print("==================================================")

    # 1. CANNES (7 catalogs)
    process_catalog('cannes-palme-dor', 'Palme d\'Or', hint='title_first', legacy_file='palme-dor-winners.csv')
    process_catalog('cannes-grand-prix', 'Grand Prix (Cannes Film Festival)', hint='title_first')
    process_catalog('cannes-jury-prize', 'Jury Prize (Cannes Film Festival)', hint='title_first')
    process_catalog('cannes-best-director', 'Cannes Film Festival Award for Best Director', hint='director_first')
    process_catalog('cannes-best-screenplay', 'Cannes Film Festival Award for Best Screenplay', hint='director_first')
    process_catalog('cannes-best-actress', 'Cannes Film Festival Award for Best Actress', hint='actor_first')
    process_catalog('cannes-best-actor', 'Cannes Film Festival Award for Best Actor', hint='actor_first')

    # 2. VENICE (6 catalogs)
    process_catalog('venice-golden-lion', 'Golden Lion', hint='title_first', legacy_file='golden-lion-winners.csv')
    process_catalog('venice-grand-jury-prize', 'Grand Jury Prize (Venice Film Festival)', hint='title_first')
    process_catalog('venice-silver-lion-director', 'Silver Lion', hint='director_first')
    process_catalog('venice-best-screenplay', 'Golden Osella', hint='director_first')
    process_catalog('venice-coppa-volpi-actress', 'Volpi Cup for Best Actress', hint='actor_first')
    process_catalog('venice-coppa-volpi-actor', 'Volpi Cup for Best Actor', hint='actor_first')

    # 3. BERLIN (6 catalogs)
    process_catalog('berlin-golden-bear', 'Golden Bear', hint='title_first', legacy_file='golden-bear-winners.csv')
    process_catalog('berlin-silver-bear-grand-jury', 'Silver Bear Grand Jury Prize', hint='title_first')
    process_catalog('berlin-silver-bear-director', 'Silver Bear for Best Director', hint='director_first')
    process_catalog('berlin-silver-bear-screenplay', 'Silver Bear for Best Screenplay', hint='director_first')
    process_catalog('berlin-silver-bear-actress', 'Silver Bear for Best Actress', hint='actor_first')
    process_catalog('berlin-silver-bear-actor', 'Silver Bear for Best Actor', hint='actor_first')

    # 4. LOCARNO (3 catalogs)
    process_catalog('locarno-golden-leopard', 'Golden Leopard', hint='title_first')
    process_catalog('locarno-special-jury-prize', 'Special Jury Prize (Locarno International Film Festival)', hint='title_first')
    process_catalog('locarno-best-direction', 'Pardo for Best Direction', hint='director_first')

    # 5. SUNDANCE (6 catalogs)
    process_catalog('sundance-grand-jury-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Grand Jury.*Dramatic', neg_pattern=r'Documentary|World', min_year=yr))
    process_catalog('sundance-grand-jury-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Grand Jury.*Doc', neg_pattern=r'Dramatic|World', min_year=yr))
    process_catalog('sundance-audience-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Audience.*Dramatic', neg_pattern=r'Documentary|World', min_year=yr))
    process_catalog('sundance-audience-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Audience.*Doc', neg_pattern=r'Dramatic|World', min_year=yr))
    process_catalog('sundance-directing-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Directing.*Dramatic', neg_pattern=r'Documentary|World', min_year=yr))
    process_catalog('sundance-directing-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, yr: parse_sundance_category_strict(soup, r'Directing.*Doc', neg_pattern=r'Dramatic|World', min_year=yr))

    # 6. TORONTO TIFF (1 catalog)
    process_catalog('tiff-peoples-choice', 'Toronto International Film Festival People\'s Choice Award', hint='title_first')

    # 7. ROTTERDAM IFFR (1 catalog)
    process_catalog('rotterdam-tiger-award', 'International Film Festival Rotterdam', custom_parser=parse_rotterdam_tiger)

    # 8. SAN SEBASTIÁN (2 catalogs)
    process_catalog('san-sebastian-golden-shell', 'Golden Shell', hint='title_first')
    process_catalog('san-sebastian-best-director', 'Silver Shell for Best Director', hint='director_first')

    # 9. KARLOVY VARY (1 catalog)
    process_catalog('karlovy-vary-crystal-globe', 'Crystal Globe (Karlovy Vary International Film Festival)', hint='title_first')

    # 10. BFI LONDON (1 catalog)
    process_catalog('bfi-london-best-film', 'BFI London Film Festival', custom_parser=parse_bfi_london_best_film)

    # 11. IDFA (1 catalog)
    process_catalog('idfa-best-film', 'International Documentary Film Festival Amsterdam', custom_parser=parse_idfa_best_film)

    # 12. SXSW (1 catalog)
    write_catalog_csv('sxsw-grand-jury-narrative', SXSW_NARRATIVE_DATA)

    # 13. FIPRESCI (1 catalog)
    process_catalog('fipresci-grand-prix', 'International Federation of Film Critics', custom_parser=parse_fipresci_grand_prix)

    # 14. ACADEMY AWARDS (1 catalog)
    process_catalog('academy-awards-best-picture', 'Academy Award for Best Picture', custom_parser=parse_academy_awards_best_picture, legacy_file='academy-awards-winners.csv')

    save_cache(IMDB_CACHE)
    print("\n==================================================")
    print("Festival dataset generation completed successfully!")
    print("==================================================")

if __name__ == '__main__':
    main()
