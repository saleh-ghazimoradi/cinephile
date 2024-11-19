package gateway

import (
	"errors"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"net/http"
	"time"
)

type userHandler struct {
	userService  service.User
	mailService  service.Mail
	tokenService service.Token
}

func (u *userHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload service_models.UserPayload
	if err := readJSON(w, r, &userPayload); err != nil {
		badRequestResponse(w, r, err)
	}

	user := &service_models.User{
		Name:      userPayload.Name,
		Email:     userPayload.Email,
		Activated: false,
	}

	if err := user.Password.Set(userPayload.Password); err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(userPayload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := u.userService.Insert(r.Context(), user); err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			conflictResponse(w, r, err)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := u.tokenService.New(user.ID, 3*24*time.Hour, service_models.ScopeActivation)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	background(func() {

		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		if err := u.mailService.Send(user.Email, "user_welcome.tmpl", data); err != nil {
			logger.Logger.Error(err.Error())
		}
	})

	if err := writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil); err != nil {
		serverErrorResponse(w, r, err)
		return
	}
}

func (u *userHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var tokenPayload service_models.TokenPayload
	if err := readJSON(w, r, &tokenPayload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(tokenPayload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := u.userService.GetForToken(r.Context(), service_models.ScopeActivation, tokenPayload.Plaintext)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			logger.Logger.Error(err.Error())
			badRequestResponse(w, r, err)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	if err := u.userService.Update(r.Context(), user); err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			editConflictResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = u.tokenService.DeleteAllForUser(service_models.ScopeActivation, user.ID)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func NewUserHandler(userService service.User, mailService service.Mail, tokenService service.Token) *userHandler {
	return &userHandler{
		userService:  userService,
		mailService:  mailService,
		tokenService: tokenService,
	}
}
