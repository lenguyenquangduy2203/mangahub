import requests
from config import settings
from typing import Any, Dict, List

class MangaUpdateService:
    def __init__(self):
        self.session = requests.Session()

    def _fetch (self, subUrl: str, payload: Dict[str, Any] | None = None, method: str = "POST") -> Dict[str, Any]:
        url = settings.MANGAUPDATE_BASE_URL + subUrl
        if method == "POST":
            response = self.session.post(
                url,
                json=payload,
                timeout=10
            )
        else:
            response = self.session.get(
                url,
                timeout=10
            )

        response.raise_for_status()
        return response.json()
    
    def _search_series_id(self, title: str) -> int | None:
        payload = {
            "search": title,
            "type": ["Manga"],
            "page": 1
        }
        data = self._fetch("/series/search", payload)

        results = data.get("results", [])
        if not results:
            return None
        
        # Get the most relevant series
        return results[0]["record"]["series_id"]
        

    def _search_series_latest_chapter(self, series_id: int) -> int | None:
        subUrl = f"/series/{series_id}"
        data = self._fetch(subUrl, method="GET")
        latest_chapter = data.get("latest_chapter")
        return latest_chapter

    def get_lastest_manga_updates(self, titles: List[str]) -> Dict[str, Any]:
        results = {}

        for title in titles:
            series_id = self._search_series_id(title)
            if not series_id:
                results[title] = {
                    "id": None,
                    "latest_chapter": None
                }
                continue

            latest = self._search_series_latest_chapter(series_id)

            results[title] = {
                "id": series_id,
                "latest_chapter": latest
            }

        return results
    