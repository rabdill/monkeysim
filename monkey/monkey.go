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
	Speed     float64
}

// Report - how monkeys check in with the master process
type Report struct {
	ID, Highwater int
}

type SpeedReport struct {
	ID    int
	Speed float64
}

// StartTyping - frantically typing random characters
func StartTyping(id int, target string, updates chan Report, done *sync.WaitGroup, outputTimer chan SpeedReport) {
	defer done.Done()
	rand.Seed(time.Now().UnixNano() / (int64(id) + 1)) // has to be `id+1` because we have an id 0
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	timer := make(chan int, 1000)

	currentSearch := 0
	highwater := -1
	tickLevel := 10000000

	go typingRate(timer, id, outputTimer)

	for i := 0; true; i++ {
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
		if i > tickLevel-1 {
			timer <- i
			i = 0
		}
	}
}

func typingRate(tick chan int, id int, output chan SpeedReport) {
	start := time.Now()
	for message := range tick {
		end := time.Now()
		elapsed := end.Sub(start)
		start = time.Now()
		output <- SpeedReport{id, float64(message) / (elapsed.Seconds() * 1000)}
	}
}
