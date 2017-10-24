package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/rabdill/monkeysim/monkey"
	"github.com/rabdill/monkeysim/printer"
)

func main() {
	printer.ClearScreen()
	seatCount := getSeatCount()

	printer.AtCursor(0, 0, "MONKEYSIM")
	printer.AtCursor(0, seatCount+4, "Enter command: ")

	updates := make(chan monkey.Report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}              // how we know when all the monkeys are done
	seats := []monkey.Monkey{}

	// read the target file
	file, err := ioutil.ReadFile("target.txt")
	if err != nil {
		fmt.Printf("\n\nERROR reading file: |%v|\n\n", err)
		os.Exit(1)
	}
	target := processTarget(file)

	// send in the monkeys!
	for i := 0; i < seatCount; i++ {
		go monkey.StartTyping(i, target, updates, toWait)
		toWait.Add(1)
		seats = append(seats, monkey.Monkey{
			Name:      fmt.Sprintf("Monkey %d", i),
			Highwater: -1,
		})
	}
	// once all the monkeys say they're done, close the updates channel
	go func() {
		toWait.Wait()
		close(updates)
	}()

	// listen for updates
	go func() {
		for update := range updates {
			seats[update.ID].Highwater = update.Highwater
			printer.Results(seats, target)
		}
	}()

	// keep an eye out for user input
	reader := bufio.NewReader(os.Stdin)
	for {
		input := getInput(seatCount, reader)
		printer.AtCursor(0, seatCount+7, printer.ClearingString())
		printer.AtCursor(0, seatCount+7, fmt.Sprintf("YOU ENTERED %s", input))
		printer.AtCursor(20, seatCount+4, printer.ClearingString())
		processInput(input, seatCount)
	}
}
