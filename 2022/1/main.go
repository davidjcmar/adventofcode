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
	one := 0
	two := 0
	three := 0
	for j, c := range elves {
		if j==0 || one < c {
			three = two
			two = one
			one = c
		} else if two < c {
			three = two
			two = c
		} else if three < c {
			three = c
		}
	}
	total := one + two + three
	fmt.Printf("%v %v %v\n", one, two, three)
	fmt.Printf("%v\n", total)
}
