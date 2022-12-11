package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func getRegExp () *regexp.Regexp {
	return regexp.MustCompile(`(?P<start1>\d*)-(?P<end1>\d*),(?P<start2>\d*)-(?P<end2>\d*)`)
}

func getValues(r *regexp.Regexp, s string) (int, int, int, int) {
	match := r.FindStringSubmatch(s)

	s1, _ := strconv.Atoi(match[r.SubexpIndex("start1")])
	s2, _ := strconv.Atoi(match[r.SubexpIndex("start2")])
	e1, _ := strconv.Atoi(match[r.SubexpIndex("end1")])
	e2, _ := strconv.Atoi(match[r.SubexpIndex("end2")])

	return s1, e1, s2, e2
}

func cleanup_first(scan *bufio.Scanner) int {
	r := getRegExp()
	var total int
	for scan.Scan() {
		start1, end1, start2, end2 := getValues(r, scan.Text())
		if (start1 >= start2 && end1 <= end2) || (start2 >= start1 && end2 <= end1) {
			//fmt.Printf("start1: %v end1: %v, start2: %v, end2: %v\n", start1, end1, start2, end2)
			total++
		}
	}
	return total
}

func cleanup_second(scan *bufio.Scanner) int {
	r := getRegExp()
	var total int
	for scan.Scan() {
		start1, end1, start2, end2 := getValues(r, scan.Text())
		if (start1 <= end2 && start2 <= end1) || (start2 <= end1 && start1 <= end2) {
			total++
		}
	}
	return total
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	fmt.Printf("First: %v\n", cleanup_first(scan))
	f.Close()

	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(bufio.NewReader(f))
	fmt.Printf("Second: %v\n", cleanup_second(scan))
	f.Close()
}
