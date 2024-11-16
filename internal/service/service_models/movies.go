package service_models

import (
	"net/http"
	"strconv"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MoviePayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Year    int32    `json:"year" validate:"required"`
	Runtime int32    `json:"runtime" validate:"required"`
	Genres  []string `json:"genres" validate:"required,gte=1"`
}

type UpdateMoviePayload struct {
	Title   *string  `json:"title" validate:"required,max=100"`
	Year    *int32   `json:"year" validate:"required"`
	Runtime *int32   `json:"runtime" validate:"required"`
	Genres  []string `json:"genres" validate:"required,gte=2"`
}

type Filter struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq Filter) Parse(r *http.Request) (Filter, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}
