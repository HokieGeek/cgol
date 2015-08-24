package cgol

const (
	STD_UNDERPOP       = 2
	STD_OVERPOPULATION = 3
	STD_LIVE           = 3
)

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

func getAllNeighbors(pond *Pond, organism OrganismReference) []OrganismReference {
	neighbors := getOrthogonalNeighbors(pond, organism)

	// Determine the offsets
	//above := organism.Y - 1
	//below := organism.Y + 1
	//left := organism.X - 1
	//right := organism.X + 1

	// TODO: Retrieve the rest of the neighbors

	return neighbors
}

func standard(pond *Pond, organism OrganismReference, getNeighbors func(*Pond, OrganismReference) []OrganismReference) {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it is dead
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	neighbors := getNeighbors(pond, organism)

	// Determine if current cell lives
	cell := pond.getOrganism(organism)
	// TODO: neighbors := getOrthogonalNeighbors(pond, organism)
	if *cell < 0 && *cell == STD_LIVE {
		pond.updateOrganismLivingState(organism, true)
		// TODO: for each neighbor, increase their count
		for neighbor := range neighbors {
			pond.updateNeighborCount(neighbors[neighbor], 1)
		}
	} else if *cell < STD_UNDERPOP || *cell > STD_OVERPOPULATION {
		pond.updateOrganismLivingState(organism, false)
		// TODO: for each neighbor, decrease their count
		for neighbor := range neighbors {
			pond.updateNeighborCount(neighbors[neighbor], -1)
		}
	}
}

func StandardOrthogonalNew(pond *Pond, organism OrganismReference) {
	standard(pond, organism, getOrthogonalNeighbors)
}

func StandardOrthogonal(pond *Pond, organism OrganismReference) {
	// -- Rules --
	// 1. If live cell has < 2 neighbors, it is dead
	// 2. If live cell has 2 or 3 neighbors, it lives
	// 3. If live cell has > 3 neighbors, it dies
	// 4. If dead cell has exactly 3 neighbors, it lives

	// Determine if current cell lives
	cell := pond.getOrganism(organism)
	// TODO: neighbors := getOrthogonalNeighbors(pond, organism)
	if *cell < 0 && *cell == STD_LIVE {
		pond.updateOrganismLivingState(organism, true)
		// TODO: for each neighbor, increase their count
	} else if *cell < STD_UNDERPOP || *cell > STD_OVERPOPULATION {
		pond.updateOrganismLivingState(organism, false)
		// TODO: for each neighbor, decrease their count
	}
}
