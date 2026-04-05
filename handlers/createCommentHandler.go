package handlers

import (
	"comment-service/models"
	"comment-service/repository"
	"database/sql"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type Handler struct {
	db    *sql.DB
	redis *redis.Client
}

func New(db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{db, redis}
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = repository.CreateComment(h.db, reviewID, claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	cacheKey := "comments_review" + strconv.FormatInt(reviewID, 10)
	h.redis.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusCreated)
}
