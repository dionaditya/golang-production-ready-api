package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dionaditya/go-production-ready-api/internal/comment"
	"github.com/gorilla/mux"
)

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
	if err := json.NewEncoder(w).Encode(&ResponseAPI{
		Result: &comments,
	}); err != nil {
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