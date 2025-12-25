import os

from dotenv import load_dotenv

load_dotenv(os.path.join(os.path.dirname(__file__), "..", "..", ".env"))

ANILIST_BASE_URL = "https://graphql.anilist.co"
MANGAUPDATE_BASE_URL = "https://api.mangaupdates.com/v1"

TIME_OUT = 60  # In seconds

SCHEDULE_HOURS = 6
# OUTPUT_PATH = os.path.join(os.path.dirname(__file__), "..", "output")

# INIT_FILE = "top_manga.json"
# UPDATE_FILE = "update_manga.json"

# INIT_PATH = os.path.join(OUTPUT_PATH, INIT_FILE)
# UPDATE_PATH = os.path.join(OUTPUT_PATH, UPDATE_FILE)

DB_PATH = os.getenv("PIPELINE_DB_PATH", "./data/mangahub.db")
