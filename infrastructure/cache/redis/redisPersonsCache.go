package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

type redisCache struct {
	host    string
	expires time.Duration
	client  *redis.Client
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
