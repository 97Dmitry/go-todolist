package handler

import (
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/request"
	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/response"
	"github.com/97Dmitry/go-todolist/internal/platform/logger"
	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
)

type CreateUserRequest struct {
	Email string `json:"email" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

type CreateUserResponse struct {
	UserID  string `json:"user_id"`
	Version int64  `json:"version"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responder := response.NewResponder(logger.FromContext(ctx), w)

	var createUserRequest CreateUserRequest
	if err := request.DecodeAndValidate(r, &createUserRequest); err != nil {
		responder.Error(err, "decode create user request")
		return
	}

	user, err := h.service.CreateUser(ctx, createUserRequest.Email, createUserRequest.Name)
	if err != nil {
		responder.Error(err, "create user")
		return
	}

	responder.JSON(http.StatusCreated, createUserResponseFromDomain(user))
}

func createUserResponseFromDomain(user domain.User) CreateUserResponse {
	return CreateUserResponse{
		UserID:  user.ID.String(),
		Version: user.Version,
		Email:   user.Email,
		Name:    user.Name,
	}
}
