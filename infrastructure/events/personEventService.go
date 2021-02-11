package events

type PersonEventService interface {
	Publish(exchange string, Qname string, value interface{}) error
}
