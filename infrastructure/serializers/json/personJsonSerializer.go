package json

import (
	"encoding/json"

	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

type Person struct{}

func (p *Person) Decode(input []byte) (*person.Person, error) {
	person := &person.Person{}

	if err := json.Unmarshal(input, person); err != nil {
		return nil, errors.Wrap(err, "serializer.Person.Decode")
	}

	return person, nil
}

func (p *Person) Encode(input *person.Person) ([]byte, error) {
	rawMsg, err := json.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Person.Encode")
	}

	return rawMsg, nil
}

func (p *Person) EncodeMultiple(input []*person.Person) ([]byte, error) {
	rawMsg, err := json.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Person.Encode")
	}

	return rawMsg, nil
}
