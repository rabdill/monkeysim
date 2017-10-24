package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/rabdill/monkeysim/monkey"
	"github.com/rabdill/monkeysim/printer"
)

func getSeatCount() (seatCount int) {
	var err error
	if len(os.Args) > 1 {
		seatCount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("\nFATAL: seatCount parameter could not be converted to int: %v", err)
			os.Exit(1)
		}
	} else {
		seatCount = 1
	}
	return
}

func closeChannelWhenDone(waitgroup *sync.WaitGroup, channel chan monkey.Report) {
	waitgroup.Wait()
	close(channel)
}

// processTarget - Turns file contents into a string containing only a-z characters
func getTarget(file string) (output string) {
	// read the target file
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("\n\nERROR reading file: |%v|\n\n", err)
		os.Exit(1)
	}
	output = string(contents)
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

func processInput(input string, seats []monkey.Monkey) ([]monkey.Monkey, string) {
	command := strings.Split(input, " ")
	switch command[0] {
	case "exit":
		printer.MoveCursor(0, len(seats)+9)
		os.Exit(0)
	case "rename":
		index := findMonkeyInList(seats, command[1])
		if index < 0 {
			return seats, "ERROR: Could not find monkey by that name to rename."
		}
		seats[index].Name = command[2]
		return seats, fmt.Sprintf("Renamed %s to %s.", command[1], command[2])
	}
	return seats, fmt.Sprintf("Unrecognized command: %s", input)
}

func findMonkeyInList(haystack []monkey.Monkey, needle string) int {
	for i, monkey := range haystack {
		if monkey.Name == needle {
			return i
		}
	}
	return -1
}

func getInput(seatCount int, reader *bufio.Reader) string {
	printer.MoveCursor(20, seatCount+4)
	text, _ := reader.ReadString('\n')
	printer.AtCursor(0, seatCount+7, printer.ClearingString())  // clear the last command's status message
	printer.AtCursor(20, seatCount+4, printer.ClearingString()) // clear the command line for the next command
	return strings.TrimRight(text, "\n")
}
