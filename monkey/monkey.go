package monkey

import (
	"math/rand"
	"sync"
	"time"
)

// Monkey - used to keep track of performance and monkey bio
type Monkey struct {
	ID        int
	Name      string
	Highwater int
	Speed     float64
	Profile   map[string]int
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
func (monkey Monkey) StartTyping(target string, updates chan Report, done *sync.WaitGroup, outputTimer chan SpeedReport) {
	defer done.Done()
	rand.Seed(time.Now().UnixNano() / (int64(monkey.ID) + 1)) // has to be `id+1` because we have an id 0

	possibilities := convertTypingProfile(monkey.Profile)

	timer := make(chan int, 1000)

	currentSearch := 0
	highwater := -1
	tickLevel := 10000000

	go typingRate(timer, monkey.ID, outputTimer)

	for i := 0; true; i++ {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentSearch] {
			if currentSearch > highwater {
				highwater = currentSearch
				updates <- Report{monkey.ID, highwater}
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

// ConstructTypingProfile - randomize how frequently a monkey hits a particular letter
func ConstructTypingProfile() map[string]int {
	width := 3 // how varied each letter could be

	profile := map[string]int{
		"a": rand.Intn(width) + 4,
		"b": rand.Intn(width) + 1,
		"c": rand.Intn(width) + 1,
		"d": rand.Intn(width) + 1,
		"e": rand.Intn(width) + 4,
		"f": rand.Intn(width) + 1,
		"g": rand.Intn(width) + 1,
		"h": rand.Intn(width) + 1,
		"i": rand.Intn(width) + 4,
		"j": rand.Intn(width) + 1,
		"k": rand.Intn(width) + 1,
		"l": rand.Intn(width) + 1,
		"m": rand.Intn(width) + 1,
		"n": rand.Intn(width) + 1,
		"o": rand.Intn(width) + 4,
		"p": rand.Intn(width) + 1,
		"q": rand.Intn(width) + 1,
		"r": rand.Intn(width) + 1,
		"s": rand.Intn(width) + 1,
		"t": rand.Intn(width) + 1,
		"u": rand.Intn(width) + 4,
		"v": rand.Intn(width) + 1,
		"w": rand.Intn(width) + 1,
		"x": rand.Intn(width) + 1,
		"y": rand.Intn(width) + 1,
		"z": rand.Intn(width) + 1,
		" ": rand.Intn(width) + 7,
	}
	return profile
}

func convertTypingProfile(profile map[string]int) string {
	answer := ""
	for char, weight := range profile {
		for i := 0; i < weight; i++ {
			answer = answer + char
		}
	}
	return answer
}
