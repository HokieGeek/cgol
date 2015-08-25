package cgol

//////////////////// STANDARD RULESET ///////////////////
const (
	STD_UNDERPOPULATION = 2
	STD_OVERCROWDING    = 3
	STD_REVIVE          = 3
)

func Standard(pond *Pond, organism OrganismReference) bool {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it dies
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	modified := false
	neighborCount := pond.getNeighborCount(organism)

	// Test rules
	if neighborCount < 0 {
		// Rule #4
		numLivingNeighbors := pond.calculateNeighborCount(organism)
		if numLivingNeighbors == STD_REVIVE {
			pond.setNeighborCount(organism, numLivingNeighbors)
			modified = true
		}

	} else if neighborCount < STD_UNDERPOPULATION || neighborCount > STD_OVERCROWDING {
		// Rules #1 and #3
		pond.setNeighborCount(organism, -1)
		modified = true
	}

	return modified
}
