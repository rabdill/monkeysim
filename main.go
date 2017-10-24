package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	tm "github.com/buger/goterm"
)

// report - how monkeys check in with the master process
type report struct {
	id, highwater int
}

// processTarget - Turns file contents into a string containing only a-z characters
func processTarget(input []byte) (output string) {
	output = string(input)
	output = strings.ToLower(output)

	// cut out line breaks:
	output = strings.Replace(output, "\n", " ", -1)
	// make sure we don't have long stretches of spaces
	re := regexp.MustCompile(" +")
	output = re.ReplaceAllLiteralString(output, " ")
	// make sure we only have alphabet characters now:
	re = regexp.MustCompile("[^a-z ]")
	output = re.ReplaceAllLiteralString(output, "")
	return
}

// monkey - frantically typing random characters
func monkey(id int, target string, updates chan report, done *sync.WaitGroup) {
	defer done.Done()
	rand.Seed(int64(id))
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	currentSearch := 0
	highwater := -1

	for i := 0; i < 99999999; i++ {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentSearch] {
			if currentSearch > highwater {
				highwater = currentSearch
				updates <- report{id, highwater}
			}
			currentSearch++
			continue
		}
		// if we were on a streak, but it's over
		if currentSearch > 0 {
			currentSearch = 0
		}
	}
}

func printResults(results []int, target string) {
	tm.MoveCursor(1, 2)
	for id, highwater := range results {
		tm.Print("Monkey ", id)
		tm.MoveCursor(20, id+2) // NOTE: If the header gets longer, this "2" needs to change.
		tm.Print("|", target[:highwater+1], "|")
		tm.Flush() // adds line break
	}

}

func main() {
	tm.Clear() // Clear current screen
	tm.MoveCursor(1, 1)
	tm.Print("MONKEYSIM")

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

	updates := make(chan report, 100) // how monkeys check in with us
	toWait := &sync.WaitGroup{}       // how we know when all the monkeys are done
	highwater := []int{}              // best each monkey's done

	// read the target file
	file, err := ioutil.ReadFile("target.txt")
	if err != nil {
		fmt.Printf("\n\nERROR reading file: |%v|\n\n", err)
		os.Exit(1)
	}
	target := processTarget(file)

	// send in the monkeys!
	for i := 0; i < monkeyCount; i++ {
		go monkey(i, target, updates, toWait)
		toWait.Add(1)
		highwater = append(highwater, -1) // so everybody's entry is in the right spot
	}
	// once all the monkeys say they're done, close the updates channel
	go func() {
		toWait.Wait()
		close(updates)
	}()
	// listen for updates
	for update := range updates {
		highwater[update.id] = update.highwater
		printResults(highwater, target)
	}
}
