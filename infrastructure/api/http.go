package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
	jsonSerializer "persons.com/api/infrastructure/serializer/json"
)

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
		log.Println(err)
	}
}

func (h *Handler) serializer(contentType string) person.PersonSerializer {
	return &jsonSerializer.Person{}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	id := chi.URLParam(r, "id")
	personFound, err := h.personService.FindById(id)

	if err != nil {
		if errors.Cause(err) == person.ErrPersonNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(personFound)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setupResponse(w, contentType, responseBody, http.StatusOK)

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	personsCollection, err := h.personService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).EncodeMultiple(personsCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setupResponse(w, contentType, responseBody, http.StatusCreated)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newPerson, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.personService.Create(newPerson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(newPerson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
