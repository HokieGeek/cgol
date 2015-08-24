package cgol

type GameStatus int

const (
	Active GameStatus = 1
	Stable GameStatus = 2
	Dead   GameStatus = 3
)

func (t GameStatus) String() string {
	return "TODO"
}

type Organism struct {
	x int
	y int
	// z int
}

type Universe struct {
	rows      int
	cols      int
	gameboard [][]int
	numLiving int
	ruleset   func()
	status    GameStatus
	name      string
	// living    []int // Would be best if this was a tuple :-/...
}

// func (t* Universe) applyRuleset() {
// }

// func (t* Universe) create(rows int, cols int, rules func(), init func()) {
// }
