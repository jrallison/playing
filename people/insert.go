package people

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Insert(db *sql.DB, schemas, count, attributes, segments int) {
	var insert func(string, []Person)

	insert = func(schema string, people []Person) {
		if len(people) == 0 {
			return
		}

		tx, err := db.Begin()
		exitIf("start transaction", err)

		query1 := "INSERT INTO " + schema + ".people (internal, external) VALUES "
		args1 := make([]interface{}, 0, len(people)*2)

		for i, person := range people {
			query1 += fmt.Sprint("($", i*2+1, ", $", i*2+2, "), ")
			args1 = append(args1, person.Internal, person.External)
		}

		var firstid int
		err = tx.QueryRow(strings.TrimSuffix(query1, ", ")+" RETURNING id", args1...).Scan(&firstid)
		exitIf("inserting people", err)

		aCount := 0
		query2 := "INSERT INTO " + schema + ".attributes (id, name, value, timestamp) VALUES "
		args2 := make([]interface{}, 0, len(people)*attributes*4)

		sCount := 0
		query3 := "INSERT INTO " + schema + ".segments (id, segment_id, member, timestamp) VALUES "
		args3 := make([]interface{}, 0, len(people)*segments*3)

		for i, person := range people {
			for n, v := range person.Attributes {
				query2 += fmt.Sprint("($", aCount*4+1, ", $", aCount*4+2, ", $", aCount*4+3, ", $", aCount*4+4, "), ")
				args2 = append(args2, firstid+i, n, v, int(time.Now().Unix()))
				aCount += 1
			}

			for n, v := range person.Memberships {
				parts := strings.SplitN(v, "|", 2)
				member := parts[0] == "entered"
				ts, _ := strconv.Atoi(parts[1])
				id, _ := strconv.Atoi(n)

				query3 += fmt.Sprint("($", sCount*4+1, ", $", sCount*4+2, ", $", sCount*4+3, ", $", sCount*4+4, "), ")
				args3 = append(args3, firstid+i, id, member, ts)
				sCount += 1
			}
		}

		if len(args2) > 0 {
			_, err = tx.Exec(strings.TrimSuffix(query2, ", "), args2...)
			exitIf("inserting attributes", err)
		}

		if len(args3) > 0 {
			_, err = tx.Exec(strings.TrimSuffix(query3, ", "), args3...)
			exitIf("inserting segments", err)
		}

		exitIf("commit transaction", tx.Commit())
	}

	i := 0
	start := time.Now()
	batch := make([]Person, 0, 100)

	for p := range Generate(count, attributes, segments) {
		i += 1

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
