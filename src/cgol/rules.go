package cgol

func getOrthogonalNeighbors(pond *Pond, organism OrganismReference) []OrganismReference {
	neighbors := make([]OrganismReference, 4)

	// Determine the offsets
	above := organism.Y - 1
	below := organism.Y + 1
	left := organism.X - 1
	right := organism.X + 1

	if above >= 0 {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: above})
	}

	if below <= pond.Rows {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: below})
	}

	if left >= 0 {
		neighbors = append(neighbors, OrganismReference{X: left, Y: organism.Y})
	}

	if right <= pond.Cols {
		neighbors = append(neighbors, OrganismReference{X: right, Y: organism.Y})
	}

	return neighbors
}

func getObliqueNeighbors(pond *Pond, organism OrganismReference) []OrganismReference {
	neighbors := make([]OrganismReference, 4)

	// Determine the offsets
	above := organism.Y - 1
	below := organism.Y + 1
	left := organism.X - 1
	right := organism.X + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, OrganismReference{X: left, Y: above})
		}
		if right <= pond.Cols {
			neighbors = append(neighbors, OrganismReference{X: right, Y: above})
		}
	}

	if below <= pond.Rows {
		if left >= 0 {
			neighbors = append(neighbors, OrganismReference{X: left, Y: below})
		}
		if right <= pond.Cols {
			neighbors = append(neighbors, OrganismReference{X: right, Y: below})
		}
	}

	return neighbors
}

func getAllNeighbors(pond *Pond, organism OrganismReference) []OrganismReference {
	neighbors := append(getOrthogonalNeighbors(pond, organism),
		getObliqueNeighbors(pond, organism)...)
	// neighbors := make([]OrganismReference, 8)
	// neighbors = append(neighbors, getOrthogonalNeighbors(pond, organism))
	// neighbors = append(neighbors, getObliqueNeighbors(pond, organism))

	return neighbors
}

//////////////////// STANDARD RULESET ///////////////////

const (
	STD_UNDERPOPULATION = 2
	STD_OVERCROWDING    = 3
	STD_REVIVE          = 3
)

func standard(pond *Pond, organism OrganismReference,
	getNeighbors func(*Pond, OrganismReference) []OrganismReference) {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it is dead
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	// TODO: does this logic properly handle when an organism dies?

	// Determine if current cell lives
	neighbors := getNeighbors(pond, organism)
	neighborCount := pond.getNeighborCount(organism)
	if neighborCount < 0 {
		// Rule #4
		numLivingNeighbors := 0
		for _, neighbor := range neighbors {
			if pond.getNeighborCount(neighbor) >= 0 {
				numLivingNeighbors++
			}
		}
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

func StandardOrthogonal(pond *Pond, organism OrganismReference) {
	standard(pond, organism, getOrthogonalNeighbors)
}

func StandardOblique(pond *Pond, organism OrganismReference) {
	standard(pond, organism, getObliqueNeighbors)
}

func StandardAll(pond *Pond, organism OrganismReference) {
	standard(pond, organism, getAllNeighbors)
}
