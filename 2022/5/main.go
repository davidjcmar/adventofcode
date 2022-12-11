package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func initStacks(scan *bufio.Scanner) []string {
	var stacks []string
	for scan.Scan(){
		var current_stack int
		line := scan.Text()
		if strings.Contains(line, "1") {
			break
		}

		for i:=0; i<len(line); {
			var crate bool
			for j:=0; j<4; j++ {
				r, b := utf8.DecodeRune([]byte(line[i:]))
				switch (j) {
				case 0:
					// check if we have a crate on this stack
					if r == '[' {
						crate = true
					}
				case 1:
					// add crate to data if exists
					if crate {
						for len(stacks) <= current_stack {
							stacks = append(stacks, "")
						}
						stacks[current_stack] = stacks[current_stack] + string(r)
					}

				case 2:
					// noop
				case 3:
					// reset crate bool for next stack
					crate = false
					current_stack++
				}
				// don't run off end of string
				i += b
				if i == len(line) {
					break
				}
			}
		}
	}
	scan.Scan() // eat blank line
	return stacks
}

func initInstructions(scan *bufio.Scanner) []string {
	var instructions []string
	for scan.Scan() {
		instructions = append(instructions, scan.Text())
	}
	return instructions
}

func executeInstructions(version int, instructions []string, stacks []string) []string {
	for _, instruction := range(instructions) {
		stacks = crane(version, instruction, stacks)
	}
	return stacks
}

func getTopCrates(stacks []string) string {
	var top string
	for _, v := range(stacks) {
		r, _ := utf8.DecodeRune([]byte(v))
		top = top + string(r)
	}
	return top
}

func crane(version int, instruction string, stacks []string) []string {
	r := regexp.MustCompile(`move (?P<count>\d*) from (?P<from>\d*) to (?P<to>\d*)`)
	match := r.FindStringSubmatch(instruction)

	count,_ := strconv.Atoi(match[r.SubexpIndex("count")])
	from, _ := strconv.Atoi(match[r.SubexpIndex("from")])
	to, _ := strconv.Atoi(match[r.SubexpIndex("to")])
	// subtract 1 to account for zero index in slice
	from -= 1
	to -= 1

	// avoid accessing off end of stack string
	if count > len(stacks[from]) {
		count = len(stacks[from])
	}
	switch (version) {
	case 9000:
		// move crates one at a time to destination stack
		for i:=0; i<count; i++ {
			stacks[to] = string(stacks[from][i]) + stacks[to]
		}
	case 9001:
		// move multiple crates at a time
		stacks[to] = stacks[from][:count] + stacks[to]
	}
	// cleanup source stack
	if count == len(stacks[from]) {
		stacks[from] = ""
	} else {
		stacks[from] = stacks[from][count:]
	}
	return stacks
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	stacks := initStacks(scan)
	instructions := initInstructions(scan)
	f.Close()
	executeInstructions(9000, instructions, stacks)
	fmt.Printf("First: %v\n", getTopCrates(stacks))
	
	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(bufio.NewReader(f))
	stacks = initStacks(scan)
	executeInstructions(9001, instructions, stacks)
	fmt.Printf("Second: %v\n", getTopCrates(stacks))
}
