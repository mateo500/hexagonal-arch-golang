package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
	"persons.com/api/infrastructure/cache"
	"persons.com/api/infrastructure/cache/redis"
	"persons.com/api/infrastructure/serializers"
	jsonSerializer "persons.com/api/infrastructure/serializers/json"
	messagepackSerializer "persons.com/api/infrastructure/serializers/messagePack"
	"persons.com/api/infrastructure/validators"
)

var personsCache = func() cache.PersonsCache {
	redisClient, err := redis.GetRedisClient(os.Getenv("CACHE_DB_URL"), 60)
	if err != nil {
		log.Fatal(err)
	}

	return redisClient
}()

type PersonHandler interface {
	GetById(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	GetAll(http.ResponseWriter, *http.Request)
}

type Handler struct {
	personService person.PersonService
}

func NewHandler(personService person.PersonService) PersonHandler {
	return &Handler{personService: personService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("error seting up http response: %v", err)
	}
}

func (h *Handler) serializer(contentType string) serializers.PersonSerializer {

	if contentType == "application/x-msgpack" {
		return &messagepackSerializer.Person{}
	}

	return &jsonSerializer.Person{}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	id := chi.URLParam(r, "id")

	var personFound *person.Person

	personInCache, err := personsCache.Get(id)

	if personInCache == nil {
		personFoundInDb, err := h.personService.FindById(id)
		if err != nil {
			if errors.Cause(err) == person.ErrPersonNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			InternalServerError(err, w)
		}
		personFound = personFoundInDb

		err = personsCache.Set(personFoundInDb.ID, personFoundInDb)
		InternalServerError(err, w)
	} else {
		personFound = personInCache
	}

	responseBody, err := h.serializer(contentType).Encode(personFound)
	InternalServerError(err, w)

	setupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	var personsFound []*person.Person

	personsInCache, err := personsCache.GetAll("personsCache@" + time.Now().Format("2-24-2021"))
	if err == nil {
		InternalServerError(err, w)
	}

	if personsInCache == nil {
		personsCollection, err := h.personService.GetAll()
		InternalServerError(err, w)

		personsFound = personsCollection
		err = personsCache.SetAll("personsCache@"+time.Now().Format("2-24-2021"), personsCollection)
		InternalServerError(err, w)
	} else {
		personsFound = personsInCache
	}

	responseBody, err := h.serializer(contentType).EncodeMultiple(personsFound)
	InternalServerError(err, w)

	setupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	requestBody, err := ioutil.ReadAll(r.Body)
	InternalServerError(err, w)

	newPerson, err := h.serializer(contentType).Decode(requestBody)
	InternalServerError(err, w)

	err = validators.PersonValidator(newPerson)
	BadRequest(err, w)

	err = h.personService.Create(newPerson)
	InternalServerError(err, w)

	err = personsCache.Set(newPerson.ID, newPerson)
	InternalServerError(err, w)

	responseBody, err := h.serializer(contentType).Encode(newPerson)
	InternalServerError(err, w)

	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
