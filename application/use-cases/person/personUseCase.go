package usecases

import (
	"time"

	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

//usecases port with domain services
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

func (p *PersonUseCase) Create(newPerson *person.Person) error {
	err := p.personValidator(newPerson)
	if err != nil {
		return errors.Wrap(err, "applicacion.personUseCase.Create")
	}

	err = p.personService.Create(newPerson)
	if err != nil {
		return errors.Wrap(err, "applicacion.personUseCase.Create")
	}

	err = p.personCache.Set(newPerson.ID, newPerson)
	if err != nil {
		return errors.Wrap(err, "applicacion.personUseCase.Create")
	}

	if newPerson.Age >= person.ColombianAdultAge {

		p.personEventService.CreateAdultPersonEvent(newPerson)
		if err != nil {
			return errors.Wrap(err, "applicacion.personUseCase.Create")
		}

	} else {
		p.personEventService.CreateMinorPersonEvent(newPerson)
		if err != nil {
			return errors.Wrap(err, "applicacion.personUseCase.Create")
		}
	}

	return nil
}

func (p *PersonUseCase) GetAll() ([]*person.Person, error) {
	var personsFound []*person.Person

	personsInCache, _ := p.personCache.GetAll("personsCache@" + time.Now().Format("2-24-2021"))

	if personsInCache == nil {
		personsCollection, err := p.personService.GetAll()
		if err != nil {
			return nil, errors.Wrap(err, "applicacion.personUseCase.GetAll")
		}

		personsFound = personsCollection

		err = p.personCache.SetAll("personsCache@"+time.Now().Format("2-24-2021"), personsCollection)
		if err != nil {
			return nil, errors.Wrap(err, "applicacion.personUseCase.GetAll")
		}

	} else {
		personsFound = personsInCache

	}

	return personsFound, nil
}

func (p *PersonUseCase) FindById(id string) (*person.Person, error) {

	var personFound *person.Person

	personInCache, _ := p.personCache.Get(id)

	if personInCache == nil {

		personFoundInDb, err := p.personService.FindById(id)
		if err != nil {
			return nil, errors.Wrap(err, "applicacion.personUseCase.FindById")
		}

		personFound = personFoundInDb

		err = p.personCache.Set(personFoundInDb.ID, personFoundInDb)
		if err != nil {
			return nil, errors.Wrap(err, "applicacion.personUseCase.FindById")
		}

	} else {
		personFound = personInCache

	}

	return personFound, nil

}
