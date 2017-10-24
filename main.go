package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/rabdill/monkeysim/monkey"
	"github.com/rabdill/monkeysim/printer"
)

func main() {
	printer.ClearScreen()
	var monkeyCount int
	var err error
	if len(os.Args) > 1 {
		monkeyCount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("\nFATAL: MonkeyCount parameter could not be converted to int: %v", err)
			os.Exit(1)
		}
	} else {
		monkeyCount = 1
	}
	printer.AtCursor(0, 0, "MONKEYSIM")
	printer.AtCursor(0, monkeyCount+4, "Enter command: ")

	updates := make(chan monkey.Report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}              // how we know when all the monkeys are done
	highwater := []int{}                     // best each monkey's done

	// read the target file
	file, err := ioutil.ReadFile("target.txt")
	if err != nil {
		fmt.Printf("\n\nERROR reading file: |%v|\n\n", err)
		os.Exit(1)
	}
	target := processTarget(file)

	// send in the monkeys!
	for i := 0; i < monkeyCount; i++ {
		go monkey.StartTyping(i, target, updates, toWait)
		toWait.Add(1)
		highwater = append(highwater, -1) // so everybody's entry is in the right spot
	}
	// once all the monkeys say they're done, close the updates channel
	go func() {
		toWait.Wait()
		close(updates)
	}()

	go func() {
		// listen for updates
		for update := range updates {
			highwater[update.Id] = update.Highwater
			printer.Results(highwater, target)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		input := getInput(monkeyCount, reader)
		printer.AtCursor(0, monkeyCount+7, printer.ClearingString())
		printer.AtCursor(0, monkeyCount+7, fmt.Sprintf("YOU ENTERED %s", input))
		printer.AtCursor(20, monkeyCount+4, printer.ClearingString())
		processInput(input, monkeyCount)
	}
}
