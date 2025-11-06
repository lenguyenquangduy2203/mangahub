package dtos

type MangaDetail struct {
	MangaID       string   `json:"manga_id"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"` // enum: [ongoing, completed, hiatus]
	TotalChapters int      `json:"total_chapters"`
	Description   string   `json:"description"`
	CoverURL      string   `json:"cover_url"`
}

type MangaListItem struct {
	MangaID       string `json:"manga_id"`
	Title         string `json:"title"`
	TotalChapters int    `json:"total_chapters"`
	Status        string `json:"status"` // enum: [ongoing, completed, hiatus]
}

type PaginatedMangaList struct {
	Total   int             `json:"total"`   // Total number of matching manga
	Limit   int             `json:"limit"`   // Number of items per page
	Offset  int             `json:"offset"`  // Offset used for pagination
	Results []MangaListItem `json:"results"` // The list of manga items
}
