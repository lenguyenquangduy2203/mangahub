import re
from html import unescape
from models.manga import Manga
from typing import Any, Dict, List

def generate_custom_id(title_romaji: str) -> str | None:
    if not title_romaji:
        return None
    
    replacements = {
        '×': 'x',
        '&': 'and',
        '+': 'plus',
        ':': '-',
        "'": '',
        '"': '',
        '!': '',
        '?': '',
        '.': '',
        '…': '',
    }
    for char, replacement in replacements.items():
        title_romaji = title_romaji.replace(char, replacement)
        
    title_slug = title_romaji.lower()

    title_slug = title_slug.replace(' ', '-')
    
    # Remove most special characters and punctuation
    # This regex keeps only alphanumeric characters (a-z, 0-9) and underscores (-)
    title_slug = re.sub(r'[^a-z0-9-]', '', title_slug)
    
    # Collapse multiple underscores into a single underscore
    title_slug = re.sub(r'-{2,}', '-', title_slug)
    
    title_slug = title_slug.strip('-')

    return title_slug

def first_paragraph(desc: str) -> str:
    if not desc:
        return desc
    first = re.split(r'(?i)<br\s*/?>|\r?\n', desc, maxsplit=1)[0]
    clean = re.sub(r'<[^>]+>', '', first)
    return unescape(clean).strip()

def title_to_str(title_obj: Dict[str, str], type: str = "romaji") -> str:
    if not title_obj:
        return ""
    return title_obj.get(type)

def _get_author_from_staffs(staffs: Dict[str, Any]) -> str:
    staff: List[Dict[str, Dict[str, Dict[str, str]]]] = staffs.get("edges")
    return staff[0].get("node").get("name").get("full", "UNKNOWN")

def map_to_manga_model(item: Dict[str, Any]) -> Manga:
    title_obj: Dict[str, str] = item.get("title")
    
    m_id = generate_custom_id(title_obj.get("romaji"))
    title = title_obj.get("english") or title_obj.get("romaji")

    author = _get_author_from_staffs(item.get("staff"))
    genres = ", ".join(item.get("genres") or [])
    status = item.get("status")
    total_chapters: int = item.get("chapters")
    description = item.get("description")

    mangaupdates_id: int = item.get("mangaupdates_id")

    return Manga(
        id=m_id,
        title=title,
        author=author,
        genres=genres,
        status=status,
        total_chapters=total_chapters,
        description=description,
        mangaupdates_id=mangaupdates_id
    )
    