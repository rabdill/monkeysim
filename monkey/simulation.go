package monkey

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func KickOffSim() {
	seatCount := getSeatCount()
	ClearScreen(seatCount)

	updates := make(chan Report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}       // how we know when all the monkeys are done
	seats := []*Monkey{}
	target := getTarget("target.txt")
	speedReports := make(chan SpeedReport, 500) // receiving speed reading from monkeys

	monkeyClient := Client{
		Target:      target,
		Updates:     updates,
		Done:        toWait,
		OutputTimer: speedReports,
	}
	// listen for speed reports
	go func() {
		for report := range speedReports {
			// HACK: ignore updates from monkeys that we don't know about yet
			if report.ID < len(seats) {
				seats[report.ID].Speed = report.Speed
				Results(seats, target)
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
		seats = append(seats, &newMonkey)

		go newMonkey.StartTyping(target, updates, toWait, speedReports)
		toWait.Add(1)
	}

	go closeChannelWhenDone(toWait, updates)
	// listen for updates
	go func() {
		for update := range updates {
			// HACK: ignore updates from monkeys that we don't know about yet
			if update.ID < len(seats) {
				seats[update.ID].Highwater = update.Highwater
				Results(seats, target)
			}
		}
	}()

	// keep an eye out for user input
	reader := bufio.NewReader(os.Stdin)
	var response string
	for {
		input := getInput(len(seats), reader)
		seats, response = processInput(input, seats, monkeyClient)
		ClearScreen(len(seats))
		Results(seats, target) // reprint table in case a monkey got renamed
		AtCursor(0, len(seats)+7, response)
	}
}
