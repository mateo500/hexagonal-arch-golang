package person

type PersonRepository interface {
	FindById(id string) (*Person, error)
	GetAll() ([]*Person, error)
	Create(person *Person) error
}
