import json
import sqlite3
from config.settings import DB_PATH
from models.manga import Manga
from typing import Any, Dict, List

import logging
logger = logging.getLogger(__name__)

class DbService():
    def __init__(self):
        pass

    def _get_connection(self):
        import os
        os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
        return sqlite3.connect(DB_PATH)

    def store_manga(self, mangas: List[Manga]):
        with self._get_connection() as conn:
            cursor = conn.cursor()

            for m in mangas:
                cursor.execute("""
                    INSERT OR REPLACE INTO manga
                    (id, title, author, genres, status, total_chapters, description, mangaupdates_id)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    m.id,
                    m.title,
                    m.author,
                    m.genres,     # serialize list → JSON string
                    m.status,
                    m.total_chapters,
                    m.description,
                    m.mangaupdates_id,
                ))

        logger.info(f"Stored {len(mangas)} manga records to DB: {DB_PATH}")
    
    def update_manga_chapters(self, update_map: Dict[int, Dict[str, Any]]):
        with self._get_connection() as conn:
            cursor = conn.cursor()

            updated_count = 0
            for mu_id, item in update_map.items():
                latest = item.get("latest")
                is_completed = item.get("isCompleted")
                
                if is_completed:
                    status = "COMPLETED"
                else: status = None

                cursor.execute("""
                    UPDATE manga
                    SET total_chapters = ?,
                        status = COALESCE(?, status)
                    WHERE mangaupdates_id = ?
                """, (latest, status, mu_id))

                updated_count += cursor.rowcount

        logger.info(f"Updated latest_chapter for {updated_count} manga series in DB.")

    def fetch_ongoing_mangas(self) -> List[int]:
        with self._get_connection() as conn:
            cursor = conn.cursor()
            
            cursor.execute("""
                SELECT mangaupdates_id 
                FROM manga
                WHERE mangaupdates_id IS NOT NULL 
                AND (status = 'ONGOING' OR status = 'HIATUS')
            """)

            rows = cursor.fetchall()

        return [row[0] for row in rows]