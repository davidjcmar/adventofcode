package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func priority(r int) int {
	if r >= 97 && r <= 122 {
		return r -96
	} else {
		return r -38
	}
}

func rucksackOne(scan *bufio.Scanner) int {
	var priority_total int
	for scan.Scan() {
		contents := scan.Text()
		l := len(contents)
		first_comp := contents[:l/2]
		second_comp := contents[l/2:]
		for _, r := range(first_comp) {
			if strings.Contains(second_comp, string(r)) {
				//fmt.Printf ("first: %v, second: %v, r: %c, p: %v\n", first_comp, second_comp, rune(r), priority(int(r)))
				priority_total += priority(int(r))
				break
			}
		}
	}
	//fmt.Printf("a: %v, z: %v\n", int('a'), int('z'))
	//fmt.Printf("A: %v, Z: %v\n", int('A'), int('Z'))
	return priority_total
}

func rucksackTwo(scan *bufio.Scanner) int {
	var priority_total int
	for scan.Scan() {
		ruck_one := scan.Text()
		scan.Scan()
		ruck_two := scan.Text()
		scan.Scan()
		ruck_three := scan.Text()

		for _, r := range(ruck_one) {
			if strings.Contains(ruck_two, string(r)) && strings.Contains(ruck_three, string(r)) {
				priority_total += priority(int(r))
				break
			}
		}
	}
	return priority_total
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	f2, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	scan_two := bufio.NewScanner(bufio.NewReader(f2))

	fmt.Printf("Output one: %v\n", rucksackOne(scan))
	fmt.Printf("Output two: %v\n", rucksackTwo(scan_two))

	f.Close()
}
