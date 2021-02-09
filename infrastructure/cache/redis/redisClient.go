package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"persons.com/api/infrastructure/cache"
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

func NewRedisCache(host string, expirationTime time.Duration) (cache.PersonsCache, error) {

	client, err := NewRedisClient(host)
	if err != nil {
		return nil, errors.Wrap(err, "cache.redisClient.NewRedisCache")
	}

	return &redisCache{
		host:    host,
		expires: expirationTime,
		client:  client,
	}, nil
}

func GetRedisClient(host string, expirationTime time.Duration) (cache.PersonsCache, error) {
	var redisCache, err = NewRedisCache(host, expirationTime)
	if err != nil {
		return nil, errors.Wrap(err, "cache.redisClient.GetRedisClient")
	}

	return redisCache, nil
}
