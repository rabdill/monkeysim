package monkey

import (
	"math/rand"
	"time"
)

// Monkey is the main struct used to keep track of
// performance and statistics for an individual typist
type Monkey struct {
	id        int
	name      string
	highwater int
	speed     float64
	profile   map[string]int
	seated    bool
	client    *client
}

// report is the format of the updates monkeys send to
// the parent process when they make progress on their assignment
type report struct {
	id, highwater int
}

// speedReport is the format of the updates monkeys send
// to the parent process about how fast they are typing.
type speedReport struct {
	id    int
	speed float64
}

// client is used for storing parameters we need to make a new monkey
type client struct {
	target      string
	outputTimer chan speedReport
}

type seat struct {
	keyboard string
	monkey   *Monkey
}

var seats map[int]seat // keeping track of who's sitting where

// createNew spawns a new monkey that's already typing
func (client *client) createNew(name string, id int) *Monkey {
	newMonkey := Monkey{
		id:        id,
		name:      name,
		highwater: -1,
		profile:   constructTypingProfile(),
		seated:    true,
		client:    client,
	}
	go newMonkey.startTyping()
	Bullpen = append(Bullpen, &newMonkey)
	return &newMonkey
}

// startTyping is a method that tells a monkey to start simulating
// key presses.
func (monkey *Monkey) startTyping() {
	rand.Seed(time.Now().UnixNano() / (int64(monkey.id) + 1)) // has to be `id+1` because we have an id 0
	possibilities := convertTypingProfile(monkey.profile)
	timer := make(chan int, 1000)
	currentSearch := 0
	tickLevel := 10000000

	go typingRate(timer, monkey.id, monkey.client.outputTimer) // speedReports sent straight to monitoring process

	for i := 0; monkey.seated; i++ {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == monkey.client.target[currentSearch] {
			if currentSearch > monkey.highwater {
				monkey.highwater = currentSearch
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
	close(timer) // don't keep listening to speed reports once the monkey stops typing
}

func (monkey *Monkey) stand() error {
	for i := 0; i < len(seats); i++ {
		if seats[i].monkey == monkey {
			seats[i] = seat{
				keyboard: seats[i].keyboard,
				monkey:   nil,
			}
			break
		}
	}
	monkey.seated = false
	return nil
}

func (monkey *Monkey) sit() error {
	for i := 0; i < len(seats); i++ {
		if seats[i].monkey == nil {
			seats[i] = seat{
				keyboard: seats[i].keyboard,
				monkey:   monkey,
			}
			break
		}
	}
	monkey.seated = true
	monkey.startTyping()
	return nil
}

// typingRate accepts an input channel to which a monkey sends a
// "tick" once every 1,000 key presses. The function then determines
// that monkey's typing speed and generates a report that's sent
// back to the parent process.
func typingRate(tick chan int, id int, output chan speedReport) {
	start := time.Now()
	for message := range tick {
		end := time.Now()
		elapsed := end.Sub(start)
		start = time.Now()
		output <- speedReport{id, float64(message) / (elapsed.Seconds() * 1000)}
	}
}

// constructTypingProfile randomizes how frequently a monkey hits a particular letter
func constructTypingProfile() map[string]int {
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

// convertTypingProfile translates the weight map generated by
// constructTypingProfile() into a long string that can be selected
// from at random to implement the weighting.
func convertTypingProfile(profile map[string]int) string {
	answer := ""
	for char, weight := range profile {
		for i := 0; i < weight; i++ {
			answer = answer + char
		}
	}
	return answer
}
