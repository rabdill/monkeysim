package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	tm "github.com/buger/goterm"
	"github.com/rabdill/monkeysim/monkey"
)

func main() {
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
	tm.Clear() // Clear current screen
	tm.MoveCursor(1, 1)
	tm.Print("MONKEYSIM")
	tm.MoveCursor(1, monkeyCount+4)
	tm.Print("Enter command: ")
	tm.Flush() // adds line break

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
			printResults(highwater, target)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		input := getInput(monkeyCount, reader)
		tm.MoveCursor(1, monkeyCount+7)
		tm.Print("YOU ENTERED ", input)
	}
}

func getInput(monkeyCount int, reader *bufio.Reader) string {
	tm.MoveCursor(10, monkeyCount+4)
	text, _ := reader.ReadString('\n')
	return strings.TrimRight(text, "\n")
}
