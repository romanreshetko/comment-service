package handlers

import (
	"comment-service/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) getCommentsByReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	nocache := r.URL.Query().Get("nocache") == "true"
	ctx := r.Context()
	cacheKey := "comments_review" + strconv.FormatInt(reviewID, 10)

	if !nocache {
		cached, err := h.redis.Get(ctx, cacheKey).Result()

		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(cached)); err == nil {
				return
			}
			log.Printf("failed to write cache response: %v", err)
		}
		log.Println("No cache, going to DB")
	}

	comments, err := repository.GetCommentsByReview(h.db, reviewID)
	if err != nil {
		log.Println("search comment error", err)
		http.Error(w, "search comments error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
