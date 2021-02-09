package person

type PersonRepository interface {
	FindById(id string) (*Person, error)
	Create(person *Person) error
}
