package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	usecases "persons.com/api/application/use-cases/person"
	"persons.com/api/domain/person"
)

type redisCache struct {
	host    string
	expires time.Duration
	client  *redis.Client
}

func NewRedisPersonsCache(host string, expirationTime time.Duration) (usecases.PersonsCacheService, error) {

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

func GetRedisClient(host string, expirationTime time.Duration) (usecases.PersonsCacheService, error) {
	var redisCache, err = NewRedisPersonsCache(host, expirationTime)
	if err != nil {
		return nil, errors.Wrap(err, "cache.redisClient.GetRedisClient")
	}

	return redisCache, nil
}

func (c *redisCache) Set(key string, person *person.Person) error {

	json, err := json.Marshal(person)
	if err != nil {
		return errors.Wrap(err, "cache.Person.Set")
	}

	var _, errorRedis = c.client.Set(key, json, c.expires*time.Second).Result()
	if errorRedis != nil {
		return errors.Wrap(err, "cache.Person.Set")
	}

	return nil

}

func (c *redisCache) Get(key string) (*person.Person, error) {
	personFound := &person.Person{}

	data, err := c.client.Get(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(person.ErrPersonNotFound, "cache.Person.Get")
	}

	err = json.Unmarshal([]byte(data), personFound)
	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	return personFound, nil
}

func (c *redisCache) GetAll(key string) ([]*person.Person, error) {
	personsFound := make([]*person.Person, 0)

	data, err := c.client.Get(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(person.ErrPersonNotFound, "cache.Person.Get")
	}

	err = json.Unmarshal([]byte(data), &personsFound)
	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	if err != nil {
		return nil, errors.Wrap(err, "cache.Person.Get")
	}

	return personsFound, nil
}

func (c *redisCache) SetAll(key string, persons []*person.Person) error {

	json, err := json.Marshal(persons)
	if err != nil {
		return errors.Wrap(err, "cache.Person.Set")
	}

	var _, errorRedis = c.client.Set(key, json, c.expires*time.Second).Result()
	if errorRedis != nil {
		return errors.Wrap(err, "cache.Person.Set")
	}

	return nil

}
