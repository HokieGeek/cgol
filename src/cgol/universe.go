package cgol

type Universe struct {
	rows      int
	cols      int
	gameboard [][]int
	ruleset   func()
}

// func (t* Universe) create(rows int, cols int, rules func()) {
// }
