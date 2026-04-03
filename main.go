package main

import (
	"comment-service/auth"
	"comment-service/cache"
	DB "comment-service/db"
	"comment-service/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	cnf := DB.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := DB.ConnectWithRetry(cnf)
	if err != nil {
		log.Fatal(err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := cache.NewRedis(redisAddr)
	if err := cache.ConnectRedisWithRetry(rdb); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	publicKey, err := auth.LoadPublicKey("./keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := auth.AuthMiddleware(publicKey)

	h := handlers.New(db, rdb)
	mux := http.NewServeMux()
	mux.Handle("/comment/update/status", authMiddleware(http.HandlerFunc(h.UpdateCommentStatusHandler)))
	log.Println("Comment service started on port 8080")
	log.Println(http.ListenAndServe(":8080", mux))
}
