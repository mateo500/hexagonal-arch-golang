package messagepack

import (
	"github.com/pkg/errors"
	messagePackEncoder "github.com/vmihailenco/msgpack/v5"
	"persons.com/api/domain/person"
)

type Person struct{}

func (p *Person) Decode(input []byte) (*person.Person, error) {
	person := &person.Person{}

	if err := messagePackEncoder.Unmarshal(input, person); err != nil {
		return nil, errors.Wrap(err, "serializer.MessagePack.Person.Decode")
	}

	return person, nil
}

func (p *Person) Encode(input *person.Person) ([]byte, error) {
	rawMsg, err := messagePackEncoder.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.MessagePack.Person.Encode")
	}

	return rawMsg, nil
}

func (p *Person) EncodeMultiple(input []*person.Person) ([]byte, error) {
	rawMsg, err := messagePackEncoder.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.MessagePack.Person.Encode")
	}

	return rawMsg, nil
}
