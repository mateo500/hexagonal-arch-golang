package cache

import "persons.com/api/domain/person"

type PersonsCache interface {
	Set(key string, value *person.Person) error
	Get(key string) (*person.Person, error)
}
