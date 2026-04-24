package handlers

import (
	"comment-service/models"
	"comment-service/repository"
	serviceIntegrations "comment-service/service-integrations"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	err = repository.DeleteComment(h.db, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	cacheKey := "comments_review" + strconv.FormatInt(reviewID, 10)
	h.redis.Del(ctx, cacheKey)

	go serviceIntegrations.UpdateUserPoints(claims.UserID, -5)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteCommentsByReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "service" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	err = repository.DeleteCommentsByReview(h.db, reviewID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
