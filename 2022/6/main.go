package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func findSignal(scan *bufio.Scanner) (int, error) {
	scan.Scan()
	in := scan.Text()

	for i:=4; i<len(in); i++ {
		if len(map[byte]byte{in[i]: in[i], in[i-1]:in[i-1], in[i-2]:in[i-2], in[i-3]:in[i-3]}) == 4 {
			return i+1, nil
		}
	}
	return 0, errors.New("No marker found")
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	marker, err := findSignal(scan)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("First: %v\n", marker)
}
