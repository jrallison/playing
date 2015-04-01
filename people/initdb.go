package people

import (
	"database/sql"
	"log"
	"strconv"
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

func InitDb(name string, schemas int, clean, index bool) *sql.DB {
	db, err := sql.Open("postgres", "dbname=postgres sslmode=disable")
	exitIf("open db", err)

	if clean {
		db.Exec("DROP DATABASE " + name)
	}

	db.Exec("CREATE DATABASE " + name)
	db.Close()

	db, err = sql.Open("postgres", "dbname="+name+" sslmode=disable")
	exitIf("open db", err)

	for i := 0; i < schemas; i++ {
		schema := "s" + strconv.Itoa(i)

		check(db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema))
		check(db.Exec("CREATE EXTENSION IF NOT EXISTS hstore"))
		check(db.Exec("CREATE TABLE IF NOT EXISTS " + schema + ".people (id serial PRIMARY KEY, internal varchar UNIQUE, external varchar UNIQUE, attributes hstore, memberships hstore)"))

		if index {
			check(db.Exec("CREATE OR REPLACE FUNCTION hstore_text(hstore) RETURNS text LANGUAGE sql IMMUTABLE AS $$ SELECT string_agg(key || '-' || value, ' ') from each($1) $$;"))
			check(db.Exec("CREATE INDEX attrs_index on " + schema + ".people using gin(to_tsvector('english', hstore_text(attributes)))"))
		}
	}

	return db
}
