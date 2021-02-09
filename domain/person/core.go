package person

import (
	"errors"
	"strconv"

	errorDispatcher "github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

const ColombianAdultAge int = 18
const NameMinLength int = 3
const LastNameMinLength int = 1
const uniqueIdMinLength int = 9

var ErrPersonNotFound = errors.New("Person not found")
var ErrInvalidNameLength = errors.New("Invalid name length, must be at least:" + strconv.Itoa(NameMinLength))
var ErrInvalidLastNameLength = errors.New("Invalid last name length, must be at least:" + strconv.Itoa(LastNameMinLength))
var ErrInvalidUniqueIdLength = errors.New("Invalid unique id length, must be at least:" + strconv.Itoa(uniqueIdMinLength))

type personService struct {
	personRepo PersonRepository
}

func NewPersonService(personRepo PersonRepository) PersonService {
	return &personService{
		personRepo,
	}
}

func (p *personService) Create(person *Person) error {

	if len(person.Name) < NameMinLength {
		return errorDispatcher.Wrap(ErrInvalidNameLength, "core.Person.Create")
	}

	if len(person.LastName) < LastNameMinLength {
		return errorDispatcher.Wrap(ErrInvalidLastNameLength, "core.Person.Create")
	}

	person.ID = shortid.MustGenerate()
	return p.personRepo.Create(person)
}

func (p *personService) FindById(id string) (*Person, error) {

	if len(id) < uniqueIdMinLength {
		return nil, errorDispatcher.Wrap(ErrInvalidUniqueIdLength, "core.Person.FindById")
	}

	return p.personRepo.FindById(id)
}

func (p *personService) GetAll() ([]*Person, error) {
	return p.personRepo.GetAll()
}
