package usecases

import (
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

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
