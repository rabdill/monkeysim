package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

func processInput(input string, seatCount int) {
	switch input {
	case "exit":
		printer.MoveCursor(0, seatCount+9)
		os.Exit(0)
	}
}

func getInput(seatCount int, reader *bufio.Reader) string {
	printer.MoveCursor(20, seatCount+4)
	text, _ := reader.ReadString('\n')
	return strings.TrimRight(text, "\n")
}
