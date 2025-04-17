package model

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=4,max=20"`
	Content string   `json:"content" validate:"required,min=8,max=100"`
	Tags    []string `json:"tags" validate:"required,notempty,dive,required,min=4,max=20"`
}
type UpdatePostPayload struct {
	Id      uuid.UUID `json:"id" validate:"required"`
	Title   string    `json:"title" validate:"required,min=4,max=20"`
	Content string    `json:"content" validate:"required,min=8,max=100"`
	Tags    []string  `json:"tags" validate:"required,notempty,dive,required,min=4,max=20"`
}

type PostParams struct {
	PerPage int      `json:"per_page" validate:"gte=1"`
	Page    int      `json:"page" validate:"gte=1"`
	Sort    string   `json:"sort" validate:"oneof=ASC DESC"`
	Tags    []string `json:"tags" validate:"dive,required,min=4,max=20"`
	Search  string   `json:"search" validate:"omitempty,min=1"`
	Since   string   `json:"since" validate:"omitempty,datetime=2006-01-02"`
	Until   string   `json:"until" validate:"omitempty,datetime=2006-01-02"`
}

func ParseTime(str string) string {
	t, err := time.Parse(time.DateTime, str)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
func (params *PostParams) Parse(r *http.Request) (*PostParams, error) {
	qs := r.URL.Query()

	per_page := qs.Get("per_page")
	if per_page != "" {
		l, err := strconv.Atoi(per_page)
		if err != nil {
			return nil, err
		}
		params.PerPage = l
	}
	page := qs.Get("page")
	if page != "" {
		l, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		params.Page = l
	}
	sort := qs.Get("sort")
	if sort != "" {
		params.Sort = sort
	}
	tags := qs.Get("tags")
	if tags != "" {
		params.Tags = strings.Split(tags, ",")
	}
	search := qs.Get("search")
	if search != "" {
		params.Search = search
	}
	since := qs.Get("since")
	if since != "" {
		// params.Since = ParseTime(since)
		params.Since = since
	}
	until := qs.Get("until")
	if until != "" {
		// params.Until = ParseTime(until)
		params.Until = until
	}
	return params, nil
}
