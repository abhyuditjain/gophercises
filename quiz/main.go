package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a CSV file name in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz (in seconds)")

	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}
	_ = file

	r := csv.NewReader(file)

	lines, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("Failed to read CSV file with error: %v", err))
	}

	problems := parseProblems(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	score := 0
	for i, p := range problems {
		// Print problem
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		ansChan := make(chan string)
		go func() {
			var ans string
			// Scanf blocks so if we don't enter anything, timer won't affect
			// Hence we put it in a goroutine
			_, _ = fmt.Scanf("%s\n", &ans)
			ansChan <- ans
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYour score = %d/%d\n", score, len(problems))
			return
		case ans := <-ansChan:
			if ans == p.a {
				score++
			}
		}
	}

	fmt.Printf("Your score = %d/%d\n", score, len(problems))
}

func parseProblems(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, v := range lines {
		problems[i] = problem{
			q: v[0],
			a: strings.TrimSpace(v[1]),
		}
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
