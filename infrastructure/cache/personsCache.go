package cache

import "persons.com/api/domain/person"

type PersonsCache interface {
	Set(key string, person *person.Person) error
	SetAll(key string, persons []*person.Person) error
	Get(key string) (*person.Person, error)
	GetAll(key string) ([]*person.Person, error)
}
