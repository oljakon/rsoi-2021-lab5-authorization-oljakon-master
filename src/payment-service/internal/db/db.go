package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const dsn = "postgresql://hippo:,7_lQeiIM%5DgiO%3EF9danx%29%3Aoh@hippo-primary.postgres-operator.svc:5432/hippo"

//const dsn = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

// CreateConnection to persons db
func CreateConnection() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
