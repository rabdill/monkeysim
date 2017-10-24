package monkey

import (
	"math/rand"
	"sync"
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
	rand.Seed(int64(id))
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	currentSearch := 0
	highwater := -1

	// for i := 0; i < 99999999; i++ {
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
		if currentSearch > 0 {
			currentSearch = 0
		}
	}
}
