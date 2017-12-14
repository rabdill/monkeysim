package monkey

import (
	"fmt"
	"sync"
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
	updates := make(chan report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}       // how we know when all the monkeys are done
	Bullpen = []*Monkey{}
	Target = getTarget("target.txt")
	speedReports := make(chan speedReport, 500) // receiving speed reading from monkeys

	monkeyClient = client{
		target:      Target,
		updates:     updates,
		done:        toWait,
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
			keyboard: "standard",
			monkey:   newMonkey,
		}
	}

	go closeChannelWhenDone(toWait, updates)
	// listen for updates
	go func() {
		for update := range updates {
			// HACK: ignore updates from monkeys that we don't know about yet
			if update.id < len(Bullpen) {
				Bullpen[update.id].highwater = update.highwater
				// *update screen
			}
		}
	}()
}

// AddMonkey processes user requests to add more monkeys
func AddMonkey() (*Monkey, error) {
	i := len(Bullpen)
	monkey := monkeyClient.createNew(fmt.Sprintf("Monkey%d", i), i)
	return monkey, nil
}

// StandUp processes user requests to get a monkey out of its seat
func StandUp(id int) (err error) {
	err = Bullpen[id].standUp()
	return
}

// Answer is the minified monkey entry sent to the HTML template.
type Answer struct {
	Seat     int
	ID       int
	Name     string
	Speed    float64
	Progress string
	Seated   bool
}

// FetchResults turns the collection of monkey stats into a format
// that can be read by the HTML templates.
func FetchResults() []Answer {
	results := []Answer{}

	for i := 0; i < len(seats); i++ { // so they show up in order
		monkey := seats[i].monkey
		if monkey == nil {
			results = append(results, Answer{Seat: i})
		} else {
			results = append(results, Answer{i, monkey.id, monkey.name, monkey.speed, fmt.Sprintf("|%s|", Target[:monkey.highwater+1]), monkey.seated})
		}

	}
	return results
}

// FetchAll turns the collection of monkey stats into a format
// that can be read by the HTML templates.
func FetchAll() []Answer {
	results := []Answer{}

	for i := 0; i < len(Bullpen); i++ { // so they show up in order
		monkey := Bullpen[i]
		if monkey == nil {
			results = append(results, Answer{Seat: i})
		} else {
			results = append(results, Answer{i, monkey.id, monkey.name, monkey.speed, fmt.Sprintf("|%s|", Target[:monkey.highwater+1]), monkey.seated})
		}

	}
	return results
}
