package main

import (
	"database/sql"
	"flag"
	"log"
	"strconv"

	"github.com/jrallison/playing/people"
	_ "github.com/lib/pq"
)

var dbname = flag.String("db", "", "name of database")
var schema = flag.Int("schema", 0, "number of the database schema (defaults to 0)")
var count = flag.Int("count", 1000000, "number of people to insert (defaults to 1,000,000)")
var segment = flag.Int("segment", 1, "segment id to add (defaults to 1)")
var percent = flag.Int("percentage", 50, "integer precentage of customers in the segment")

func main() {
	flag.Parse()

	if *dbname == "" {
		log.Fatal("Must provide at least a database name")
	}

	db, err := sql.Open("postgres", "dbname="+*dbname+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	schemaName := "s" + strconv.Itoa(*schema)

	people.CreateSegment(db, schemaName, *count, *segment, *percent)
}
