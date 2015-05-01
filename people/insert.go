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

// Insert generates 'count' people and inserts them into the people table.
func Insert(db *sql.DB, schemas, count, attributes, segments int) {
	var insert func(string, []Person)

	insert = func(schema string, people []Person) {
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

			attrs := make(map[string]sql.NullString)
			memberships := make(map[string]sql.NullString)

			for n, v := range person.Attributes {
				attrs[n] = sql.NullString{v, true}
			}

			for n, v := range person.Memberships {
				memberships[n] = sql.NullString{strconv.Itoa(v), true}
			}

			args = append(args, person.Internal, person.External, hstore.Hstore{attrs}, hstore.Hstore{memberships})
		}

		r, err := db.Exec(query, args...)
		exitIf("inserting people", err)

		if num, _ := r.RowsAffected(); num != int64(len(people)) {
			log.Fatal("insert didn't insert?", r)
		}

		exitIf("commit transaction", tx.Commit())
	}

	i := 0
	start := time.Now()
	batch := make([]Person, 0, 100)

	for p := range Generate(count, attributes, segments) {
		i++

		if i%10000 == 0 {
			log.Println("inserting person", i)
		}

		batch = append(batch, p)

		if len(batch) >= 100 {
			insert(randomSchema(schemas), batch)
			batch = make([]Person, 0, 100)
		}
	}

	insert(randomSchema(schemas), batch)

	log.Println("inserted", count, "persons in", time.Since(start))
}

func randomSchema(n int) string {
	return "s" + strconv.Itoa(rand.Intn(n))
}
