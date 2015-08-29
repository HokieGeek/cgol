package cgol

//////////////////// STANDARD RULESET ///////////////////
func Standard(numNeighbors int, isAlive bool) bool {
	// fmt.Printf("Standard(%d)\n", numNeighbors)
	// TODO: figure out go documentation
	// Returns true if the organism lives

	// -- Rules --
	// 1. If live cell has < 2 neighbors, it dies
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	// TODO: would be nice if these could be mucked with
	const (
		STD_UNDERPOPULATION = 2
		STD_OVERCROWDING    = 3
		STD_REVIVE          = 3
	)

	// Rule #4
	if !isAlive && numNeighbors == STD_REVIVE {
		// fmt.Printf("Reviving: %s\n", organism.String())
		return true

	} else if numNeighbors >= 0 &&
		(numNeighbors < STD_UNDERPOPULATION || numNeighbors > STD_OVERCROWDING) {
		// Rules #1 and #3
		// fmt.Printf("Killing: %s\n", organism.String())
		return false
	}

	return isAlive
}
