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
	profile   []int
	seated    bool
	client    *client
}

// AddSeatInput is the information sent from a client requesting
// to add a new, empty seat for a monkey
type AddSeatInput struct {
	Layout string
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
	layout string
	monkey *Monkey
}

var seats map[int]seat // keeping track of who's sitting where

// createNew spawns a new monkey that's already typing
func (client *client) createNew(name string, id int) *Monkey {
	newMonkey := Monkey{
		id:        id,
		name:      name,
		highwater: -1,
		profile:   constructTypingProfile(),
		seated:    false,
		client:    client,
	}
	Bullpen = append(Bullpen, &newMonkey)
	return &newMonkey
}

// startTyping is a method that tells a monkey to start simulating
// key presses.
func (monkey *Monkey) startTyping(seat int) {
	rand.Seed(time.Now().UnixNano() / (int64(monkey.id) + 1)) // has to be `id+1` because we have an id 0
	possibilities := convertTypingProfile(monkey.profile, seats[seat].layout)
	timer := make(chan int, 1000)
	currentSearch := 0
	tickLevel := 1000 * 1000

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
				layout: seats[i].layout,
				monkey: nil,
			}
			break
		}
	}
	monkey.seated = false
	return nil
}

func (monkey *Monkey) sit() error {
	var i int
	for i = 0; i < len(seats); i++ {
		if seats[i].monkey == nil {
			seats[i] = seat{
				layout: seats[i].layout,
				monkey: monkey,
			}
			break
		}
	}
	monkey.seated = true
	monkey.startTyping(i)
	return nil
}

// typingRate accepts an input channel to which a monkey sends a
// "tick" once every million key presses. The function then determines
// that monkey's typing speed and generates a report that's sent
// back to the parent process.
func typingRate(tick chan int, id int, output chan speedReport) {
	start := time.Now()
	for message := range tick {
		end := time.Now()
		elapsed := end.Sub(start)
		start = time.Now()
		// We divide by 1,000 here to get a more readable number. We COULD
		// do the same thing by multiplying the "tickLevel" by 1000, but that
		// means getting updates 1,000 times less often. This does the same thing
		// but updates faster.
		output <- speedReport{id, float64(message) / (elapsed.Seconds() * 1000)}
	}
}

// constructTypingProfile randomizes how frequently a monkey hits a particular letter
func constructTypingProfile() (answer []int) {
	width := 3 // how varied each letter could be

	for i := 0; i < 27; i++ {
		answer = append(answer, rand.Intn(width)+1)
	}
	return
}

// convertTypingProfile translates the weight map generated by
// constructTypingProfile() into a long string that can be selected
// from at random to implement the weighting.
func convertTypingProfile(profile []int, layout string) string {
	answer := ""

	keyboards := map[string]string{
		"qwerty":  "qwertyuiopasdfghjklzxcvbnm ",
		"dvorak":  "pyfgcrlaoeuidhtnsqjkxbmwvz ",
		"colemak": "qwfpgjluyarstdhneiozxcvbkm ",
	}
	for i := 0; i < 27; i++ {
		for j := 0; j < profile[i]; j++ {
			answer = answer + string(keyboards[layout][i])
		}
	}
	return answer
}
