package serializers

import "persons.com/api/domain/person"

type PersonSerializer interface {
	Decode(input []byte) (*person.Person, error)
	Encode(input *person.Person) ([]byte, error)
	EncodeMultiple(input []*person.Person) ([]byte, error)
}
