package dtos

type UserMangaAddRequest struct {
	MangaID        string `json:"manga_id" binding:"required"`
	CurrentChapter int    `json:"current_chapter" binding:"min=0"`
}

type UserMangaUpdateRequest struct {
	MangaID        string `json:"manga_id" binding:"required"`
	CurrentChapter int    `json:"current_chapter" binding:"required,min=1"`
}

type UserMangaGetRequest struct {
	Status string `form:"status"` // enum: [reading, completed, plan_to_read]
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}
