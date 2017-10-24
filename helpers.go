package main

import (
	"fmt"
	"regexp"
	"strings"
)

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

// TODO: why write the monkeys' name over and over?
func printResults(results []int, target string) {
	headerSize := 2 // NOTE: If the header gets longer, this "2" needs to change.
	fmt.Print("\033[0;2H")
	for id, highwater := range results {
		PrintAtCursor(0, id+headerSize, fmt.Sprintf("Monkey %d", id))
		PrintAtCursor(20, id+headerSize, fmt.Sprintf("|%s|", target[:highwater+1]))
		// go back to soliciting user input once we're done printing:
		MoveCursor(20, len(results)+4)
	}
}

// MoveCursor - shift terminal printer to particular coordinate
func MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// PrintAtCursor - shift terminal printer to particular coordinate and print something
func PrintAtCursor(x, y int, toPrint string) {
	fmt.Printf("\033[%d;%dH%s", y, x, toPrint)
}

func ClearScreen() {
	for i := 0; i < 50; i++ {
		PrintAtCursor(0, i, "                                                             ")
	}
}
