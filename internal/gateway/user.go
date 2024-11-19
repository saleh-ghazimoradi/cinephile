package gateway

import (
	"errors"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"net/http"
)

type userHandler struct {
	userService service.User
	mailService service.Mail
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

	background(func() {
		if err := u.mailService.Send(user.Email, "user_welcome.tmpl", user); err != nil {
			logger.Logger.Error(err.Error())
		}
	})

	if err := writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil); err != nil {
		serverErrorResponse(w, r, err)
		return
	}
}

func NewUserHandler(userService service.User, mailService service.Mail) *userHandler {
	return &userHandler{
		userService: userService,
		mailService: mailService,
	}
}
