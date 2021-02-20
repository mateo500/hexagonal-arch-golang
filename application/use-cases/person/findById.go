package usecases

import (
	"persons.com/api/domain/person"
)

func (p *PersonUseCase) FindById(id string) (*person.Person, error) {

	var personFound *person.Person

	personInCache, _ := p.personCache.Get(id)

	if personInCache == nil {

		personFoundInDb, err := p.personService.FindById(id)
		if err != nil {
			return nil, err
		}

		personFound = personFoundInDb

		err = p.personCache.Set(personFoundInDb.ID, personFoundInDb)
		if err != nil {
			return nil, err
		}

	} else {
		personFound = personInCache

	}

	return personFound, nil

}
