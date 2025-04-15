package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	PerPage int    `json:"per_page" validate:"gte=1"`
	Page    int    `json:"page" validate:"gte=1"`
	Sort    string `json:"sort" validate:"oneof=ASC DESC"`
}

func (paginatedQuery *PaginatedFeedQuery) Parse(r *http.Request) (*PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	per_page := qs.Get("per_page")
	if per_page != "" {
		l, err := strconv.Atoi(per_page)
		if err != nil {
			return nil, err
		}
		paginatedQuery.PerPage = l
	}
	page := qs.Get("page")
	if page != "" {
		l, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		paginatedQuery.Page = l
	}
	sort := qs.Get("sort")
	if sort != "" {
		paginatedQuery.Sort = sort
	}
	return paginatedQuery, nil
}
