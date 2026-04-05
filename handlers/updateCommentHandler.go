package handlers

import (
	"comment-service/models"
	"comment-service/repository"
	"encoding/json"
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

	commentID, err := strconv.ParseInt(r.URL.Query().Get("comment_id"), 10, 64)
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

func (h *Handler) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.ParseInt(r.URL.Query().Get("comment_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect comment_id", http.StatusBadRequest)
		return
	}

	userID, reviewID, err := repository.GetUserIdAndReviewIDByComment(h.db, commentID)
	if err != nil {
		if err.Error() == "incorrect commentID" {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if claims.Role != "user" || claims.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = repository.UpdateComment(h.db, commentID, req)
	if err != nil {
		if err.Error() == "comment not found" {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating comment", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	cacheKey := "comments_review" + strconv.FormatInt(reviewID, 10)
	h.redis.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
}
