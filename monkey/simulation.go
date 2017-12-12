package monkey

import (
	"fmt"
	"sync"
)

// Seats holds the current roster of monkeys and their stats.
var Seats []*monkey

//Target holds the goal text that the monkeys are working toward.
var Target string

// KickOffSim starts the simulation and sets the first set of
// monkeys off to typing.
func KickOffSim() {
	seatCount := getSeatCount()

	updates := make(chan report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}       // how we know when all the monkeys are done
	Seats = []*monkey{}
	Target = getTarget("target.txt")
	speedReports := make(chan speedReport, 500) // receiving speed reading from monkeys

	// monkeyClient := client{
	// 	target:      Target,
	// 	updates:     updates,
	// 	done:        toWait,
	// 	outputTimer: speedReports,
	// }
	// listen for speed reports
	go func() {
		for report := range speedReports {
			// HACK: ignore updates from monkeys that we don't know about yet
			// *tk when would this happen?!
			if report.id < len(Seats) {
				Seats[report.id].speed = report.speed
				// *update screen
			}
		}
	}()

	// send in the monkeys!
	for i := 0; i < seatCount; i++ {
		newMonkey := monkey{
			id:        i,
			name:      fmt.Sprintf("Monkey%d", i),
			highwater: -1,
			profile:   constructTypingProfile(),
		}
		Seats = append(Seats, &newMonkey)

		go newMonkey.startTyping(Target, updates, toWait, speedReports)
		toWait.Add(1)
	}

	go closeChannelWhenDone(toWait, updates)
	// listen for updates
	go func() {
		for update := range updates {
			// HACK: ignore updates from monkeys that we don't know about yet
			if update.id < len(Seats) {
				Seats[update.id].highwater = update.highwater
				// *update screen
			}
		}
	}()
}

// Answer is the minified monkey entry sent to the HTML template.
type Answer struct {
	Name     string
	Speed    float64
	Progress string
}

// FetchResults turns the collection of monkey stats into a format
// that can be read by the HTML templates.
func FetchResults() []Answer {

	results := []Answer{}

	for _, monkey := range Seats {
		results = append(results, Answer{monkey.name, monkey.speed, fmt.Sprintf("|%s|", Target[:monkey.highwater+1])})
	}
	return results
}
