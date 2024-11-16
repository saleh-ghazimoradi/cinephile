package service_models

import (
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
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

type UpdateMoviePayload struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *int32   `json:"runtime"`
	Genres  []string `json:"genres"`
}

type WithPagination struct {
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
	Filter `json:"filter"`
}

type Filter struct {
	Page     int    `json:"page" validate:"gte=0,lte=10000000"`
	PageSize int    `json:"page_size" validate:"gte=1,lte=20"`
	Sort     string `json:"sort" validate:"oneof=asc desc"`
}

type MovieWithMetaData struct {
	Movie
	MovieCount int `json:"movie_count"`
}
