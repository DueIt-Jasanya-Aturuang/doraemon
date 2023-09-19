package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisImpl struct {
	Client *redis.Client
}

func NewRedisConnection() *RedisImpl {
	host := fmt.Sprintf("%s:%s", RedisHost, RedisPort)
	rDB := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: RedisPass,
		DB:       RedisDB,
	})

	ctx := context.TODO()
	ping, err := rDB.Ping(ctx).Result()
	if err != nil {
		log.Err(err).Msg("ping redis error")
		os.Exit(1)
	}

	log.Info().Msgf("connection to redis : %s", ping)

	return &RedisImpl{
		Client: rDB,
	}
}
