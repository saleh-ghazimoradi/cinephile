package gateway

import (
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"net/http"
	"time"
)

type moviePayload struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

type movieHandler struct {
	movieService service.Movie
}

func (m *movieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload moviePayload
	if err := readJSON(w, r, payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%v\n", payload)
}

func (m *movieHandler) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	movie := service_models.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Black List",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	if err = writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}

func NewMovieHandler(movieService service.Movie) *movieHandler {
	return &movieHandler{movieService: movieService}
}
