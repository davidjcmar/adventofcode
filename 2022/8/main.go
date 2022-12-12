package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	up = iota
	down
	left
	right
)

func newTree(h int, v bool) *tree {
	t := tree{height: h, visible: v}
	return &t
}

func (t *tree) setTreeVis() {
	t.visible = true
}

type tree struct {
	height int
	visible bool
} 

func initTreeMap(scan *bufio.Scanner) [][]*tree {
	var treeMap [][]*tree
	for scan.Scan() {
		text := scan.Text()
		row := make([]*tree, len(text))
		for i, v := range(text) {
			h, _ := strconv.Atoi(string(v))
			row[i] = newTree(h, false)
		}
		treeMap = append(treeMap, row)
	}
	return treeMap
}

func checkDirection(direction, x, y int, treeMap[][]*tree, c chan<- bool) {
	visible := true
	switch(direction) {
	case up:
		for i:=x-1; i>=0; i-=1 {
			if treeMap[i][y].height >= treeMap[x][y].height {
				visible = false
				break
			}
		}

	case down:
		for i:=x+1; i<len(treeMap); i+=1 {
			if treeMap[i][y].height >= treeMap[x][y].height {
				visible = false
				break
			}
		}
	
	case left:
		for i:=y-1; i>=0; i-=1 {
			if treeMap[x][i].height >= treeMap[x][y].height {
				visible = false
				break
			}
		}
	
	case right:
		for i:=y+1; i<len(treeMap[x]); i+=1 {
			if treeMap[x][i].height >= treeMap[x][y].height {
				visible = false
				break
			}
		}
	}

	c <- visible
}

func checkViewingDistance(direction, x, y int, treeMap [][]*tree, d chan<- int) {
	distance := 0
	switch(direction) {
	case up:
		for i:=x; i>0; i-= 1 {
			if treeMap[i][y].height >= treeMap[x][y].height && !(i == x){
				break
			}
			distance++
		}
	case down:
		for i:=x; i<len(treeMap)-1; i+=1 {
			if treeMap[i][y].height >= treeMap[x][y].height && !(i == x){
				break
			}
			distance++
		}

	case left:
		for i:=y; i>0; i-=1 {
			if treeMap[x][i].height >= treeMap[x][y].height && !(i == y){
				break
			}
			distance++
		}

	case right:
		for i:=y; i<len(treeMap[x])-1; i+=1 {
			if treeMap[x][i].height >= treeMap[x][y].height && !(i == y){
				break
			}
			distance++
		}
	}
	d <- distance
}

func setTreeMapVis(treeMap [][]*tree) (int, [][]*tree) {
	var visCount int
	for i, _ := range treeMap {
		for j, _ := range treeMap[i] {
			c := make(chan bool)

			if i == 0 || j == 0 || i == len(treeMap) || j == len(treeMap[i]) {
				treeMap[i][j].setTreeVis()
				visCount++
				close(c)
			} else {
				go checkDirection(up, i, j, treeMap, c)
				go checkDirection(down, i, j, treeMap, c)
				go checkDirection(left, i, j, treeMap, c)
				go checkDirection(right, i, j, treeMap, c)

				visible := make(map[bool]bool)
				count := 0
				for {
					v := <-c
					count++
					visible[v] = v
					if count == 4 {

						close(c)
						break
					}
				}
				if visible[true] {
					visCount++
				}
			}
		}
	}
	return visCount, treeMap
}

func getMaxTreeVisDistance(treeMap [][]*tree) int {
	var maxDistance int
	for i, _ := range treeMap {
		for j, _ := range treeMap[i] {
			d := make(chan int)

			go checkViewingDistance(up, i, j, treeMap, d)
			go checkViewingDistance(down, i, j, treeMap, d)
			go checkViewingDistance(left, i, j, treeMap, d)
			go checkViewingDistance(right, i, j, treeMap, d)

			count := 0
			distanceScore := 1
			for {
				distance := <-d
				count++
				distanceScore = distanceScore * distance
				if count == 4 {
					close(d)
					//fmt.Printf ("i: %d  j: %d  tree: %d\n", i, j, treeMap[i][j].height)
					if maxDistance == 0 || distanceScore > maxDistance {
						//fmt.Printf("maxD: %d  newMax: %d\n", maxDistance, distanceScore)
						maxDistance = distanceScore
					}
					break
				}
			}
		}
	}
	return maxDistance
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	var visibleCount int
	treeMap := initTreeMap(bufio.NewScanner(bufio.NewReader(f)))
	visibleCount, treeMap = setTreeMapVis(treeMap)
	/* print tree map */
	/*
	for i, _ := range(treeMap) {
		for j, _ := range(treeMap[i]) {
			fmt.Printf("%d ", treeMap[i][j].height)
		}
		fmt.Println()
	}
	*/
	fmt.Printf("First: %d\n", visibleCount)

	fmt.Printf("Second: %d\n", getMaxTreeVisDistance(treeMap))
}