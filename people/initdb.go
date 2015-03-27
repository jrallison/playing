package people

import (
	"database/sql"
	"log"
)

func exitIf(context string, err error) {
	if err != nil {
		log.Fatal(context, " ", err)
	}
}

func check(r sql.Result, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitDb(name, schema string, clean bool) *sql.DB {
	db, err := sql.Open("postgres", "dbname=postgres sslmode=disable")
	exitIf("open db", err)

	if clean {
		db.Exec("DROP DATABASE " + name)
	}

	db.Exec("CREATE DATABASE " + name)
	db.Close()

	db, err = sql.Open("postgres", "dbname="+name+" sslmode=disable")
	exitIf("open db", err)

	check(db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema))
	check(db.Exec("CREATE EXTENSION IF NOT EXISTS hstore"))
	check(db.Exec("CREATE TABLE IF NOT EXISTS " + schema + ".people (id serial PRIMARY KEY, internal varchar UNIQUE, external varchar UNIQUE, attributes hstore, memberships hstore)"))

	return db
}
