package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dionaditya/go-production-ready-api/internal/models"
	"github.com/gorilla/mux"
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

func (h *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.ParseUint(id, 10, 64)

	w.Header().Set("Content-Type", "application/json; charse=UTF-8")

	type newData struct {
		Username string
	}

	var userData newData

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
	}

	user, err := h.Service.UpdateUser(uint(userID), userData)

	if err != nil {
		sendErrorResponse(w, "failed to update comments with ID"+id, err)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		panic(err)
	}
}
