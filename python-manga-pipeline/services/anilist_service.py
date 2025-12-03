import requests
from config import settings
from typing import Any, Dict, List, Tuple
from utils import utils

class AniListClient:
    TOP_100_MANGA_BY_POPULARITY_QUERY = """
    query ($type: MediaType, $format: MediaFormat, $sort: [MediaSort], $staffSort: [StaffSort], $page: Int, $staffPerPage: Int, $asHtml: Boolean, $countryOfOrigin: CountryCode) {
        Page(page: $page) {
            pageInfo {
                perPage
                currentPage
                lastPage
                hasNextPage
            }
            media (type: $type, format: $format, sort: $sort, countryOfOrigin: $countryOfOrigin) {
                title {
                    romaji
                    english
                }
                description(asHtml: $asHtml)
                genres
                status
                chapters
                staff(sort: $staffSort perPage: $staffPerPage) {
                    edges {
                        node {
                            name {
                                full
                            }
                        }
                        role
                    }
                }
            }
        }
    } 
    """

    def __init__(self):
        self.session = requests.Session()

    def _fetch (self, query: str, variables: Dict[str, Any]) -> Dict[str, Any]:
        response = self.session.post(
            settings.ANILIST_BASE_URL,
            json={"query": query, "variables": variables},
            timeout=10
        )

        response.raise_for_status()
        return response.json()

    def get_100_manga_by_popularity(self) -> Tuple[List[Dict[str, Any]], List[Dict[str, Any]]]:
        all_manga = []
        all_not_completed_manga_title = []

        payload = {
            "type": "MANGA",
            "format": "MANGA",
            "sort": "POPULARITY_DESC",
            "staffSort2": "RELEVANCE",
            "staffPerPage": 1,
            "asHtml": False,
            "countryOfOrigin": "JP",
        }

        # 1 page has 50 records from anilist
        for page in [1, 2]:
            payload["page"] = page
            response_data = self._fetch(self.TOP_100_MANGA_BY_POPULARITY_QUERY, payload)
            page_manga = response_data["data"]["Page"]["media"]

            # Clean data
            for item in page_manga:
                # Clean Description
                desc = utils.first_paragraph(item.get("description", ""))
                item["description"] = desc

                # Clean Status & Get list of non-COMPLETED
                match item["status"]:
                    case "RELEASING":
                        item["status"] = "ONGOING"
                        all_not_completed_manga_title.append(item["title"])

                    case "HIATUS":
                        all_not_completed_manga_title.append(item["title"])
                        pass
                    
                    case "FINISHED":
                        item["status"] = "COMPLETED"                    
                    
                    case "CANCELLED":
                        item["status"] = "COMPLETED"
                    
                    case _:
                        item["status"] = "ONGOING"
                        item["chapters"] = 0   

            all_manga.extend(page_manga)

        return all_manga, all_not_completed_manga_title
    