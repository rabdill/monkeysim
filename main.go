package main

import (
	"fmt"
	"math/rand"
)

func main() {
	rand.Seed(9)
	possibilities := "abcdefghij"
	target := "abcdefghij"

	currentPosition := 0
	highwater := -1

	for i := 0; i < 99999999; i++ {
		keyPress := possibilities[rand.Intn(len(possibilities))]
		if keyPress == target[currentPosition] {
			if currentPosition > highwater {
				fmt.Print(keyPress)
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
	fmt.Printf("\n\nCOMPLETE!\nLongest achievement: ")
	for p := 0; p <= highwater; p++ {
		fmt.Print(string(target[p]))
	}
	fmt.Print("\n\n")
}
