package usecases

import (
	"time"

	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

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
