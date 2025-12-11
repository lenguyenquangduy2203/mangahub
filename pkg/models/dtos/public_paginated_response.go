package dtos

import "mangahub/pkg/pagination"

type PaginationResponse struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// Public structures for paginated responses
type PaginatedResponse[T any] struct {
	Results    []T                `json:"results"`
	Pagination PaginationResponse `json:"pagination"`
}

func BuildPaginatedResponse[T any](data pagination.Paginated) PaginatedResponse[T] {
	limit := data.GetLimit()
	offset := data.GetOffset()
	total := data.GetTotal()

	page := (offset / limit) + 1

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return PaginatedResponse[T]{
		Results: data.GetResults().([]T),
		Pagination: PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
}
