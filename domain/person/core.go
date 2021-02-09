package person

import (
	"errors"

	"github.com/teris-io/shortid"
)

const ColombianAdultAge int = 18
const NameMinLength int = 3
const LastNameMinLength int = 1

var ErrPersonNotFound = errors.New("Person not found")

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

func (p *personService) GetAll() ([]*Person, error) {
	return p.personRepo.GetAll()
}
