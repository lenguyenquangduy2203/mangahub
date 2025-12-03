from dataclasses import dataclass

@dataclass
class Manga:
    id: str #one-piece
    title: str
    author: str
    genres: list
    status: str
    total_chapters: int
    desciption: str
    mangaupdates_id: int
