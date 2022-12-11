package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func findSignal(input string, count int) (int, error) {
	for i:=count; i<len(input); i++ {
		m := make(map[byte]byte)
		for j:=0; j<count; j++ {
			m[input[i-j]] = input[i-j]
		}
		if len(m) == count {
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
	scan.Scan()
	input := scan.Text()
	p_marker, m_marker := 4, 14
	packet_marker, err := findSignal(input, p_marker)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("First: %v\n", packet_marker)

	message_marker, err := findSignal(input, m_marker)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Second: %v\n", message_marker)
}
