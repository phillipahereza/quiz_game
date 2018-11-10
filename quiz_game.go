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
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format `question,answer`")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file")
	}

	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	for i, v := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, v.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYour time has elapsed. You scored %d out of %d\n", correct, len(problems))
			return

		case answer := <-answerCh:
			if answer == v.answer {
				correct++
			}
		}
	}

	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{line[0], strings.TrimSpace(line[1])}

	}

	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(-1)

}
