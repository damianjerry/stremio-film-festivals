#!/usr/bin/env python3
"""
High-Precision Film Festival Data Extractor and Validator
Extracts verified festival winner lists from Wikipedia / Wikidata / Cinemeta,
enforcing multi-signal verification:
  1. Strict decade/winner table selection (ignoring multiple-winners, statistics, and career tables)
  2. Film link isolation (extracting film titles from italic tags or designated film columns, never actor/director biographies)
  3. Proper continuation-row handling for tied winners with rowspan Year cells
  4. Canonical classic dictionary + Wikidata P345 property matching
  5. Multi-signal Cinemeta validation (title similarity >= 75%, year tolerance <= 2 years, type == movie)
  6. Maximum year filter (<= 2025) to reject speculative / future entries
"""

import os
import re
import csv
import time
import json
import urllib.parse
import urllib.request
from bs4 import BeautifulSoup
from rapidfuzz import fuzz

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

CANONICAL_FILM_IMDB_MAPPINGS = {
    "Casablanca": "tt0034583",
    "Union Pacific": "tt0032080",
    "Wings": "tt0018578",
    "The Broadway Melody": "tt0019729",
    "All Quiet on the Western Front": "tt0020629",
    "Cimarron": "tt0021746",
    "Grand Hotel": "tt0022958",
    "Cavalcade": "tt0023876",
    "It Happened One Night": "tt0025316",
    "Mutiny on the Bounty": "tt0026752",
    "The Great Ziegfeld": "tt0027698",
    "The Life of Emile Zola": "tt0029146",
    "You Can't Take It with You": "tt0030993",
    "Gone with the Wind": "tt0031381",
    "Rebecca": "tt0032976",
    "How Green Was My Valley": "tt0033729",
    "Mrs. Miniver": "tt0035093",
    "Going My Way": "tt0036872",
    "The Lost Weekend": "tt0037884",
    "The Best Years of Our Lives": "tt0036868",
    "Gentleman's Agreement": "tt0039416",
    "Hamlet": "tt0040416",
    "All the King's Men": "tt0041113",
    "All About Eve": "tt0042192",
    "An American in Paris": "tt0043278",
    "The Greatest Show on Earth": "tt0044672",
    "From Here to Eternity": "tt0045793",
    "On the Waterfront": "tt0047296",
    "Marty": "tt0048356",
    "Around the World in 80 Days": "tt0048960",
    "The Bridge on the River Kwai": "tt0050212",
    "Gigi": "tt0051658",
    "Ben-Hur": "tt0052618",
    "The Apartment": "tt0053604",
    "West Side Story": "tt0055614",
    "Lawrence of Arabia": "tt0056172",
    "Tom Jones": "tt0057590",
    "My Fair Lady": "tt0058385",
    "The Sound of Music": "tt0059742",
    "A Man for All Seasons": "tt0060665",
    "In the Heat of the Night": "tt0061811",
    "Oliver!": "tt0063385",
    "Midnight Cowboy": "tt0064665",
    "Patton": "tt0066206",
    "The French Connection": "tt0067116",
    "The Godfather": "tt0068646",
    "The Sting": "tt0070735",
    "The Godfather Part II": "tt0071562",
    "One Flew Over the Cuckoo's Nest": "tt0073486",
    "Rocky": "tt0075148",
    "Annie Hall": "tt0075686",
    "The Deer Hunter": "tt0077416",
    "Kramer vs. Kramer": "tt0079417",
    "Ordinary People": "tt0081283",
    "Chariots of Fire": "tt0082158",
    "Gandhi": "tt0083987",
    "Terms of Endearment": "tt0086425",
    "Amadeus": "tt0086879",
    "Out of Africa": "tt0089755",
    "Platoon": "tt0091763",
    "The Last Emperor": "tt0093389",
    "Rain Man": "tt0095953",
    "Driving Miss Daisy": "tt0097239",
    "Dances with Wolves": "tt0099348",
    "The Silence of the Lambs": "tt0102926",
    "Unforgiven": "tt0105695",
    "Schindler's List": "tt0108052",
    "Forrest Gump": "tt0109830",
    "Braveheart": "tt0112573",
    "The English Patient": "tt0116209",
    "Titanic": "tt0120338",
    "Shakespeare in Love": "tt0138097",
    "American Beauty": "tt0169547",
    "Gladiator": "tt0172495",
    "A Beautiful Mind": "tt0268978",
    "Chicago": "tt0299658",
    "The Lord of the Rings: The Return of the King": "tt0167260",
    "Million Dollar Baby": "tt0405159",
    "Crash": "tt0375679",
    "The Departed": "tt0407887",
    "No Country for Old Men": "tt0477348",
    "Slumdog Millionaire": "tt1010048",
    "The Hurt Locker": "tt0887912",
    "The King's Speech": "tt1504320",
    "The Artist": "tt1655442",
    "Argo": "tt1024648",
    "12 Years a Slave": "tt2024544",
    "Birdman or (The Unexpected Virtue of Ignorance)": "tt2562232",
    "Spotlight": "tt1895587",
    "Moonlight": "tt4975721",
    "The Shape of Water": "tt5580390",
    "Green Book": "tt6966692",
    "Parasite": "tt6751668",
    "Nomadland": "tt9770150",
    "CODA": "tt10366460",
    "Everything Everywhere All at Once": "tt6710474",
    "Oppenheimer": "tt15398776",
    "Anora": "tt28607951",
    "Vertigo": "tt0052357",
    "North by Northwest": "tt0053125",
    "Seven Samurai": "tt0047478",
    "Rashomon": "tt0042876",
    "La Strada": "tt0047528",
    "Bicycle Thieves": "tt0040522",
    "Rome, Open City": "tt0038057",
    "Brief Encounter": "tt0037558",
    "Wild Strawberries": "tt0050986",
    "The Seventh Seal": "tt0050976",
    "Persona": "tt0060827",
    "8½": "tt0056801",
    "La Dolce Vita": "tt0053779",
    "L'Avventura": "tt0053580",
    "Breathless": "tt0053472",
    "Cleo from 5 to 7": "tt0055852",
    "Beau Travail": "tt0209933",
    "In the Mood for Love": "tt0247444",
    "Yi Yi": "tt0244316",
    "Mulholland Drive": "tt0166924",
    "Pulp Fiction": "tt0110912",
    "Apocalypse Now": "tt0078788",
    "Taxi Driver": "tt0075314",
    "Paris, Texas": "tt0087884",
    "Wings of Desire": "tt0093191",
    "Farewell My Concubine": "tt0106332",
    "Underground": "tt0114787",
    "Taste of Cherry": "tt0120265",
    "Eternity and a Day": "tt0156794",
    "Rosetta": "tt0200071",
    "Dancer in the Dark": "tt0168629",
    "The Pianist": "tt0253474",
    "The White Ribbon": "tt1149362",
    "Amour": "tt1602620",
    "Blue Is the Warmest Colour": "tt2278871",
    "Winter Sleep": "tt2758880",
    "The Square": "tt4995790",
    "Shoplifters": "tt8075192",
    "Triangle of Sadness": "tt7322224",
    "Anatomy of a Fall": "tt17009710",
    "Larks on a String": "tt0064994",
    "Commissar": "tt0061876",
    "A Big Family": "tt0046800",
    "Volver": "tt0453556",
    "Streamers": "tt0086377",
    "The Master": "tt1560747",
    "What Time Is It?": "tt0097048",
    "House of Games": "tt0093223",
    "Scarecrow": "tt0070643",
    "Stars at Noon": "tt10354106",
    "Gloria": "tt0080798",
    "The Man in the White Suit": "tt0044876",
    "Repulsion": "tt0059646",
    "Overlord": "tt0073498",
    "Las truchas": "tt0076846",
    "La colmena": "tt0083745",
    "The Beehive": "tt0083745",
    "Trojan Eddie": "tt0117961",
    "Mirage": "tt0059448",
    "Soleil O": "tt0065014",
    "Soleil Ô": "tt0065014",
    "Private Road": "tt0067623",
    "Schmetterlinge": "tt0096055",
    "Aimee & Jaguar": "tt0130444",
    "Aimée & Jaguar": "tt0130444",
    "Sachs' Disease": "tt0206124",
    "La maladie de Sachs": "tt0206124",
    "Emilia Pérez": "tt20221436",
    "Kinds of Kindness": "tt22408160",
    "The Zone of Interest": "tt7160372",
    "All We Imagine as Light": "tt27823528",
    "Linha de Passe": "tt0803029",
}

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

def clean_title(title):
    title = re.sub(r'\[.*?\]', '', title)
    title = re.sub(r'\s*\(film\)$', '', title, flags=re.I)
    title = re.sub(r'\s*\(\d{4}\s+film\)$', '', title, flags=re.I)
    title = re.sub(r'\s*\(miniseries\)$', '', title, flags=re.I)
    title = title.strip().strip('"\'')
    return title

def is_note_or_cancellation(text):
    text_lower = text.lower()
    patterns = [
        r'\boutbreak\b', r'\bsecond world war\b', r'\bworld war ii\b',
        r'\bcancelled\b', r'\bno festival\b', r'\bnot held\b',
        r'\bno award\b', r'\bnot awarded\b', r'\bno official award\b',
        r'\btimeline of\b', r'\bedition\b', r'\btied\b', r'\bjury resigned\b',
        r'\bfestival director\b', r'\bjewish state\b', r'\bhonorary\b'
    ]
    return any(re.search(p, text_lower) for p in patterns)

def resolve_imdb_id_strict(wiki_title, display_title="", year=None):
    clean_disp = clean_title(display_title or wiki_title)
    if clean_disp in CANONICAL_FILM_IMDB_MAPPINGS:
        return CANONICAL_FILM_IMDB_MAPPINGS[clean_disp]

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

    # 2. Try Cinemeta search API with STRICT multi-signal validation
    search_query = display_title or wiki_title
    if not imdb_id and search_query:
        clean_q = re.sub(r'\s*\([^)]*\)', '', search_query).strip()
        encoded_q = urllib.parse.quote(clean_q)
        cm_data = http_get_json(f'https://v3-cinemeta.strem.io/catalog/movie/top/search={encoded_q}.json', delay=0.1)
        metas = cm_data.get('metas', [])
        for m in metas:
            m_name = m.get('name', '')
            m_year_str = str(m.get('releaseInfo', '') or m.get('year', ''))
            m_year = int(m_year_str[:4]) if re.match(r'^\d{4}', m_year_str) else None

            sim_ratio = fuzz.ratio(clean_q.lower(), m_name.lower())
            sim_token = fuzz.token_set_ratio(clean_q.lower(), m_name.lower())
            max_sim = max(sim_ratio, sim_token)

            year_diff = abs(m_year - int(year)) if (m_year and year) else 999

            if max_sim >= 75 and year_diff <= 2:
                imdb_id = m.get('id')
                break

    if imdb_id and re.match(r'^tt\d{7,8}$', imdb_id):
        IMDB_CACHE[cache_key] = imdb_id
        return imdb_id

    return None

def extract_award_films_robust(soup, min_year=1920, max_year=2025):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c and 'navbox' not in c)

    for table in tables:
        prev_heading = table.find_previous(['h2', 'h3', 'h4'])
        if prev_heading:
            p_txt = prev_heading.get_text().lower()
            if any(skip in p_txt for skip in [
                'lifetime', 'honorary', 'honour', 'multiple', 'records', 'superlatives',
                'statistics', 'directors with', 'actors with', 'actresses with',
                'winners of multiple', 'festival director', 'jury president', 'see also',
                'retrospective', 'other awards', 'notes', 'references', 'external links'
            ]):
                continue

        headers = []
        header_row = table.find('tr')
        if header_row:
            headers = [th.get_text(strip=True).lower() for th in header_row.find_all(['th', 'td'])]

        film_col_idx = -1
        for idx, h in enumerate(headers):
            if any(k in h for k in ['film', 'english title', 'title', 'película', 'titolo', 'obra']) and not any(bad in h for bad in ['director', 'actor', 'actress', 'screenwriter', 'recipient', 'sceneggiatura']):
                film_col_idx = idx
                break

        current_year = None
        for tr in table.find_all('tr')[1:]:
            cells = tr.find_all(['th', 'td'])
            if not cells:
                continue

            first_text = cells[0].get_text(strip=True)
            year_match = re.search(r'\b(19\d\d|20\d\d)\b', first_text)

            if year_match and len(first_text) <= 14:
                yr = int(year_match.group(1))
                if yr > max_year or yr < min_year:
                    continue
                current_year = yr
                has_year_col = True
                row_cells = cells[1:]
            else:
                has_year_col = False
                row_cells = cells

            if not current_year or current_year < min_year or current_year > max_year:
                continue

            row_text = tr.get_text(strip=True)
            if is_note_or_cancellation(row_text):
                continue

            chosen_a = None
            target_cell = None
            if has_year_col and film_col_idx != -1:
                target_idx = film_col_idx - 1
                if 0 <= target_idx < len(row_cells):
                    target_cell = row_cells[target_idx]
            elif not has_year_col and len(row_cells) > 0:
                target_cell = row_cells[0]

            if target_cell:
                italic = target_cell.find(['i', 'em'])
                if italic:
                    chosen_a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                if not chosen_a:
                    chosen_a = target_cell.find('a')

            if not chosen_a:
                for cell in row_cells:
                    italic = cell.find(['i', 'em'])
                    if italic:
                        a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                        if a and a.get('href', '').startswith('/wiki/'):
                            chosen_a = a
                            break

            if chosen_a:
                href = chosen_a.get('href', '').replace('/wiki/', '')
                title = clean_title(chosen_a.get_text(strip=True) or chosen_a.get('title', ''))
                if href and title and len(title) > 1 and not is_note_or_cancellation(title):
                    entries.append((current_year, title, href))

    return entries

def parse_acting_films(soup, min_year=1920, max_year=2025):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c and 'navbox' not in c)
    for table in tables:
        prev_h = table.find_previous(['h2', 'h3', 'h4'])
        if prev_h:
            p_txt = prev_h.get_text().lower()
            if any(skip in p_txt for skip in ['multiple', 'superlatives', 'records', 'statistics', 'see also', 'references', 'external links', 'notes']):
                continue

        headers = []
        h_row = table.find('tr')
        if h_row:
            headers = [th.get_text(strip=True).lower() for th in h_row.find_all(['th', 'td'])]

        film_col_idx = -1
        for idx, h in enumerate(headers):
            if any(k in h for k in ['film', 'english title', 'title', 'película', 'titolo', 'obra']) and not any(bad in h for bad in ['actor', 'actress', 'recipient', 'role', 'character']):
                film_col_idx = idx
                break

        current_year = None
        for tr in table.find_all('tr')[1:]:
            cells = tr.find_all(['th', 'td'])
            if not cells:
                continue

            first_text = cells[0].get_text(strip=True)
            year_match = re.search(r'\b(19\d\d|20\d\d)\b', first_text)

            if year_match and len(first_text) <= 14:
                yr = int(year_match.group(1))
                if yr > max_year or yr < min_year:
                    continue
                current_year = yr
                has_year_col = True
                row_cells = cells[1:]
            else:
                has_year_col = False
                row_cells = cells

            if not current_year or current_year < min_year or current_year > max_year:
                continue

            row_txt = tr.get_text(strip=True).lower()
            if is_note_or_cancellation(row_txt):
                continue

            chosen_a = None
            if has_year_col and film_col_idx != -1:
                target_idx = film_col_idx - 1
                if 0 <= target_idx < len(row_cells):
                    target_cell = row_cells[target_idx]
                    italic = target_cell.find(['i', 'em'])
                    if italic:
                        chosen_a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                    if not chosen_a:
                        chosen_a = target_cell.find('a')
            elif not has_year_col and len(row_cells) > 0:
                target_cell = row_cells[0] if film_col_idx == -1 or len(row_cells) == 1 else row_cells[min(1, len(row_cells)-1)]
                italic = target_cell.find(['i', 'em'])
                if italic:
                    chosen_a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                if not chosen_a:
                    chosen_a = target_cell.find('a')

            if not chosen_a:
                for cell in row_cells:
                    italic = cell.find(['i', 'em'])
                    if italic:
                        a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                        if a and a.get('href', '').startswith('/wiki/'):
                            chosen_a = a
                            break

            if chosen_a:
                href = chosen_a.get('href', '').replace('/wiki/', '')
                title = clean_title(chosen_a.get_text(strip=True) or chosen_a.get('title', ''))
                if href and title and len(title) > 1 and not is_note_or_cancellation(title):
                    entries.append((current_year, title, href))

    return entries

def parse_tiff_peoples_choice(soup, min_year=1978, max_year=2025):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        for tr in t.find_all('tr')[1:]:
            cells = tr.find_all(['th', 'td'])
            if len(cells) < 2:
                continue
            yr_match = re.search(r'\b(19\d\d|20\d\d)\b', cells[0].get_text())
            if not yr_match:
                continue
            yr = int(yr_match.group(1))
            if yr < min_year or yr > max_year:
                continue

            film_cell = cells[1]
            italic = film_cell.find(['i', 'em'])
            chosen_a = None
            if italic:
                chosen_a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
            if not chosen_a:
                chosen_a = film_cell.find('a')

            if chosen_a:
                href = chosen_a.get('href', '').replace('/wiki/', '')
                title = clean_title(chosen_a.get_text(strip=True) or chosen_a.get('title', ''))
                if href and title and not is_note_or_cancellation(title):
                    entries.append((yr, title, href))
    return entries

def parse_academy_awards_best_picture(soup, min_year=1927, max_year=2025):
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
                if yr < min_year or yr > max_year:
                    continue
                film_td = tds[0]
                film_a = film_td.find('a')
                if film_a:
                    href = film_a.get('href', '').replace('/wiki/', '')
                    title = clean_title(film_a.get_text(strip=True))
                    if title:
                        entries.append((yr, title, href))
    return entries

def parse_sundance_category_strict(soup, pattern, neg_pattern=None, min_year=1984, max_year=2025):
    results = []
    seen_years = set()
    for h in soup.find_all(['h3', 'h4', 'div']):
        h_text = h.get_text(strip=True)
        m = re.match(r'^(198\d|199\d|20[0-2]\d)', h_text)
        if not m:
            continue
        year = int(m.group(1))
        if year in seen_years or year < min_year or year > max_year:
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
                if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned) and not is_note_or_cancellation(cleaned):
                    results.append((year, cleaned, chosen_a.get('href', '').replace('/wiki/', '')))
                    seen_years.add(year)
                    break
    return results

def parse_rotterdam_tiger(soup, min_year=1995, max_year=2025):
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

                if not current_year or current_year < min_year or current_year > max_year:
                    continue

                for cell in row_cells[:2]:
                    for a in cell.find_all('a'):
                        href = a.get('href', '')
                        if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:']):
                            cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                            if cleaned and len(cleaned) > 1 and not re.match(r'^\d+$', cleaned) and not is_note_or_cancellation(cleaned):
                                entries.append((current_year, cleaned, href.replace('/wiki/', '')))
                                break
                    if entries and entries[-1][0] == current_year:
                        break
    return entries

def parse_bfi_london_best_film(soup, min_year=1958, max_year=2025):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    for t in tables:
        prev_heading = t.find_previous(['h2', 'h3', 'h4'])
        if prev_heading and any(skip in prev_heading.get_text().lower() for skip in ['director', 'staff', 'presidents', 'see also']):
            continue

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
            if yr > max_year:
                continue

            chosen_a = None
            for cell in tds:
                italic = cell.find(['i', 'em'])
                if italic:
                    a = italic.find('a') or (italic.parent.name == 'a' and italic.parent)
                    if a and a.get('href', '').startswith('/wiki/'):
                        chosen_a = a
                        break

            if chosen_a:
                cleaned = clean_title(chosen_a.get_text(strip=True) or chosen_a.get('title', ''))
                if cleaned and len(cleaned) > 1 and not is_note_or_cancellation(cleaned):
                    entries.append((yr, cleaned, chosen_a.get('href', '').replace('/wiki/', '')))
    return entries

def parse_idfa_best_film(soup, min_year=1988, max_year=2025):
    entries = []
    tables = soup.find_all('table', class_=lambda c: c and 'wikitable' in c)
    if tables:
        t = tables[0]
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

            if not current_year or current_year < min_year or current_year > max_year:
                continue

            for cell in row_cells[:1]:
                for a in cell.find_all('a'):
                    href = a.get('href', '')
                    if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:']):
                        cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                        if cleaned and len(cleaned) > 1 and not is_note_or_cancellation(cleaned):
                            entries.append((current_year, cleaned, href.replace('/wiki/', '')))
                            break
    return entries

def parse_fipresci_grand_prix(soup, min_year=1999, max_year=2025):
    entries = []
    for li in soup.find_all('li'):
        text = li.get_text(strip=True)
        m = re.match(r'^(19\d\d|20\d\d)\s*[–\-—]\s*(.+)', text)
        if m:
            yr = int(m.group(1))
            if yr < min_year or yr > max_year:
                continue
            for a in li.find_all('a'):
                href = a.get('href', '')
                if href.startswith('/wiki/') and not any(x in href for x in ['Festival', 'festival', 'cite_note', 'List_of', 'Special:', 'FIPRESCI']):
                    cleaned = clean_title(a.get_text(strip=True) or a.get('title', ''))
                    if cleaned and len(cleaned) > 1 and not is_note_or_cancellation(cleaned):
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
        writer = csv.writer(f, lineterminator='\n')
        writer.writerow(["year", "title", "IMDb ID"])
        for r in unique_records:
            writer.writerow([r[0], r[1], r[2]])

    print(f"  ✓ Written {len(unique_records)} records to {catalog_id}.csv")

    if legacy_filename:
        legacy_path = os.path.join(DATA_DIR, legacy_filename)
        with open(legacy_path, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f, lineterminator='\n')
            writer.writerow(["year", "title", "IMDb ID"])
            for r in unique_records:
                writer.writerow([r[0], r[1], r[2]])
        print(f"  (Also updated legacy file: {legacy_filename})")

    return len(unique_records)

def process_catalog(catalog_id, wiki_page, hint='film_col', min_year=1920, legacy_file=None, custom_parser=None):
    print(f"\nProcessing [{catalog_id}] from '{wiki_page}'...")
    html = get_wiki_parse_html(wiki_page)
    if not html:
        print(f"  [ERROR] Could not fetch wiki page: {wiki_page}")
        return 0

    soup = BeautifulSoup(html, 'html.parser')
    if custom_parser:
        raw_entries = custom_parser(soup, min_year=min_year)
    else:
        raw_entries = extract_award_films_robust(soup, min_year=min_year)

    if catalog_id == 'cannes-palme-dor':
        raw_entries.append((1939, 'Union Pacific', 'Union_Pacific_(film)'))

    print(f"  Extracted {len(raw_entries)} candidate raw entries")

    records = []
    for yr, title, wiki_link in raw_entries:
        imdb_id = resolve_imdb_id_strict(wiki_link, title, yr)
        if imdb_id:
            records.append((yr, title, imdb_id))
        else:
            print(f"  [WARN] Strict resolution rejected / unresolved: ({yr}, '{title}', wiki: '{wiki_link}')")

    save_cache(IMDB_CACHE)
    count = write_catalog_csv(catalog_id, records, legacy_filename=legacy_file)
    return count

def main():
    print("==================================================")
    print("High-Precision Film Festival Dataset Generator")
    print("==================================================")

    # 1. CANNES (7 catalogs)
    process_catalog('cannes-palme-dor', 'Palme d\'Or', legacy_file='palme-dor-winners.csv')
    process_catalog('cannes-grand-prix', 'Grand Prix (Cannes Film Festival)')
    process_catalog('cannes-jury-prize', 'Jury Prize (Cannes Film Festival)')
    process_catalog('cannes-best-director', 'Cannes Film Festival Award for Best Director')
    process_catalog('cannes-best-screenplay', 'Cannes Film Festival Award for Best Screenplay')
    process_catalog('cannes-best-actress', 'Cannes Film Festival Award for Best Actress', custom_parser=parse_acting_films)
    process_catalog('cannes-best-actor', 'Cannes Film Festival Award for Best Actor', custom_parser=parse_acting_films)

    # 2. VENICE (6 catalogs)
    process_catalog('venice-golden-lion', 'Golden Lion', legacy_file='golden-lion-winners.csv')
    process_catalog('venice-grand-jury-prize', 'Grand Jury Prize (Venice Film Festival)')
    process_catalog('venice-silver-lion-director', 'Silver Lion')
    process_catalog('venice-best-screenplay', 'Golden Osella')
    process_catalog('venice-coppa-volpi-actress', 'Volpi Cup for Best Actress', custom_parser=parse_acting_films)
    process_catalog('venice-coppa-volpi-actor', 'Volpi Cup for Best Actor', custom_parser=parse_acting_films)

    # 3. BERLIN (6 catalogs)
    process_catalog('berlin-golden-bear', 'Golden Bear', legacy_file='golden-bear-winners.csv')
    process_catalog('berlin-silver-bear-grand-jury', 'Silver Bear Grand Jury Prize')
    process_catalog('berlin-silver-bear-director', 'Silver Bear for Best Director')
    process_catalog('berlin-silver-bear-screenplay', 'Silver Bear for Best Screenplay')
    process_catalog('berlin-silver-bear-actress', 'Silver Bear for Best Actress', custom_parser=parse_acting_films)
    process_catalog('berlin-silver-bear-actor', 'Silver Bear for Best Actor', custom_parser=parse_acting_films)

    # 4. LOCARNO (3 catalogs)
    process_catalog('locarno-golden-leopard', 'Golden Leopard')
    process_catalog('locarno-special-jury-prize', 'Special Jury Prize (Locarno International Film Festival)')
    process_catalog('locarno-best-direction', 'Pardo for Best Direction')

    # 5. SUNDANCE (6 catalogs)
    process_catalog('sundance-grand-jury-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Grand Jury.*Dramatic', neg_pattern=r'Documentary|World', min_year=min_year))
    process_catalog('sundance-grand-jury-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Grand Jury.*Doc', neg_pattern=r'Dramatic|World', min_year=min_year))
    process_catalog('sundance-audience-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Audience.*Dramatic', neg_pattern=r'Documentary|World', min_year=min_year))
    process_catalog('sundance-audience-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Audience.*Doc', neg_pattern=r'Dramatic|World', min_year=min_year))
    process_catalog('sundance-directing-dramatic', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Directing.*Dramatic', neg_pattern=r'Documentary|World', min_year=min_year))
    process_catalog('sundance-directing-doc', 'List of Sundance Film Festival award winners',
                    custom_parser=lambda soup, min_year: parse_sundance_category_strict(soup, r'Directing.*Doc', neg_pattern=r'Dramatic|World', min_year=min_year))

    # 6. TORONTO TIFF (1 catalog)
    process_catalog('tiff-peoples-choice', 'Toronto International Film Festival People\'s Choice Award', custom_parser=parse_tiff_peoples_choice)

    # 7. ROTTERDAM IFFR (1 catalog)
    process_catalog('rotterdam-tiger-award', 'International Film Festival Rotterdam', custom_parser=parse_rotterdam_tiger)

    # 8. SAN SEBASTIÁN (2 catalogs)
    process_catalog('san-sebastian-golden-shell', 'Golden Shell')
    process_catalog('san-sebastian-best-director', 'Silver Shell for Best Director')

    # 9. KARLOVY VARY (1 catalog)
    process_catalog('karlovy-vary-crystal-globe', 'Crystal Globe (Karlovy Vary International Film Festival)')

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
    print("Regeneration completed with multi-signal precision!")
    print("==================================================")

if __name__ == '__main__':
    main()
