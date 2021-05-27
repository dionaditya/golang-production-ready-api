package transportHTTP

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dionaditya/go-production-ready-api/internal/comment"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

type Response struct {
	Message string
	Error   string
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charse=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "I am alive"}); err != nil {
			panic(err)
		}

	})
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)

	w.Header().Set("Content-Type", "application/json; charse=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		sendErrorResponse(w, "Failed to parse ID", err)
	}

	comment, err := h.Service.GetComment(uint(i))

	if err != nil {
		sendErrorResponse(w, "Failed to retrieve comment by ID", err)
	}

	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		panic(err)
	}
}

func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetComments()
	if err != nil {
		sendErrorResponse(w, "failed to retrieve all comments", err)
	}

	w.Header().Set("Content-Type", "application/json; charse=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		panic(err)
	}

}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charse=UTF-8")
	w.WriteHeader(http.StatusOK)

	var comment comment.Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
	}

	comment, err := h.Service.PostCmoment(comment)

	if err != nil {
		sendErrorResponse(w, "failed to post new comments", err)
	}

	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)

	w.Header().Set("Content-Type", "application/json; charse=UTF-8")
	w.WriteHeader(http.StatusOK)

	var comment comment.Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "failed to decode json", err)
	}

	comment, err = h.Service.UpdateComemnt(uint(commentID), comment)

	if err != nil {
		sendErrorResponse(w, "failed to update comments with ID"+id, err)
	}

	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		sendErrorResponse(w, "failed to parse ID", err)
	}

	w.Header().Set("Content-Type", "application/json; charse=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = h.Service.DeleteComment(uint(commentID))

	if err != nil {
		sendErrorResponse(w, "failed to delete comment with ID"+id, err)

	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Comment successfully deleted"}); err != nil {
		panic(err)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(Response{Message: "Comment successfully deleted", Error: err.Error()}); err != nil {
		panic(err)
	}
}
