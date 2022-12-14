package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func (t *tail) print() {
	fmt.Println("Record")
	for p, m := range(t.record) {
		fmt.Printf("%v -- %d\n", p, m)
	}

	dedup := make(map[position]int)
	fmt.Println("History")
	for _, p := range(t.history) {
		dedup[p]++
		fmt.Printf("%v\n", p)
	}
	fmt.Println("Diff")
	fmt.Printf("record: %d -- history: %d\n", len(t.record), len(dedup))

}

func (t *tail) follow(p position) {
	var dx, dy, abx, aby int
	dx = p.x - t.pos.x
	dy = p.y - t.pos.y
	//fmt.Printf("dx: %d dy: %d -- h_x: %d h_y: %d ", dx, dy, p.x, p.y)
	if dx < 0 {
		abx = -dx
	} else {
		abx = dx
	}
	if dy < 0 {
		aby = -dy
	} else {
		aby = dy
	}

	if (abx > 1 && aby >= 1) || (abx >= 1 && aby > 1) {
		if dx > 0 {
			t.pos.x++
		} else {
			t.pos.x--
		}
		if dy > 0 {
			t.pos.y++
		} else {
			t.pos.y--
		}
	} else if abx > 1 {
		if dx > 0 {
			t.pos.x++
		} else {
			t.pos.x--
		}
	} else if aby > 1 {
		if dy > 0 {
			t.pos.y++
		} else {
			t.pos.y--
		}
	}
	t.record[t.pos]++
	t.history = append(t.history, t.pos)
	//fmt.Printf("-- tail x: %d y: %d\n", t.pos.x, t.pos.y)
}

func (t *tail) position_count() int {
	return len(t.record)
}

func (t *tail) move(direction string) {
	switch direction {
	case "U":
		t.pos.y += 1
	case "D":
		t.pos.y -= 1
	case "L":
		t.pos.x -= 1
	case "R":
		t.pos.x += 1
	}
}

func newTail() *tail {
	t := tail{}
	t.record = make(map[position]int)
	return &t
}

func newHead() *head {
	h := head{}
	return &h
}

func (h *head) move(direction string) {
	switch direction {
	case "U":
		h.pos.y += 1
	case "D":
		h.pos.y -= 1
	case "L":
		h.pos.x -= 1
	case "R":
		h.pos.x += 1
	}
}

type position struct {
	x, y int
}

type head struct {
	pos position
}

type tail struct {
	pos    position
	record map[position]int
	history []position
}

func chase(scan *bufio.Scanner) int {
	h := head{}
	t := newTail()

	for scan.Scan() {
		direction := scan.Text()
		scan.Scan()
		count, _ := strconv.Atoi(scan.Text())
		//fmt.Printf("dir: %v -- count: %v\n", direction, count)

		for i:=0; i < count; i++ {
			h.move(direction)
			t.follow(h.pos)
		}
	}
	//t.print()
	return t.position_count()
}

func broken(scan *bufio.Scanner, n int) int {
	var t []*tail

	for i:=0; i<n; i++ {
		t = append(t, newTail())
	}

	for scan.Scan() {
		direction := scan.Text()
		scan.Scan()
		count, _ := strconv.Atoi(scan.Text())

		//fmt.Printf("direction: %s -- count: %d\n", direction, count)
		for i:=0; i<count; i++ {
			t[0].move(direction)
			for j:=1; j<n; j++ {
				t[j].follow(t[j-1].pos)
			}
		}
	}
	return t[n-1].position_count()
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	scan.Split(bufio.ScanWords)

	first := chase(scan)
	f.Close()
	fmt.Printf("First: %d\n", first)

	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(bufio.NewReader(f))
	scan.Split(bufio.ScanWords)
	second := broken(scan, 10)
	f.Close()
	fmt.Printf("Second: %d\n", second)
}
