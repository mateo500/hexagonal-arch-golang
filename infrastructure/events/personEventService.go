package events

type PersonEventService interface {
	Publish(target string, busName string, value interface{}) error
}
