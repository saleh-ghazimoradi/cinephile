package gateway

import (
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"net/http"
)

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Logger.Error(err.Error(), "method", method, "uri", uri)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err)
}

func editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	errorResponse(w, r, http.StatusConflict, message)
}
