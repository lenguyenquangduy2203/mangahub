import requests
from config.settings import MANGAUPDATE_BASE_URL, TIME_OUT
from typing import Any, Dict, List, Tuple

class MangaUpdateService:
    def __init__(self):
        self.session = requests.Session()

    def _fetch (self, subUrl: str, payload: Dict[str, Any] | None = None, method: str = "POST") -> Dict[str, Any]:
        url = MANGAUPDATE_BASE_URL + subUrl
        if method == "POST":
            response = self.session.post(
                url,
                json=payload,
                timeout=TIME_OUT
            )
        else:
            response = self.session.get(
                url,
                timeout=TIME_OUT
            )

        response.raise_for_status()
        return response.json()
    
    def _search_series_id(self, title: str, year: str) -> int | None:
        payload = {
            "search": title,
            "type": ["Manga"],
            "year": year,
            "page": 1
        }
        data = self._fetch("/series/search", payload)

        results = data.get("results", [])
        if not results:
            return None
        
        # Get the most relevant series
        return results[0]["record"]["series_id"]
        

    def _search_series_latest_chapter(self, series_id: int) -> Tuple[int, bool]:
        subUrl = f"/series/{series_id}"
        data = self._fetch(subUrl, method="GET")
        latest_chapter = data.get("latest_chapter")
        
        if (not data.get("completed")):
            return (latest_chapter, False)
        
        return (latest_chapter, True)


    def get_latest_manga_updates_init(self, non_completed: List[Dict[str, Any]]) -> Dict[str, Any]:
        results = {}

        for obj in non_completed:
            series_id = self._search_series_id(obj["title"], str(obj["year"]))
            if not series_id:
                results[obj["title"]] = {
                    "id": None,
                    "latest_chapter": None
                }
                continue

            latest, _ = self._search_series_latest_chapter(series_id)
            
            results[obj["title"]] = {
                    "id": series_id,
                    "latest_chapter": latest
                }
                

        return results

    def get_latest_manga_updates_update(self, series_ids: List[int]) -> Dict[int, Dict[str, Any]]:
        result = {}

        for sid in series_ids:
            latest, isComplete = self._search_series_latest_chapter(sid)
            result[sid] = {
                "latest": latest,
                "isCompleted": isComplete
                }
        
        return result