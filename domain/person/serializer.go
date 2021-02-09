package person

type PersonSerializer interface {
	Decode(input []byte) (*Person, error)
	Encode(input *Person) ([]byte, error)
	EncodeMultiple(input []*Person) ([]byte, error)
}
