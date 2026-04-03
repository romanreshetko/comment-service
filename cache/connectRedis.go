package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func ConnectRedisWithRetry(rdb *redis.Client) error {
	var err error

	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Println("connected to redis")
			return nil
		}

		log.Printf("redis not ready, retry %d/5: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return err
}
