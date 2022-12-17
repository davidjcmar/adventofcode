package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	monkeyIt = iota
	itemsIt
	operationIt
	testIt
	conditionIt
)

const (
	addition = iota
	subtraction
	multiplication
	division
)

func regexByLineType(lt string) *regexp.Regexp {
	var regex *regexp.Regexp
	switch lt {
	case "Monkey":
		regex = regexp.MustCompile(`Monkey (?P<monkeyId>\d+):`)
	case "Starting":
		regex = regexp.MustCompile(`\s*Starting items: (?P<itemList>.*)`)
	case "Operation":
		regex = regexp.MustCompile(`\s*Operation: new = (?P<operation>.*)`)
	case "Test":
		regex = regexp.MustCompile(`\s*Test: divisible by (?P<test>\d*)`)
	case "If":
		regex = regexp.MustCompile(`\s*If (?P<condition>true|false): throw to monkey (?P<target>\d*)`)
	}
	return regex
}

func parseLine(line string) (int, []int) {
	var rType int
	var rList []int

	lineTypeRegEx := regexp.MustCompile(`\s*(?P<type>\w*):{0,1}\s.*`)

	matchType := lineTypeRegEx.FindStringSubmatch(line)

	lineType := matchType[lineTypeRegEx.SubexpIndex("type")]
	regex := regexByLineType(lineType)
	if regex == nil {
		fmt.Println("Problem with regexByLineType")
		os.Exit(1)
	}
	match := regex.FindStringSubmatch(line)

	switch lineType {
	case "Monkey":
		rType = monkeyIt
		mId, _ := strconv.Atoi(match[regex.SubexpIndex("monkeyId")])
		rList = append(rList, mId)
	case "Starting":
		rType = itemsIt
		iStringList := strings.Split(match[regex.SubexpIndex("itemList")], ", ")
		for _, v := range iStringList {
			itemInt, _ := strconv.Atoi(v)
			rList = append(rList, itemInt)
		}
	case "Operation":
		rType = operationIt
		oStringList := strings.Split(match[regex.SubexpIndex("operation")], " ")
		for _, v := range oStringList {
			switch v {
			case "old":
				rList = append(rList, 0)
			case "+":
				rList = append(rList, addition)
			case "*":
				rList = append(rList, multiplication)
			case "-":
				rList = append(rList, subtraction)
			case "/":
				rList = append(rList, division)
			default:
				oInt, _ := strconv.Atoi(v)
				rList = append(rList, oInt)
			}
		}
	case "Test":
		rType = testIt
		tInt, _ := strconv.Atoi(match[regex.SubexpIndex("test")])
		rList = append(rList, tInt)
	case "If":
		rType = conditionIt
		condString := match[regex.SubexpIndex("condition")]
		targetInt, _ := strconv.Atoi(match[regex.SubexpIndex("target")])
		switch condString {
		case "true":
			rList = append(rList, 1)
		case "false":
			rList = append(rList, 0)
		}
		rList = append(rList, targetInt)
	}

	return rType, rList
}

func newOperation(operation []int, tdp int) func(int) int {
	operand2 := operation[2]
	op := operation[1]

	if op == addition {
		if operand2 != 0 {
			return func(w int) int {
				if tdp != 0 {
					return (w + operand2) % tdp
				} else {
					return w + operand2
				}
			}
		} else {
			return func(w int) int {
				if tdp != 0 {
					return (w * w) % tdp
				} else {
					return w + w
				}
			}
		}
	} else if op == multiplication {
		if operand2 != 0 {
			return func(w int) int {
				if tdp != 0 {
					return (w * operand2) % tdp
				} else {
					return w * operand2
				}
			}
		} else {
			return func(w int) int {
				if tdp != 0 {
					return (w * w) % tdp
				} else {
					return w * w
				}
			}
		}
	} else if op == subtraction {
		if operand2 != 0 {
			return func(w int) int {
				return w - operand2
			}
		} else {
			return func(w int) int {
				return 0
			}
		}
	} else if op == division {
		if operand2 != 0 {
			return func(w int) int {
				return w / operand2
			}
		} else {
			return func(w int) int {
				return 1
			}
		}
	}
	return nil
}

func newMonkeyBusiness(operation []int, test, cond1, cond0, tdp int) monkeyBusiness {
	var target int

	op := newOperation(operation, tdp)

	return func(w, wMod int) (int, int) {
		worry := op(w)
		if wMod > 0 {
			worry = worry / wMod
		} else {

		}
		//fmt.Printf("%d\n", worry)
		if worry%test == 0 {
			target = cond1
		} else {
			target = cond0
		}
		return worry, target
	}
}

type monkeyBusiness func(int, int) (int, int)

func newMonkey(mId, test, cond1, cond0 int, operation, items []int, tdp int) *monkey {
	mb := newMonkeyBusiness(operation, test, cond1, cond0, tdp)
	m := monkey{id: mId, mb: mb, items: items, inspectCount: 0, testDivisor: test}
	return &m
}

type monkey struct {
	id, inspectCount, testDivisor int
	items                         []int
	mb                            monkeyBusiness
}

func checkMonkeyDone(m map[string]bool) bool {
	var ok bool
	if _, ok = m["monkeyId"]; !ok {
		return false
	} else if _, ok = m["items"]; !ok {
		return false
	} else if _, ok = m["operation"]; !ok {
		return false
	} else if _, ok = m["test"]; !ok {
		return false
	} else if _, ok = m["cond1"]; !ok {
		return false
	} else if _, ok = m["cond0"]; !ok {
		return false
	}
	return true
}

func generateMonkeys(scan *bufio.Scanner, tdp int) []*monkey {
	var monkeys []*monkey
	var mId, test, cond1, cond0 int
	var items, op []int

	for scan.Scan() {
		if scan.Text() == "" {
			continue
		}
		monkeyDone := false
		monkeyDoneMap := make(map[string]bool)
		for !monkeyDone {
			//fmt.Printf("Line: %v\n", scan.Text())
			l, v := parseLine(scan.Text())
			switch l {
			case monkeyIt:
				mId = v[0]
				monkeyDoneMap["monkeyId"] = true
			case itemsIt:
				items = v
				monkeyDoneMap["items"] = true
			case operationIt:
				op = v
				monkeyDoneMap["operation"] = true
			case testIt:
				test = v[0]
				monkeyDoneMap["test"] = true
			case conditionIt:
				if v[0] == 0 {
					cond0 = v[1]
					monkeyDoneMap["cond0"] = true
				} else {
					cond1 = v[1]
					monkeyDoneMap["cond1"] = true
				}
			}
			if checkMonkeyDone(monkeyDoneMap) {
				//fmt.Printf("Making monkey: %d\n", mId)
				monkeys = append(monkeys, newMonkey(mId, test, cond1, cond0, op, items, tdp))
				monkeyDone = true
				break
			}
			scan.Scan()
		}
	}

	return monkeys
}

func getTestDivisorsProduct(monkeys []*monkey) int {
	divisorProduct := 1
	for _, v := range monkeys {
		divisorProduct *= v.testDivisor
	}
	return divisorProduct
}

func keepAway(monkeys []*monkey, iterations, worryMod int) []*monkey {
	//fmt.Println("Enter keepaway")
	for i := 0; i < iterations; i++ {
		for _, monkey := range monkeys {
			//fmt.Printf("Monkey %v Items %v\n", monkey, monkey.items)
			for _, item := range monkey.items {
				w, t := monkey.mb(item, worryMod)
				//fmt.Printf("Monkey %d had item with worry %d and tossed it to monkey %d now has worry %d\n", monkey.id, item, t, w)
				monkey.inspectCount++
				monkey.items = monkey.items[1:]
				monkeys[t].items = append(monkeys[t].items, w)
			}
		}
	}
	return monkeys
}

func findMostActive(monkeys []*monkey) int {
	var activitySorted []int
	for _, v := range monkeys {
		activitySorted = append(activitySorted, v.inspectCount)
	}
	sort.Sort(sort.IntSlice(activitySorted))
	mostActive := activitySorted[len(activitySorted)-1] * activitySorted[len(activitySorted)-2]
	return mostActive
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	m := generateMonkeys(scan, 0)
	tdp := getTestDivisorsProduct(m)
	f.Close()
	//fmt.Printf("Monkeys:\n%v\n", m)
	m = keepAway(m, 20, 3)
	/*
		for _, monkey := range m {
			fmt.Printf("Monkey %d inspected %d items\n", monkey.id, monkey.inspectCount)
		}
	*/
	fmt.Printf("First: %d\n", findMostActive(m))

	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(bufio.NewReader(f))
	m = generateMonkeys(scan, tdp)
	f.Close()

	m = keepAway(m, 10000, 0)
	/*
		for _, monkey := range m {
			fmt.Printf("Monkey %d inspected %d items\n", monkey.id, monkey.inspectCount)
		}
	*/
	fmt.Printf("Second: %d\n", findMostActive(m))
}
