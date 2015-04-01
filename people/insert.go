package people

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/lib/pq/hstore"
)

func Insert(db *sql.DB, schemas, count, attributes, segments int) {
	type person struct {
		internal    string
		external    string
		attributes  map[string]sql.NullString
		memberships map[string]sql.NullString
	}

	var insert func(string, []person)

	insert = func(schema string, people []person) {
		if len(people) == 0 {
			return
		}

		tx, err := db.Begin()
		exitIf("start transaction", err)

		query := "INSERT INTO " + schema + ".people (internal, external, attributes, memberships) VALUES "
		args := make([]interface{}, 0, len(people)*4)

		for i, person := range people {
			query += fmt.Sprint("($", i*4+1, ", $", i*4+2, ", $", i*4+3, ", $", i*4+4, ")")
			if i != len(people)-1 {
				query += ", "
			}

			args = append(args, person.internal, person.external, hstore.Hstore{person.attributes}, hstore.Hstore{person.memberships})
		}

		r, err := db.Exec(query, args...)
		exitIf("inserting people", err)

		if num, _ := r.RowsAffected(); num != int64(len(people)) {
			log.Fatal("insert didn't insert?", r)
		}

		exitIf("commit transaction", tx.Commit())
	}

	start := time.Now()
	batch := make([]person, 0, 100)

	for i := 0; i < count; i++ {
		if i%10000 == 0 {
			log.Println("inserting person", i)
		}

		p := person{
			"i" + strconv.Itoa(i),
			"e" + strconv.Itoa(i),
			make(map[string]sql.NullString),
			make(map[string]sql.NullString),
		}

		for j := 0; j < attributes; j++ {
			p.attributes["attr"+strconv.Itoa(j)] = sql.NullString{"value" + strconv.Itoa(rand.Intn(attributes*10)), true}
		}

		for j := 0; j < segments; j++ {
			status := "entered|"

			if rand.Intn(2) == 0 {
				status = "left|"
			}

			p.memberships[strconv.Itoa(j)] = sql.NullString{status + strconv.Itoa(int(time.Now().Unix())-rand.Intn(24*60*60)), true}
		}

		batch = append(batch, p)

		if len(batch) >= 100 {
			insert(randomSchema(schemas), batch)
			batch = make([]person, 0, 100)
		}
	}

	insert(randomSchema(schemas), batch)

	log.Println("inserted", count, "persons in", time.Since(start))
}

func randomSchema(n int) string {
	return "s" + strconv.Itoa(rand.Intn(n))
}
