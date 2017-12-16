package monkey

import (
	"fmt"
)

// Bullpen holds the current roster of monkeys and their stats.
var Bullpen []*Monkey

//Target holds the goal text that the monkeys are working toward.
var Target string

var monkeyClient client

// KickOffSim starts the simulation and sets the first set of
// monkeys off to typing.
func KickOffSim() {
	seatCount := getSeatCount()
	seats = make(map[int]seat)
	Bullpen = []*Monkey{}
	Target = getTarget("target.txt")
	speedReports := make(chan speedReport, 500) // receiving speed reading from monkeys

	monkeyClient = client{
		target:      Target,
		outputTimer: speedReports,
	}
	// listen for speed reports
	go func() {
		for report := range speedReports {
			// HACK: ignore updates from monkeys that we don't know about yet
			if report.id < len(Bullpen) {
				Bullpen[report.id].speed = report.speed
			}
		}
	}()

	// send in the monkeys!
	for i := 0; i < seatCount; i++ {
		newMonkey := monkeyClient.createNew(fmt.Sprintf("Monkey%d", i), i)
		seats[i] = seat{ // (The sim starts with all seats filled with monkeys)
			layout: "qwerty",
			monkey: newMonkey,
		}
		newMonkey.seated = true
		go newMonkey.startTyping(i)
	}
}

// AddMonkey processes user requests to add another monkey
func AddMonkey() (*Monkey, error) {
	i := len(Bullpen)
	monkey := monkeyClient.createNew(fmt.Sprintf("Monkey%d", i), i)
	return monkey, nil
}

// AddSeat processes user requests to add another seat
func AddSeat(input AddSeatInput) (seat, error) {
	newSeat := seat{layout: input.Layout}
	seats[len(seats)] = newSeat
	return newSeat, nil
}

// Stand processes user requests to get a monkey out of its seat
func Stand(id int) (err error) {
	err = seats[id].monkey.stand()
	return
}

// Sit processes user requests to get a monkey out of its seat
func Sit(id int) (err error) {
	err = Bullpen[id].sit()
	return
}

// Answer is the minified monkey entry sent to the HTML template.
type Answer struct {
	Seat     int
	ID       int
	Name     string
	Speed    float64
	Progress string
	Keyboard string
	Seated   bool
}

// FetchResults turns the collection of monkey stats into a format
// that can be read by the HTML templates. This returns ONLY the monkeys
// that are seated and working.
func FetchResults() []Answer {
	results := []Answer{}

	for i := 0; i < len(seats); i++ { // so they show up in order
		monkey := seats[i].monkey
		if monkey == nil {
			results = append(results, Answer{Seat: i})
		} else {
			toAdd := Answer{
				Seat:     i,
				ID:       monkey.id,
				Name:     monkey.name,
				Speed:    monkey.speed,
				Progress: fmt.Sprintf("|%s|", Target[:monkey.highwater+1]),
				Keyboard: seats[i].layout,
				Seated:   monkey.seated,
			}
			results = append(results, toAdd)
		}

	}
	return results
}

// FetchAll turns the collection of monkey stats into a format
// that can be read by the HTML templates. This returns ALL monkeys.
func FetchAll() []Answer {
	results := []Answer{}

	for i := 0; i < len(Bullpen); i++ { // so they show up in order
		monkey := Bullpen[i]
		if monkey == nil {
			results = append(results, Answer{Seat: i})
		} else {
			toAdd := Answer{
				Seat:     i,
				ID:       monkey.id,
				Name:     monkey.name,
				Speed:    monkey.speed,
				Progress: fmt.Sprintf("|%s|", Target[:monkey.highwater+1]),
				Keyboard: seats[i].layout,
				Seated:   monkey.seated,
			}
			results = append(results, toAdd)
		}

	}
	return results
}
