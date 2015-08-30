package cgol

//////////////////// STANDARD RULESET ///////////////////

// RulesConwayLife is a function that applies the "standard" rules for Conway's Game of Life
// Returns true if the organism lives
//
// -- Rules --
// 1. If live cell has < 2 neighbors, it dies
// 2. If live cell has 2 or 3 neighbors, it lives
// 3. If live cell has > 3 neighbors, it dies
// 4. If dead cell has exactly 3 neighbors, it lives
func RulesConwayLife(numNeighbors int, isAlive bool) bool {
	const (
		STD_UNDERPOPULATION = 2
		STD_OVERCROWDING    = 3
		STD_REVIVE          = 3
	)

	// Rule #4
	if !isAlive && numNeighbors == STD_REVIVE {
		return true

	} else if numNeighbors >= 0 &&
		(numNeighbors < STD_UNDERPOPULATION || numNeighbors > STD_OVERCROWDING) {
		// Rules #1 and #3
		return false
	}

	// Rule #2
	return isAlive
}
