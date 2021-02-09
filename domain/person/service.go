package person

type PersonService interface {
	FindById(id string) (*Person, error)
	Create(person *Person) error
}
