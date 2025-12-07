import logging
import time
from services.anilist_service import AniListClient
from services.db_service import DbService
from services.mangaupdate_service import MangaUpdateService
from utils import utils

logger = logging.getLogger(__name__)

class PipelineRunner:
    def __init__(self):
        pass
    
    def run_top_100_init_pipline(self):
        start_time = time.time()

        ani = AniListClient()
        mu = MangaUpdateService()
        db = DbService()
        
        logger.info("(1) Fetching metadata from AniList...")
        top_100, non_completed_titles_obj = ani.get_100_manga_by_popularity()

        logger.info("(2) Fetching latest chapter of ONGOING or HIATUS series from MangaUpdates...")
        mu_map = mu.get_latest_manga_updates_init(non_completed_titles_obj)

        logger.info("(3) Enrinch top_100 items with MU info where available...")
        for item in top_100:
            title = utils.title_to_str(item.get("title"))
            mu_info = mu_map.get(title)

            if mu_info:
                item["chapters"] = mu_info.get("latest_chapter")
                item["mangaupdates_id"] = mu_info.get("id")

        logger.info("(4) Map Dict -> Manga class...")
        manga_models = []
        for item in top_100:
            try:
                manga = utils.map_to_manga_model(item)
                manga_models.append(manga)
            except Exception as e:
                logger.warning(f"Failed to map item {item.get('title').get('english')}: {e}")

        logger.info("(5) Store models...")
        db.store_manga(manga_models)
        
        end_time = time.time()

        logger.info(f"Pipeline finished ({(end_time - start_time):.2f}s): stored {len(manga_models)} manga records.")    

    def run_update_non_complete_pipeline(self):
        start_time = time.time()

        mu = MangaUpdateService()
        db = DbService() 

        logger.info("(1) Fetch info from db...")
        rows = db.fetch_ongoing_mangas()

        logger.info(f"(2) Update {len(rows)} id from MangaUpdates...")
        update_map = mu.get_latest_manga_updates_update(rows)

        logger.info("(3) Store the updated update...")
        db.update_manga_chapters(update_map)

        end_time = time.time()

        logger.info(f"Pipeline finished ({(end_time - start_time):.2f}s): update {len(update_map)} manga series.")