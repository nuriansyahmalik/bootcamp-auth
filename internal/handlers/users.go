package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/users"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type UsersHandler struct {
	UsersService   users.UsersService
	AuthMiddleware *middleware.Authentication
}

func ProvideUsersBarBazHandler(usersService users.UsersService, authMiddleware *middleware.Authentication) UsersHandler {
	return UsersHandler{
		UsersService:   usersService,
		AuthMiddleware: authMiddleware,
	}
}

func (h *UsersHandler) Router(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Register)
		r.Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(jwt.AuthMiddleware)
			r.Get("/profile", h.Profile)
			r.Put("/profile", h.UpdateProfile)
			r.Get("/validate-auth", h.ValidateUsers)
		})
	})
}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat users.UserRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.NewV4()
	if err != nil {
		logger.ErrorWithStack(err)
	}
	foo, err := h.UsersService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, foo)
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat users.LoginRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}

	foo, err := h.UsersService.Login(requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}
func (h *UsersHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	user, err := h.UsersService.ResolveByID(claims.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, failure.Unauthorized("Unauthorized"))
		return
	}

	response.WithJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat users.UserRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		logger.ErrorWithStack(err)
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	user, err := h.UsersService.Update(claims.ID, requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) ValidateUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response.WithJSON(w, http.StatusOK, claims)
}
