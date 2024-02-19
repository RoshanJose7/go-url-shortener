package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CreateClient() *sql.DB {
	db, err := sql.Open("sqlite3", "main.db")

	if err != nil {
		panic(err)
	}

	sqlStmt := `create table if not exists urls (
					id integer not null primary key autoincrement,
					url text,
					shortened text,
					expires_at datetime
				);`

	_, err = db.Exec(sqlStmt)

	if err != nil {
		panic(err)
	}

	return db
}
