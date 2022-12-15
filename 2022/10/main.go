package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	addx = 2
)

func (s signal) calc() int {
	return s.i * s.v
}
type signal struct {
	i, v int
}

func newSprite(i int) *sprite {
	s := sprite{x: i-1, y: i, z: i+1}
	return &s
}
func (s *sprite) move(i int) {
	s.y = i
	s.x = s.y - 1
	s.z = s.y + 1
}
type sprite struct {
	x, y, z int
}

func check_interesting(c, x int, s []signal) []signal {
	if (c - 20) % 40 == 0 {
		//fmt.Printf("Interesting! x: %d cycle: %d v: %d\n", x, c, x*c)
		i := signal{i: c, v: x}
		s = append(s, i)
	}
	return s
}

func sum_cycles(scan *bufio.Scanner) int {
	var sum int
	cycle, x := 1, 1
	var signals []signal

	for scan.Scan() {
		signals = check_interesting(cycle, x, signals)
		t := scan.Text()
		if t == "noop" {
			cycle++
		} else if t == "addx" {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			cycle++
			signals = check_interesting(cycle, x, signals)
			cycle++
			x += v
			//fmt.Printf("x == %d, cycle %d\n", x, cycle)
		}
	}
	//fmt.Printf("signals: %v\n", signals)
	for _, v := range(signals) {
		sum += v.calc()
	}
	return sum
}

func check_pixel(p, w int, s *sprite) bool {
	px := p % w
	if s.x == px || s.y == px || s.z == px {
		return true
	}
	return false
}

func draw(scan *bufio.Scanner) {
	var crt []bool
	p := 0
	x := 1
	w := 40
	s_pos := newSprite(x)
	for scan.Scan() {
		crt = append(crt, false)
		crt[p] = check_pixel(p, w, s_pos)
		t := scan.Text()
		if t == "noop" {
			p++
			continue
		} else if t == "addx" {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			crt = append(crt, false)
			p++
			crt[p] = check_pixel(p, w, s_pos)
			p++
			x += v
			s_pos.move(x)
		}
	}
	for i, v := range(crt) {
		if i % w == 0 {
			fmt.Println()
		}
		if v {
			fmt.Printf("#")
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Println()
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	scan.Split(bufio.ScanWords)

	first := sum_cycles(scan)
	fmt.Printf("First: %d\n", first)
	f.Close()

	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(bufio.NewReader(f))
	scan.Split(bufio.ScanWords)

	draw(scan)
	
}