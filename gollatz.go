package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// This map will store all of the numbers and how many steps it takes to complete their Collatz run.
	var scores map[uint64]uint
	scores = make(map[uint64]uint)

	// Check if the user actually provided an argument
	if len(os.Args) != 2 {
		fmt.Printf("Usage: collatz [unsigned 64-bit integer]\n")
	} else {
		// Parse the command line argument
		max, err := strconv.ParseUint(os.Args[1], 10, 64)

		// Make sure the user wasn't stupid.
		if checkErr(err) {
			fmt.Printf("You dun goofed. You probably put in a non-integer value.")
		} else {
			// This is so we can time the program. We only want to time the main logic.
			// Using defer will execute the function once all the logic is finished.
			// However, this also tracks the Printf output.
			// For now, trackTime will return a string instead of printing.
			//defer trackTime(time.Now())
			start := time.Now()

			// This is the main loop of the program.
			var value uint64
			for value = 1; value < max; value++ {
				scores[value] = collatz(value)
			} // End of main loop
			trackTime("The Collatz portion", start)

			// This is where we figure out which one took the most steps.
			var highestScore uint
			var highestIndex uint64

			// Reuse value because, hey, free memory.
			for value = 0; value < max; value++ {
				if scores[value] > highestScore {
					highestScore = scores[value]
					highestIndex = value
				}
			}
			trackTime("The program", start)
			fmt.Printf("%d has takes the most steps at %d.\n", highestIndex, highestScore)
		} // End of program

	}
}

// Classic 3n+1 conjecture.
func collatz(num uint64) uint {
	var count uint
	for num > 1 {
		if num%2 == 0 {
			num /= 2
		} else {
			num = num*3 + 1
		}
		count++
	}
	return count
}

// This function tracks time. Yay!
func trackTime(name string, start time.Time) {
	fmt.Printf("%s took %s to execute.\n", name, time.Since(start))
}

// Simply returns false if there is no error, and true if there is one. Technically, this is error handling!
func checkErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
