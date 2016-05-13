package life

import (
	"bytes"
	"strconv"
)

type Rules struct {
	Survive []int // The number of neighbors an alive cell needs to have to survive
	Born    []int // The number of neighbors a dead cell needs to have to be born
}

func (t *Rules) String() string {
	var buf bytes.Buffer

	for _, val := range t.Survive {
		buf.WriteString(strconv.Itoa(val))
	}

	buf.WriteString("/")

	for _, val := range t.Born {
		buf.WriteString(strconv.Itoa(val))
	}

	return buf.String()
}

// This function tests the number of neighbors that a cell has against
// the rules given. If the cell lives (or continues to live) based on the
// rules, then the function returns true.
func testRule(numNeighbors int, isAlive bool, rules *Rules) bool {
	list := rules.Survive
	if !isAlive {
		list = rules.Born
	}

	for _, val := range list {
		if val == numNeighbors {
			return true
		}
	}

	return false
}

// Returns a RulesTest function that uses the given ruleset
func RulesTester(rules *Rules) func(int, bool) bool {
	return func(numNeighbors int, isAlive bool) bool {
		return testRule(numNeighbors, isAlive, rules)
	}
}

// Returns a Rules struct filled with the normal Conway rules of 23/3
//	-- Rules --
// 	1. If live cell has < 2 neighbors, it dies
// 	2. If live cell has 2 or 3 neighbors, it lives
// 	3. If live cell has > 3 neighbors, it dies
// 	4. If dead cell has exactly 3 neighbors, it lives
func GetConwayRules() *Rules {
	return &Rules{Survive: []int{2, 3}, Born: []int{3}}
}

// Returns a rules tester with Conway Normal rules
func ConwayTester() func(int, bool) bool {
	return RulesTester(GetConwayRules())
}
