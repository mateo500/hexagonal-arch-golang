package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

type mysqlRepository struct {
	mysqlClient *sql.DB
}

func NewMysqlClient(dbUser string, dbPass string, dbName string) (*sql.DB, error) {

	database, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMysqlClient")
	}

	return database, nil
}

func NewMysqlRepository(db *sql.DB) person.PersonRepository {
	return &mysqlRepository{
		mysqlClient: db,
	}
}

func (m *mysqlRepository) Create(person *person.Person) error {
	sql := `
		INSERT INTO persons (id, name, last_name, age)
		VALUES (?, ?, ?, ?)
	`

	statement, err := m.mysqlClient.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "repository.Person.Create")
	}

	result, err := statement.Exec(person.ID, person.Name, person.LastName, person.Age)
	if err != nil {
		return errors.Wrap(err, "repository.Person.Create")
	}

	log.Println(result)

	return nil
}

func (m *mysqlRepository) FindById(id string) (*person.Person, error) {

	newPerson := &person.Person{}

	sql := "SELECT * FROM persons WHERE id = ?"

	result, err := m.mysqlClient.Query(sql, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Person.FindById")
	}

	for result.Next() {
		var age int
		var id, name, last_name string

		err := result.Scan(&id, &name, &last_name, &age)
		if err != nil {
			return nil, errors.Wrap(err, "repository.Person.FindById")
		}

		newPerson.ID = id
		newPerson.Age = age
		newPerson.LastName = last_name
		newPerson.Name = name
	}

	return newPerson, nil

}
