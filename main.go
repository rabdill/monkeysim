package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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

// monkey - frantically typing random characters
func monkey(id int, done chan bool) {
	fmt.Printf("\n=======NEW MONKEY %d!========\n", id)
	rand.Seed(int64(id))
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	file, err := ioutil.ReadFile("target.txt")
	if err != nil {
		fmt.Printf("\n\nERROR=======\n|%v|\n\n", err)
		os.Exit(1)
	}
	target := processTarget(file)

	currentPosition := 0
	highwater := -1

	for {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentPosition] {
			if currentPosition > highwater {
				fmt.Printf("\n%d - NEW HIGH POINT: was %v, now %v: |%v|", id, highwater, currentPosition, target[:currentPosition+1])
				highwater = currentPosition
			}
			currentPosition++
			continue
		}
		// if we were on a streak, but it's over
		if currentPosition > 0 {
			if highwater > 5 && highwater == currentPosition { // if we were close
				fmt.Printf("\n%d - just missed: %s%s", id, target[:currentPosition], string(keyPress))
			}
			currentPosition = 0
		}
	}
}

func main() {
	monkeyCount := 8

	done := make(chan bool, monkeyCount)

	for i := 0; i < monkeyCount; i++ {
		go monkey(i, done)
	}
	<-done // note: this will never actually happen
}
