package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

func ConnectToRedis(host, password string, port int, log *slog.Logger) *redis.Client {
	const requestTimeout = 3 * time.Second

	const maxRetries = 5

	const retryDelay = 3 * time.Second

	var counts uint8

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	for {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: password,
			DB:       0,
		})
		_, err := rdb.Ping(ctx).Result()

		if err != nil {
			log.Warn("Redis not connected yet", err)
			counts++
		} else {
			log.Info("Connected to Redis!")
			return rdb
		}
		if counts > maxRetries {
			log.Error("Unable connect to Redis", err)
			return nil
		}

		time.Sleep(retryDelay)

		continue
	}

}
