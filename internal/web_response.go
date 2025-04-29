package internal

import "math"

type WebResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type PaginationResponse struct {
	Items      any            `json:"items,omitempty"`
	Pagination PaginationMeta `json:"pagination"`
}
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

func NewPaginationMeta(currentPage, perPage, totalItems int) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	return PaginationMeta{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
}
