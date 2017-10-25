package monkey

import (
	"math/rand"
	"sync"
	"time"
)

// Monkey - used to keep track of performance and monkey bio
type Monkey struct {
	Name      string
	Highwater int
}

// Report - how monkeys check in with the master process
type Report struct {
	ID, Highwater int
}

// StartTyping - frantically typing random characters
func StartTyping(id int, target string, updates chan Report, done *sync.WaitGroup) {
	defer done.Done()
	rand.Seed(time.Now().UnixNano() / (int64(id) + 1)) // has to be `id+1` because we have an id 0
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	currentSearch := 0
	highwater := -1

	for {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentSearch] {
			if currentSearch > highwater {
				highwater = currentSearch
				updates <- Report{id, highwater}
			}
			currentSearch++
			continue
		}
		// if we were on a streak, but it's over
		currentSearch = 0
	}
}
