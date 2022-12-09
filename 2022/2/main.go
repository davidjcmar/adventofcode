package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var cipher = map[string]int {
		"A": 1,
		"B": 2,
		"C": 3,
		"X": 1,
		"Y": 2,
		"Z": 3,
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	scan.Split(bufio.ScanWords)
	var score int
	for scan.Scan() {
		opp := cipher[scan.Text()]
		scan.Scan()
		me := cipher[scan.Text()]
		switch outcome := opp - me; {
			// tie
			case outcome == 0:
				score += 3
				score += me
				fmt.Printf("Draw opp: %v  me: %v  outcome: %v\n", opp, me, outcome)
			// win
			case outcome == -1, outcome == 2:
				score += 6
				score += me
				fmt.Printf("Win opp: %v  me: %v  outcome: %v\n", opp, me, outcome)
			// loss
			case outcome == 1, outcome == -2:
				score += me
				fmt.Printf("Loss opp: %v  me: %v  outcome: %v\n", opp, me, outcome)
		}
	}
	fmt.Println(score)
}
