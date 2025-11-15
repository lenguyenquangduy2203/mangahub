package dtos

import "time"

type UserLibrary struct {
	ReadingList ReadingList `json:"reading_lists"`
}

type ReadingList struct {
	Reading    []UserLibraryItem `json:"reading"`
	Completed  []UserLibraryItem `json:"completed"`
	PlanToRead []UserLibraryItem `json:"plan_to_read"`
}

type UserLibraryItem struct {
	MangaID        string    `json:"manga_id"`
	CurrentChapter int       `json:"current_chapter"`
	Status         string    `json:"status"`
	LastUpdated    time.Time `json:"last_updated"`
}

type PaginatedUserLibrary struct {
	Total   int               `json:"total"`   // Total number of user's manga with that status
	Limit   int               `json:"limit"`   // Number of items per page
	Offset  int               `json:"offset"`  // Offset used for pagination
	Results []UserLibraryItem `json:"results"` // The list of user's manga items
}
