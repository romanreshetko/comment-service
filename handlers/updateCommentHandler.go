package handlers

import (
	"comment-service/models"
	"comment-service/repository"
	"net/http"
	"slices"
	"strconv"
)

func (h *Handler) UpdateCommentStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	validStatuses := []string{"published", "moderating", "blocked", "reported", "blocked_reported", "undefined", "moderation_error"}
	if !slices.Contains(validStatuses, status) {
		http.Error(w, "incorrect status", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect comment_id", http.StatusBadRequest)
		return
	}
	if claims.Role == "user" && status != "reported" && status != "blocked_reported" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if claims.Role != "user" && claims.Role != "moderator" && claims.Role != "admin" && claims.Role != "service" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = repository.UpdateCommentStatus(h.db, commentID, status)
	if err != nil {
		if err.Error() == "comment not found" {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
