package main

import (
	"regexp"
	"strings"

	tm "github.com/buger/goterm"
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

func printResults(results []int, target string) {
	tm.MoveCursor(1, 2)
	for id, highwater := range results {
		tm.Print("Monkey ", id)
		tm.MoveCursor(20, id+2) // NOTE: If the header gets longer, this "2" needs to change.
		tm.Print("|", target[:highwater+1], "|")
		tm.MoveCursor(10, len(results)+4)
	}
	tm.Flush() // adds line break
}
