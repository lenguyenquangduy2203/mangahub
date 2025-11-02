package dtos

// MangaSearchQuery represents the parameters accepted for searching manga.
type MangaSearchQuery struct {
	Title  string `form:"title"`
	Author string `form:"author"`
	Genre  string `form:"genre"`
	Status string `form:"status"` // enum: [ongoing, completed, hiatus]
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type MangaIDQuery struct {
	ID string `uri:"id" binding:"required"`
}
