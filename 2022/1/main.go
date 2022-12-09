package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	scan := bufio.NewScanner(r)
	var t string
	var elves []int
	i := 0
	for scan.Scan() {
		t = scan.Text()
		if i == len(elves) {
			elves = append(elves, 0)
		}
		if t != "" {
			c, _ := strconv.Atoi(t)
			elves[i] += c
		} else {
			i++
		}
	}
	e := 0
	for j, c := range elves {
		if j==0 || e < c {
			e = c
		}
	}
	fmt.Printf("%v\n", e)
}
