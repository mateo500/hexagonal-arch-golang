package redis

import (
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisClient(redisUrl string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisUrl)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(redisUrl string) (person.PersonRepository, error) {
	repository := &redisRepository{}
	client, err := NewRedisClient(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}

	repository.client = client
	return repository, nil
}

func (r *redisRepository) Create(person *person.Person) error {
	key := person.ID
	data := map[string]interface{}{
		"id":       person.ID,
		"name":     person.Name,
		"lastName": person.LastName,
		"age":      person.Age,
	}

	_, err := r.client.HMSet(key, data).Result()

	if err != nil {
		return errors.Wrap(err, "repository.Person.Create")
	}

	return nil

}

func (r *redisRepository) FindById(id string) (*person.Person, error) {
	newPerson := &person.Person{}
	key := id

	data, err := r.client.HGetAll(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.Person.FindById")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(person.ErrPersonNotFound, "repository.Person.FindById")
	}

	age, err := strconv.Atoi(data["age"])

	newPerson.ID = data["id"]
	newPerson.Name = data["name"]
	newPerson.LastName = data["lastName"]
	newPerson.Age = age

	return newPerson, nil
}

func (r *redisRepository) GetAll() ([]*person.Person, error) {

	return nil, errors.Wrap(errors.New("method not supported with redis repository"), "repository.Person.GetAll")
}
