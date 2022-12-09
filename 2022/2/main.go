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
	var score, second_score int
	for scan.Scan() {
		opp := cipher[scan.Text()]
		scan.Scan()
		me := cipher[scan.Text()]
		// first half
		switch outcome := opp - me; {
			// draw
			case outcome == 0:
				score += 3
				score += me
				// fmt.Printf("Draw\t\tOpp: %v\t\tMe: %v\n", opp, me)
			// win
			case outcome == -1, outcome == 2:
				score += 6
				score += me
			// loss
			case outcome == 1, outcome == -2:
				score += me
		}
		// second half
		switch (me) {
			// loss
			case 1:
				x := (opp + 2) % 3
				if x == 0 {
					x = 3
				}
				second_score += x
				fmt.Printf("Loss opp: %v  me: %v\n", opp, me)
			// draw
			case 2:
			    second_score += opp
				second_score += 3
				fmt.Printf("Draw opp: %v  me: %v\n", opp, me)
			// win
			case 3:
				x := (opp + 1) % 3
				if x == 0 {
					x = 3
				}
				second_score += x
				second_score += 6
				fmt.Printf("Win opp: %v  me: %v\n", opp, me)
		}
	}
	fmt.Println(score)
	fmt.Printf("Second score: %v\n", second_score)
}
