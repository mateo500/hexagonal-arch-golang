package mysql

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
	"persons.com/api/domain/person"
)

type mysqlRepository struct {
	mysqlClient *sql.DB
}

func (m *mysqlRepository) NewMysqlClient(dbUser string, dbPass string, dbName string) (*sql.DB, error) {

	database, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMysqlClient")
	}

	sql := `DROP TABLE IF EXISTS persons;
	CREATE TABLE persons (
		id int(15) NOT NULL,
	  	name varchar(30) NOT NULL,
	  	last_name varchar(30) NOT NULL,
		age int(3) NOT NULL,  
	  	PRIMARY KEY (id)
	);`

	statement, err := database.Prepare(sql)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Person.Create")
	}

	statement.Exec()

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
		VALUES ($1, $2, $3, $4)
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

	defer m.mysqlClient.Close()
	return nil
}

func (m *mysqlRepository) FindById(id string) (*person.Person, error) {

}
