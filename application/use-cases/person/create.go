package usecases

import (
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

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
