package handlers

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	db    *sql.DB
	redis *redis.Client
}

func New(db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{db, redis}
}
