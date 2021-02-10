package events

type EventService interface {
	Publish(target string, busName string, value interface{}) error
}
