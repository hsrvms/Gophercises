package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// parse command line flags
	filename := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	timeLimit := flag.Int("limit", 10, "the time limit for the quiz in seconds")
	flag.Parse()

	// Open the csv file
	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file %s", *filename))
	}

	// Create a new csv reader
	r := csv.NewReader(file)

	// Read all records from the csv file
	lines, err := r.ReadAll()
	if err != nil {
		exit("Error reading CSV file")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Keep track of the number of correct and incorrect answers
	correct := 0

	// Loop through each record and ask the question
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	// Output the number of correct answers and total number of questions
	fmt.Printf("Correct %d out of %d\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for idx, line := range lines {
		problems[idx] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func exit(msg string) {
	log.Println(msg)
	os.Exit(1)
}
