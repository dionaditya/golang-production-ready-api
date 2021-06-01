package http

import (
	"encoding/json"
	"net/http"

	"github.com/dionaditya/go-production-ready-api/internal/models"
)

func (h *UserController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charse=UTF-8")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
		return
	}

	payload, err := h.Service.Register(user)

	if err != nil {
		sendErrorResponse(w, "failed to register", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&payload); err != nil {
		panic(err)
	}
}

func (h *UserController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charse=UTF-8")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
		return
	}

	payload, err := h.Service.Login(user.Email, user.Password)

	if err != nil {
		sendErrorResponse(w, "failed to login", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&payload); err != nil {
		panic(err)
	}
}

func (h *UserController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charse=UTF-8")

	type RefreshToken struct {
		RefreshToken string
	}
	var refreshToken RefreshToken
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
		return
	}

	payload, err := h.Service.RefreshToken(refreshToken.RefreshToken)

	if err != nil {
		sendErrorResponse(w, "failed to refresh token", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&payload); err != nil {
		panic(err)
	}
}
