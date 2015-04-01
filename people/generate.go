package people

import (
	"math/rand"
	"strconv"
	"time"
)

type Person struct {
	Internal    int
	External    string
	Attributes  map[string]string
	Memberships map[string]string
}

func Generate(count, attrs, segments int) <-chan Person {
	ret := make(chan Person)

	go (func() {
		for i := 0; i < count; i++ {
			p := Person{
				i + 1,
				"e" + strconv.Itoa(i),
				make(map[string]string),
				make(map[string]string),
			}

			for j := 0; j < attrs; j++ {
				p.Attributes["attr"+strconv.Itoa(j)] = "value" + p.External + strconv.Itoa(rand.Intn(attrs*10))
			}

			for j := 0; j < segments; j++ {
				status := "entered|"

				if rand.Intn(2) == 0 {
					status = "left|"
				}

				p.Memberships[strconv.Itoa(j)] = status + strconv.Itoa(int(time.Now().Unix())-rand.Intn(24*60*60))
			}

			ret <- p
		}

		close(ret)
	})()

	return ret
}
