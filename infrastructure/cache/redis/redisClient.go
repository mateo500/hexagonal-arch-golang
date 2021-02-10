package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func NewRedisClient(redisUrl string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cache.redisClient.FindNewRedisClientById")
	}

	client := redis.NewClient(options)

	_, err = client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "cache.redisClient.FindNewRedisClientById")
	}

	return client, nil
}
