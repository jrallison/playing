package main

import (
	"flag"
	"log"

	"github.com/jrallison/playingwithpostgres/people"
	_ "github.com/lib/pq"
)

var dbname = flag.String("db", "", "name of database")
var schemas = flag.Int("schemas", 1, "number of database schemas (defaults to 1)")
var clean = flag.Bool("clean", false, "drop database before starting (defaults to false)")
var index = flag.Bool("index", false, "whether or not to create full text indexes for attributes/segments (defaults to false)")
var count = flag.Int("count", 1000000, "number of people to insert (defaults to 1,000,000)")
var segments = flag.Int("segments", 200, "number of segments to create memberships (defaults to 200)")
var attributes = flag.Int("attributes", 50, "number of attributes to create per person (defaults to 50)")

func main() {
	flag.Parse()

	if *dbname == "" {
		log.Fatal("Must provide at least a database name")
	}

	db := people.InitDb(*dbname, *schemas, *clean, *index)
	people.Insert(db, *schemas, *count, *attributes, *segments)
}
