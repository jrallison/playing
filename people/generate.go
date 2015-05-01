package people

import (
	"math/rand"
	"strconv"
	"time"
)

// Person contains generated information about person record
type Person struct {
	Internal    int
	External    string
	Attributes  map[string]string
	Memberships map[string]int
}

// Generate builds 'count' persons based on random data
func Generate(count, attrs, segments int) <-chan Person {
	ret := make(chan Person)

	go (func() {
		for i := 0; i < count; i++ {
			p := Person{
				i + 1,
				"e" + strconv.Itoa(i),
				make(map[string]string),
				make(map[string]int),
			}

			for j := 0; j < attrs; j++ {
				p.Attributes["attr"+strconv.Itoa(j)] = "value" + strconv.Itoa(rand.Intn(attrs*10))
			}

			for j := 0; j < segments; j++ {
				status := "entered|"

				if rand.Intn(2) == 0 {
					status = "left|"
				}

				p.Memberships[status+strconv.Itoa(j)] = int(time.Now().Unix()) - rand.Intn(24*60*60)
			}

			ret <- p
		}

		close(ret)
	})()

	return ret
}
