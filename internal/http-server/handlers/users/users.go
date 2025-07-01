package usershandler

import (
	"context"
	"log/slog"
	"net/http"

	ssov1 "github.com/Swiffy12/protos/gen/go/sso"
	"github.com/Swiffy12/taskify/internal/clients/sso/grpc"
	"github.com/Swiffy12/taskify/internal/http-server/models"
	resp "github.com/Swiffy12/taskify/internal/lib/api/response"
	"github.com/Swiffy12/taskify/internal/lib/logger/sl"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	log *slog.Logger
	sso *grpc.Client
}

type UserService interface {
	Register(email, password string) ssov1.RegisterResponse
}

func New(log *slog.Logger, sso *grpc.Client) *UserHandler {
	return &UserHandler{log: log, sso: sso}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.Register"

	log := u.log.With(
		slog.String("op", op),
	)

	var user models.User
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to decode request body"))
		return
	}

	if err := validator.New().Struct(user); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, resp.ValidationError(validateErr))
		return
	}

	log.Info("register user", slog.String("email", user.Email))
	userId, err := u.sso.Register(context.Background(), user.Email, user.Password)
	if err != nil {
		log.Error("internal error", sl.Err(err))
		render.JSON(w, r, resp.Error("internal error"))
		return
	}
	log.Info("registered user", slog.String("email", user.Email))

	render.JSON(w, r, resp.OK(userId))
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.Login"

	log := u.log.With(
		slog.String("op", op),
	)

	var user models.User
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to decode request body"))
		return
	}

	if err := validator.New().Struct(user); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, resp.ValidationError(validateErr))
		return
	}

	log.Info("login user", slog.String("email", user.Email))
	token, err := u.sso.Login(context.Background(), user.Email, user.Password, u.sso.AppId)
	if err != nil {
		log.Error("internal error", sl.Err(err))
		render.JSON(w, r, resp.Error("internal error"))
		return
	}
	log.Info("user is logged in", slog.String("email", user.Email))

	render.JSON(w, r, resp.OK(token))
}
