import os

ANILIST_BASE_URL = "https://graphql.anilist.co"
MANGAUPDATE_BASE_URL = "https://api.mangaupdates.com/v1"

SCHEDULE_HOURS = 6
OUTPUT_PATH = os.path.join(os.path.dirname(__file__), "..", "output", "top_manga.json")