package monkey

// Results - write a table of the monkeys' high scores to the terminal
func Results(results []*Monkey, target string) {
	// headerSize := 2 // NOTE: If the header gets longer, this "2" needs to change.

	// for i, monkey := range results {
	// 	AtCursor(0, i+headerSize, monkey.Name) // TODO: why write the monkeys' name over and over?
	// 	AtCursor(15, i+headerSize, fmt.Sprintf("%.2f kpms", monkey.Speed))
	// 	AtCursor(35, i+headerSize, fmt.Sprintf("|%s|", target[:monkey.Highwater+1]))
	// 	// go back to soliciting user input once we're done printing:
	// 	AtCursor(0, len(results)+4, "Enter command: ")
	// }
}

// MoveCursor - shift terminal printer to particular coordinate
func MoveCursor(x, y int) {
	// fmt.Printf("\033[%d;%dH", y, x)
}

// AtCursor - shift terminal printer to particular coordinate and print something
func AtCursor(x, y int, toPrint string) {
	// fmt.Printf("\033[%d;%dH%s", y, x, toPrint)
}

// ClearScreen - wipe out most stuff that's going to be in the way
func ClearScreen(seatCount int) {
	// for i := 0; i < 50; i++ {
	// 	AtCursor(0, i, ClearingString())
	// }
	// AtCursor(0, 0, "MONKEYSIM")
	// AtCursor(0, seatCount+4, "Enter command: ")
}

// ClearingString - a long blank string for zeroing out the display
func ClearingString() string {
	return "                                                                                                                  "
}
