import json
import logging
import os
import time
from config.settings import OUTPUT_PATH
from models.manga import Manga
from services.anilist_service import AniListClient
from services.mangaupdate_service import MangaUpdateService
from typing import List
from utils import utils

logger = logging.getLogger(__name__)

class PipelineRunner:
    def __init__(self):
        pass

    def _store_manga_list(self, mangas: List[Manga]):
        data = [m.__dict__ for m in mangas]
        os.makedirs(os.path.dirname(OUTPUT_PATH), exist_ok=True)
        with open(OUTPUT_PATH, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        logger.info(f"Wrote {len(mangas)} records to {OUTPUT_PATH}")

    def run_top_100_init_pipline(self):
        start_time = time.time()

        ani = AniListClient()
        mu = MangaUpdateService()
        
        logger.info("(1) Fetching metadata from AniList...")
        top_100, non_completed_titles_obj = ani.get_100_manga_by_popularity()

        logger.info("(2) Build list of titles...")
        non_completed_titles = []
        for t_obj in non_completed_titles_obj:
            title = utils.title_to_str(t_obj)
            if title:
                non_completed_titles.append(title)

        logger.info("(3) Fetching latest chapter of ONGOING or HIATUS series from MangaUpdates...")
        mu_map = mu.get_lastest_manga_updates(non_completed_titles)

        logger.info("(4) Enrinch top_100 items with MU info where available...")
        for item in top_100:
            title = utils.title_to_str(item.get("title"))
            mu_info = mu_map.get(title)

            if mu_info:
                item["chapters"] = mu_info.get("latest_chapter")
                item["mangaupdates_id"] = mu_info.get("id")

        logger.info("(5) Map Dict -> Manga class...")
        manga_models = []
        for item in top_100:
            try:
                manga = utils.map_to_manga_model(item)
                manga_models.append(manga)
            except Exception as e:
                logger.warn(f"Failed to map item {item.get('title').get('english')}: {e}")

        logger.info("(6) Store models...")
        self._store_manga_list(manga_models)
        end_time = time.time()

        logger.info(f"Pipeline finished ({(end_time - start_time):.2f}s): stored {len(manga_models)} manga records.")    