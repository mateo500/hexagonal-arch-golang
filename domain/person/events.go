package person

type PersonEventsService interface {
	CreateMinorPersonEvent(person *Person) error
	CreateAdultPersonEvent(person *Person) error
}
