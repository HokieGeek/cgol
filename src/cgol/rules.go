package cgol

//////////////////// STANDARD RULESET ///////////////////

const (
	STD_UNDERPOPULATION = 2
	STD_OVERCROWDING    = 3
	STD_REVIVE          = 3
)

func Standard(pond *Pond, organism OrganismReference) {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it is dead
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	// TODO: does this logic properly handle when an organism dies?

	// Determine if current cell lives
	neighbors := pond.GetNeighbors(organism)
	neighborCount := pond.getNeighborCount(organism)
	if neighborCount < 0 {
		// Rule #4
		// numLivingNeighbors := 0
		// for _, neighbor := range neighbors {
		// 	if pond.getNeighborCount(neighbor) >= 0 {
		// 		numLivingNeighbors++
		// 	}
		// }
		numLivingNeighbors := pond.calculateNeighborCount(organism)
		if numLivingNeighbors == STD_REVIVE {
			pond.setNeighborCount(organism, numLivingNeighbors)
			for _, neighbor := range neighbors {
				pond.incrementNeighborCount(neighbor)
			}
		}

	} else if neighborCount < STD_UNDERPOPULATION || neighborCount > STD_OVERCROWDING {
		// Rules #1 and #3
		pond.setNeighborCount(organism, -1)
		for _, neighbor := range neighbors {
			pond.decrementNeighborCount(neighbor)
		}
	}
	// TODO: Rule #2?
}
