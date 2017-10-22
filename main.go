package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

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

func main() {
	rand.Seed(9)
	possibilities := "abcdefghijklmnopqrstuvqxyz    "

	file, err := ioutil.ReadFile("target.txt")
	if err != nil {
		fmt.Printf("\n\nERROR=======\n|%v|\n\n", err)
		os.Exit(1)
	}
	target := processTarget(file)

	currentPosition := 0
	highwater := -1

	for true {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentPosition] {
			if currentPosition > highwater {
				fmt.Print(target[:currentPosition+1])
				fmt.Printf("\nNEW HIGH POINT: was %v, now %v\n", highwater, currentPosition)
				highwater = currentPosition
			}
			currentPosition++
			continue
		}
		if currentPosition > 0 {
			currentPosition = 0
		}
	}
	// This section is unreachable when the for loop is infinite:
	// fmt.Printf("\n\nCOMPLETE!\nLongest achievement: ")
	// for p := 0; p <= highwater; p++ {
	// 	fmt.Print(string(target[p]))
	// }
	// fmt.Print("\n\n")
}
