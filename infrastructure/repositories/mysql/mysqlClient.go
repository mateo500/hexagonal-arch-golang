package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func NewMysqlClient(dbUser string, dbPass string, dbName string) (*sql.DB, error) {

	database, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mysqlClient.NewMysqlClient")
	}

	return database, nil
}
