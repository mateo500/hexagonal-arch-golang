package validators

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"persons.com/api/domain/person"
)

func PersonValidator(person *person.Person) error {
	return validation.ValidateStruct(person,
		validation.Field(&person.Name, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]")), is.Alpha),
		validation.Field(&person.LastName, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]")), is.Alpha),
	)
}
