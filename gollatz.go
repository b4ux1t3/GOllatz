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
			// This is so we can time the program. We only want to time the main logic.
			start := time.Now()

			// This channel will store the results from the individual goroutines,
			// and will be processed by the scoreboard function
			results := make(chan Result, bufferSize)

			// This channel will generate work for our threads
			work := make(chan uint64, max)

			// This Result is going to store our highest score
			var overallHighScore Result

			// This is where we need to start our goroutines
			go threaded(work, results)
			// go threaded(work, results)
			// go threaded(work, results)
			// go threaded(work, results)
			go scoreboard(results, overallHighScore)
			trackTime("The Collatz portion", start)

			fmt.Printf("%d takes the most steps at %d.\n", overallHighScore.Value, overallHighScore.Score)
		} // End of main program logic

	} // End of program
}

// Result will store the scores and its corresponding value of n
type Result struct {
	Score uint
	Value uint64
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

// This function will handle the main logic of the program. It will take the next value out of
// the input channel, and calculate collatz from that.
func threaded(in <-chan uint64, results chan<- Result) {
	// This is where we figure out which one took the most steps.
	value := <-in
	var result Result

	result.Score = collatz(value)
	result.Value = value

	results <- result

	fmt.Printf("%d, % d", result.Score, result.Value)

	close(results)
}

// This function keeps track of the actual highest score, and then reports it as it is updated.
func scoreboard(results <-chan Result, highScore Result) {
	var nextScore = <-results
	if nextScore.Score > highScore.Score {
		highScore = nextScore
		fmt.Printf("%d currently takes the most steps at %d\n\n", highScore.Value, highScore.Score)
	}
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
