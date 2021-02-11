package usecases

import (
	"persons.com/api/domain/person"
)

type PersonUseCases interface {
	FindById(id string) (*person.Person, error)
	GetAll() ([]*person.Person, error)
	Create(person *person.Person) error
}

//cache service port
type PersonsCacheService interface {
	Set(key string, person *person.Person) error
	SetAll(key string, persons []*person.Person) error
	Get(key string) (*person.Person, error)
	GetAll(key string) ([]*person.Person, error)
}

//validator port
type PersonValidator func(person *person.Person) error

type PersonUseCase struct {
	personService      person.PersonService
	personCache        PersonsCacheService
	personEventService person.PersonEventsService
	personValidator    PersonValidator
}

func NewPersonUseCases(personService person.PersonService, personEventService person.PersonEventsService, personCache PersonsCacheService, personValidator PersonValidator) PersonUseCases {
	return &PersonUseCase{
		personService:      personService,
		personCache:        personCache,
		personEventService: personEventService,
		personValidator:    personValidator,
	}
}
