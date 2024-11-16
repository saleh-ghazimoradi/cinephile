package gateway

import (
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"net/http"
)

type movieHandler struct {
	movieService service.Movie
}

func (m *movieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_models.MoviePayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	movie := &service_models.Movie{
		Title:   payload.Title,
		Year:    payload.Year,
		Runtime: payload.Runtime,
		Genres:  payload.Genres,
	}

	if err := m.movieService.Create(r.Context(), movie); err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err := writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		serverErrorResponse(w, r, err)
		return
	}
}

func (m *movieHandler) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	movie, err := m.movieService.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	if err = writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *movieHandler) ListMoviesHandler(w http.ResponseWriter, r *http.Request) {
	p := service_models.Filter{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := p.Parse(r)

	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	movies, err := m.movieService.GetAll(r.Context(), fq)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *movieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_models.UpdateMoviePayload
	id, err := readIDParam(r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	movie, err := m.movieService.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	if err = readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		movie.Title = *payload.Title
	}

	if payload.Year != nil {
		movie.Year = *payload.Year
	}

	if payload.Runtime != nil {
		movie.Runtime = *payload.Runtime
	}

	if payload.Genres != nil {
		movie.Genres = payload.Genres
	}

	if err = m.movieService.Update(r.Context(), movie); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			editConflictResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *movieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	if err := m.movieService.Delete(r.Context(), id); err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}

func NewMovieHandler(movieService service.Movie) *movieHandler {
	return &movieHandler{movieService: movieService}
}
