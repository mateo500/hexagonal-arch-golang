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
	httpUtils "persons.com/api/infrastructure/api/http/utils"
	"persons.com/api/infrastructure/cache"
	"persons.com/api/infrastructure/cache/redis"
	"persons.com/api/infrastructure/env"
	"persons.com/api/infrastructure/events"
	"persons.com/api/infrastructure/events/rabbitmq"
	"persons.com/api/infrastructure/serializers"
	jsonSerializer "persons.com/api/infrastructure/serializers/json"
	messagepackSerializer "persons.com/api/infrastructure/serializers/messagePack"
	"persons.com/api/infrastructure/validators"
)

var envMap map[string]string = env.NewEnvService().GetEnvs(os.Getenv("APP_MODE"))

type PersonHandler interface {
	GetById(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	GetAll(http.ResponseWriter, *http.Request)
}

type Handler struct {
	Service  person.PersonService
	Cache    cache.PersonsCache
	EventBus events.PersonEventService
}

func NewHandler(personService person.PersonService) PersonHandler {

	redisClient, err := redis.GetRedisClient(envMap["CACHE_DB_URL"], 30)
	if err != nil {
		log.Fatal(err)
	}

	rabbitEventsService, err := rabbitmq.NewRabbitMqService(envMap["Q_URL"], []string{"persons"}, []string{"minors", "adults"})
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		Service:  personService,
		Cache:    redisClient,
		EventBus: rabbitEventsService,
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

	personInCache, err := h.Cache.Get(id)

	if personInCache == nil {

		personFoundInDb, err := h.Service.FindById(id)

		if err != nil {
			if errors.Cause(err) == person.ErrPersonNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			internalServerError(err, w)
		}

		personFound = personFoundInDb

		err = h.Cache.Set(personFoundInDb.ID, personFoundInDb)
		internalServerError(err, w)

	} else {
		personFound = personInCache
	}

	responseBody, err := h.serializer(contentType).Encode(personFound)
	internalServerError(err, w)

	httpUtils.SetupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	var personsFound []*person.Person

	personsInCache, _ := h.Cache.GetAll("personsCache@" + time.Now().Format("2-24-2021"))

	if personsInCache == nil {
		personsCollection, err := h.Service.GetAll()
		internalServerError(err, w)

		personsFound = personsCollection
		err = h.Cache.SetAll("personsCache@"+time.Now().Format("2-24-2021"), personsCollection)
		internalServerError(err, w)
	} else {
		personsFound = personsInCache
	}

	responseBody, err := h.serializer(contentType).EncodeMultiple(personsFound)
	internalServerError(err, w)

	httpUtils.SetupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	exchangeType := ""

	requestBody, err := ioutil.ReadAll(r.Body)
	internalServerError(err, w)

	newPerson, err := h.serializer(contentType).Decode(requestBody)
	internalServerError(err, w)

	err = validators.PersonValidator(newPerson)
	badRequest(err, w)

	err = h.Service.Create(newPerson)
	internalServerError(err, w)

	err = h.Cache.Set(newPerson.ID, newPerson)
	internalServerError(err, w)

	if newPerson.Age >= person.ColombianAdultAge {
		exchangeType = "adults"
	} else {
		exchangeType = "minors"
	}

	err = h.EventBus.Publish(exchangeType, "persons", newPerson)
	internalServerError(err, w)

	responseBody, err := h.serializer(contentType).Encode(newPerson)
	internalServerError(err, w)

	httpUtils.SetupResponse(w, contentType, responseBody, http.StatusCreated)
}
