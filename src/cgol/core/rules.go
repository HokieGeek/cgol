package cgol

import (
	"bytes"
	"strconv"
)

type Rules struct {
	Survive []int
	Born    []int
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

func NewRulesFromString(rules string) *Rules {
	// TODO: split this: ####/#### (0 > # <= 8)
	return nil
}

// This function tests the number of neighbors that a cell has against
// the rules given. If the cell lives (or continues to live) based on the
// rules, then the function returns true.
func RulesTest(numNeighbors int, isAlive bool, rules *Rules) bool {
	// TODO: would it be weird if this function returned the rules?
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

func GetRulesTester(rules *Rules) func(int, bool) bool {
	return func(numNeighbors int, isAlive bool) bool {
		return RulesTest(numNeighbors, isAlive, rules)
	}
}

// Returns a Rules struct filled with the normal Conway rules of 23/3
func GetConwayRules() *Rules {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it dies
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives
	return &Rules{Survive: []int{2, 3}, Born: []int{3}}
}

func GetConwayTester() func(int, bool) bool {
	return GetRulesTester(GetConwayRules())
}
