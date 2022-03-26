package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFlagPtr := flag.String("csv", "problems.csv", "csv file to use")
	limitFlagPtr := flag.Int("limit", 30, "time for the quiz")
	flag.Parse()

	questions := readQuestions(csvFlagPtr)

	fmt.Println("Press Enter when ready")
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(*limitFlagPtr) * time.Second)

	stats := askQuestions(questions, timer)
	fmt.Printf("%d out of %d were correct!\n", stats, len(questions))
}

func askQuestions(questions [][]string, timer *time.Timer) int {
	correctAnswers := 0
	for _, line := range questions {
		var response string
		ch := make(chan int)
		fmt.Printf("%s = ", line[0])
		go func() {
			fmt.Scanln(&response)
			ch <- 1
		}()
		select {
		case <-timer.C:
			fmt.Println("\nTIME RAN OUT!")
			return correctAnswers
		case <-ch:
			if response == line[1] {
				correctAnswers++
			}
		}
	}
	return correctAnswers
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readQuestions(csvFlagPtr *string) [][]string {
	f, err := os.Open(*csvFlagPtr)
	check(err)
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	check(err)
	return data
}
