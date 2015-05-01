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

// CreateSegment adds a new segment membership to the first 'count' people
func CreateSegment(db *sql.DB, schema string, count, segmentID, percent int) {
	var add func(string, []int)

	add = func(schema string, people []int) {
		if len(people) == 0 {
			return
		}

		tx, err := db.Begin()
		exitIf("start transaction", err)

		query := `
		  UPDATE ` + schema + `.people
			SET memberships = memberships || myvalues.hash::hstore
			FROM (
				VALUES `

		end := `) AS myvalues (id, hash)
		  WHERE ` + schema + `.people.id = myvalues.id::integer`

		args := make([]interface{}, 0, len(people)*2)

		for i, id := range people {
			query += fmt.Sprint("($", i*2+1, ", $", i*2+2, ")")
			if i != len(people)-1 {
				query += ", "
			}

			status := "left|"

			if rand.Intn(100) <= percent {
				status = "entered|"
			}

			key := status + strconv.Itoa(segmentID)
			value := sql.NullString{strconv.Itoa(int(time.Now().Unix())), true}

			args = append(args, id, hstore.Hstore{map[string]sql.NullString{key: value}})
		}

		query += end

		r, err := db.Exec(query, args...)
		exitIf("updating people", err)

		if num, _ := r.RowsAffected(); num != int64(len(people)) {
			log.Fatal("update didn't update?", r)
		}

		exitIf("commit transaction", tx.Commit())
	}

	start := time.Now()
	batch := make([]int, 0, 100)

	for i := 0; i < count; i++ {
		id := i + 1

		if i%10000 == 0 {
			log.Println("adding to person", id)
		}

		batch = append(batch, id)

		if len(batch) >= 100 {
			add(schema, batch)
			batch = make([]int, 0, 100)
		}
	}

	add(schema, batch)

	log.Println("updated", count, "persons in", time.Since(start))
}
