package usecases

import (
	"persons.com/api/domain/person"
)

func (p *PersonUseCase) Create(newPerson *person.Person) error {
	err := p.personValidator(newPerson)
	if err != nil {
		return err
	}

	err = p.personService.Create(newPerson)
	if err != nil {
		return err
	}

	err = p.personCache.Set(newPerson.ID, newPerson)
	if err != nil {
		return err
	}

	if newPerson.Age >= person.ColombianAdultAge {

		p.personEventService.CreateAdultPersonEvent(newPerson)
		if err != nil {
			return err
		}

	} else {
		p.personEventService.CreateMinorPersonEvent(newPerson)
		if err != nil {
			return err
		}
	}

	return nil
}
