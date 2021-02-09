package person

import (
	"errors"

	"github.com/teris-io/shortid"
)

var ErrPersonNotFound = errors.New("Person not found")
var ColombianAdultAge int = 18

type personService struct {
	personRepo PersonRepository
}

func NewPersonService(personRepo PersonRepository) PersonService {
	return &personService{
		personRepo,
	}
}

func (p *personService) Create(person *Person) error {
	person.ID = shortid.MustGenerate()
	return p.personRepo.Create(person)
}

func (p *personService) FindById(id string) (*Person, error) {
	return p.personRepo.FindById(id)
}
