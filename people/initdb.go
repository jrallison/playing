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
		check(db.Exec("CREATE TABLE IF NOT EXISTS " + schema + ".people (id serial PRIMARY KEY, internal varchar UNIQUE, external varchar UNIQUE)"))
		check(db.Exec("CREATE TABLE IF NOT EXISTS " + schema + ".attributes (id integer references " + schema + ".people ON DELETE CASCADE, name varchar, value varchar, timestamp integer, PRIMARY KEY (id, name))"))
		check(db.Exec("CREATE TABLE IF NOT EXISTS " + schema + ".segments (id integer references " + schema + ".people ON DELETE CASCADE, segment_id integer, member boolean, timestamp integer, PRIMARY KEY (id, segment_id))"))

		if index {
			check(db.Exec("CREATE INDEX attrs_idx on " + schema + ".attributes (id, name, value)"))
			check(db.Exec("CREATE INDEX segs_idx ON " + schema + ".segments (id, segment_id, member, timestamp)"))
		}
	}

	return db
}
