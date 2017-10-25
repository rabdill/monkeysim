package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/rabdill/monkeysim/monkey"
	"github.com/rabdill/monkeysim/printer"
)

func main() {
	seatCount := getSeatCount()
	printer.ClearScreen(seatCount)

	updates := make(chan monkey.Report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}              // how we know when all the monkeys are done
	seats := []*monkey.Monkey{}
	target := getTarget("target.txt")
	speedReports := make(chan monkey.SpeedReport, 500) // receiving speed reading from monkeys

	monkeyClient := monkey.Client{
		Target:      target,
		Updates:     updates,
		Done:        toWait,
		OutputTimer: speedReports,
	}
	// listen for speed reports
	go func() {
		for report := range speedReports {
			// ignore updates from monkeys that we don't know about yet
			if report.ID < len(seats) {
				seats[report.ID].Speed = report.Speed
				printer.Results(seats, target)
			}
		}
	}()

	// send in the monkeys!
	for i := 0; i < seatCount; i++ {
		newMonkey := monkey.Monkey{
			ID:        i,
			Name:      fmt.Sprintf("Monkey%d", i),
			Highwater: -1,
			Profile:   monkey.ConstructTypingProfile(),
		}
		seats = append(seats, &newMonkey)

		go newMonkey.StartTyping(target, updates, toWait, speedReports)
		toWait.Add(1)
	}

	go closeChannelWhenDone(toWait, updates)
	// listen for updates
	go func() {
		for update := range updates {
			// ignore updates from monkeys that we don't know about yet
			if update.ID < len(seats) {
				seats[update.ID].Highwater = update.Highwater
				printer.Results(seats, target)
			}
		}
	}()

	// keep an eye out for user input
	reader := bufio.NewReader(os.Stdin)
	var response string
	for {
		seats, response = processInput(getInput(len(seats), reader), seats, monkeyClient)
		printer.ClearScreen(len(seats))
		printer.Results(seats, target) // reprint table in case a monkey got renamed
		printer.AtCursor(0, len(seats)+7, response)
	}
}
