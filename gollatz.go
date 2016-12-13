package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// Check if the user actually provided an argument
	if len(os.Args) != 2 {
		fmt.Printf("Usage: collatz [unsigned 64-bit integer]\n")
	} else {
		// Parse the command line argument
		max, err := strconv.ParseUint(os.Args[1], 10, 64)
		bufferSize := 1000

		// Make sure the user wasn't stupid.
		if checkErr(err) {
			fmt.Printf("You dun goofed. You probably put in a non-integer value.")
		} else {
			// This is so we can time the program. We only want to time the
			// main logic.
			start := time.Now()

			// This channel will store the results from the individual
			// goroutines, and will be processed by the scoreboard function
			results := make(chan Result, bufferSize)

			// This channel will increment with each value of n. This way,
			// we can dynamically spin up goroutines whenever they are
			// available to handle the next value of n, and we should be
			// able to achieve 100% processor usage.
			valueN := make(chan uint64)

			valueN <- max

			// This Result is going to store our highest score
			var overallHighScore Result

			// This is going to hold our next value of n so that we can avoid initializing it over and over again every time we spin up a new goroutine
			var nextValue uint64
			nextValue = iterateN(valueN)
			threaded(nextValue, results)
			go func() {
				for true {
					nextValue = iterateN(valueN)
					go threaded(nextValue, results)
					fmt.Printf("The buffer is currently %d / %d full", len(results), bufferSize)
				}
			}()

			// This is so that we can constantly update the overall high score, since we need to keep it running.
			go func() {
				for true {
					overallHighScore = scoreboard(results, overallHighScore)
				}
			}()

			trackTime("The Collatz portion", start)
			fmt.Printf("%d has takes the most steps at %d.\n", overallHighScore.highestValue, overallHighScore.highestScore)
		} // End of main program logic

	} // End of program
}

// Result will store the highest scores and its corresponding value of n
type Result struct {
	highestScore uint
	highestValue uint64
}

// Classic 3n+1 conjecture.
func collatz(num uint64) Result {
	var count uint
	var result Result
	result.highestValue = num
	for num > 1 {
		if num%2 == 0 {
			num /= 2
		} else {
			num = num*3 + 1
		}
		count++
	}
	result.highestScore = count
	return result
}

// This function takes a value, and then runs thecollatz function on it.
// Then it sticks the result of that function into the channel that holds
// our results
func threaded(value uint64, results chan Result) {

	var result = collatz(value)

	results <- result
}

// This function keeps track of the actual highest score, and then reports
// it as it is updated.
func scoreboard(results chan Result, highScore Result) Result {
	nextScore := <-results
	if nextScore.highestScore > highScore.highestScore {
		highScore = nextScore
		fmt.Printf("%d currently takes the most steps at %d\n\n", highScore.highestValue, highScore.highestScore)
	}
	return highScore
}

// This function returns the current value of the channel, and increments
// it and updates the channel with thenew value
func iterateN(valueN chan uint64) uint64 {
	currentN := <-valueN
	valueN <- currentN + 1
	return currentN
}

// This function tracks time. Yay!
func trackTime(name string, start time.Time) {
	fmt.Printf("%s took %s to execute.\n", name, time.Since(start))
}

// Simply returns false if there is no error, and true if there is one.
// Technically, this is error handling!
func checkErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
