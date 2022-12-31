package infrastructure

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func OpenRedis() *redis.Client {
	addrRedis := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	passwordRedis := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addrRedis,
		Password: passwordRedis, // no password set
		DB:       1,             // use default DB
	})

	return rdb
}
