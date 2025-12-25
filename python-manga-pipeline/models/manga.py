from dataclasses import dataclass
from typing import Optional

@dataclass
class Manga:
    id: str #one-piece
    title: str
    author: str
    genres: list
    status: str
    total_chapters: int
    description: str
    mangaupdates_id: Optional[int]
