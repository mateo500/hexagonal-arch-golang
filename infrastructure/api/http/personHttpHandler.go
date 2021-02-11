package http

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	usecases "persons.com/api/application/use-cases/person"
	"persons.com/api/domain/person"
	httpUtils "persons.com/api/infrastructure/api/http/utils"
	"persons.com/api/infrastructure/api/serializers"
	jsonSerializer "persons.com/api/infrastructure/api/serializers/json"
	messagepackSerializer "persons.com/api/infrastructure/api/serializers/messagePack"
	"persons.com/api/infrastructure/env"
	"persons.com/api/infrastructure/validators"
)

var envMap map[string]string = env.NewEnvService().GetEnvs(os.Getenv("APP_MODE"))

type Handler struct {
	personUseCases usecases.PersonUseCases
}

func NewHandler(personService person.PersonService, personEventsService person.PersonEventsService, personsCacheService usecases.PersonsCacheService) PersonHandler {

	usecases := usecases.NewPersonUseCases(personService, personEventsService, personsCacheService, validators.PersonValidator)

	return &Handler{
		personUseCases: usecases,
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

	personFound, err := h.personUseCases.FindById(id)
	if err != nil {
		if errors.Cause(err) == person.ErrPersonNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		internalServerError(err, w)
	}

	responseBody, err := h.serializer(contentType).Encode(personFound)
	internalServerError(err, w)

	httpUtils.SetupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	personsFound, err := h.personUseCases.GetAll()
	internalServerError(err, w)

	responseBody, err := h.serializer("application/json").EncodeMultiple(personsFound)
	internalServerError(err, w)

	httpUtils.SetupResponse(w, "application/json", responseBody, http.StatusOK)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	requestBody, err := ioutil.ReadAll(r.Body)
	internalServerError(err, w)

	newPerson, err := h.serializer(contentType).Decode(requestBody)
	internalServerError(err, w)

	if newPerson != nil {
		err := h.personUseCases.Create(newPerson)
		internalServerError(err, w)

		responseBody, err := h.serializer(contentType).Encode(newPerson)
		internalServerError(err, w)

		httpUtils.SetupResponse(w, contentType, responseBody, http.StatusCreated)

	} else {
		internalServerError(errors.New("error deserializing payload"), w)
	}

}
