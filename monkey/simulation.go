package monkey

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var Seats []*Monkey
var Target string

func KickOffSim() {
	seatCount := getSeatCount()
	ClearScreen(seatCount)

	updates := make(chan Report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}       // how we know when all the monkeys are done
	Seats = []*Monkey{}
	Target = getTarget("target.txt")
	speedReports := make(chan SpeedReport, 500) // receiving speed reading from monkeys

	monkeyClient := Client{
		Target:      Target,
		Updates:     updates,
		Done:        toWait,
		OutputTimer: speedReports,
	}
	// listen for speed reports
	go func() {
		for report := range speedReports {
			// HACK: ignore updates from monkeys that we don't know about yet
			if report.ID < len(Seats) {
				Seats[report.ID].Speed = report.Speed
				Results(Seats, Target)
			}
		}
	}()

	// send in the monkeys!
	for i := 0; i < seatCount; i++ {
		newMonkey := Monkey{
			ID:        i,
			Name:      fmt.Sprintf("Monkey%d", i),
			Highwater: -1,
			Profile:   ConstructTypingProfile(),
		}
		Seats = append(Seats, &newMonkey)

		go newMonkey.StartTyping(Target, updates, toWait, speedReports)
		toWait.Add(1)
	}

	go closeChannelWhenDone(toWait, updates)
	// listen for updates
	go func() {
		for update := range updates {
			// HACK: ignore updates from monkeys that we don't know about yet
			if update.ID < len(Seats) {
				Seats[update.ID].Highwater = update.Highwater
				Results(Seats, Target)
			}
		}
	}()

	// keep an eye out for user input
	reader := bufio.NewReader(os.Stdin)
	var response string
	for {
		input := getInput(len(Seats), reader)
		Seats, response = processInput(input, Seats, monkeyClient)
		ClearScreen(len(Seats))
		Results(Seats, Target) // reprint table in case a monkey got renamed
		AtCursor(0, len(Seats)+7, response)
	}
}

func FetchResults() map[string]string {
	results := map[string]string{}

	for _, monkey := range Seats {
		results[monkey.Name] = fmt.Sprintf("|%s|", Target[:monkey.Highwater+1])
	}
	return results
}
