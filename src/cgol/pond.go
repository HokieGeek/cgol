package cgol

import "fmt"

type GameStatus int

const (
	Active GameStatus = 1
	Stable GameStatus = 2
	Dead   GameStatus = 3
)

func (t GameStatus) String() string {
	s := ""

	if t&Active == Active {
		s += "Active"
	} else if t&Stable == Stable {
		s += "Stable"
	} else if t&Dead == Dead {
		s += "Dead"
	}

	return s
}

type OrganismReference struct {
	X int
	Y int
	Z int
}

// type PondStats struct {
// }

type Pond struct {
	Name         string
	Rows         int
	Cols         int
	NumLiving    int
	Status       GameStatus
	gameboard    [][]int
	processQueue []OrganismReference
	ruleset      func(*Pond, OrganismReference)
	initializer  func(*Pond)
}

func (t *Pond) getOrthogonalNeighbors(organism OrganismReference) []OrganismReference {
	neighbors := make([]OrganismReference, 4)

	// Determine the offsets
	above := organism.Y - 1
	below := organism.Y + 1
	left := organism.X - 1
	right := organism.X + 1

	if above >= 0 {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: above})
	}

	if below <= t.Rows {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: below})
	}

	if left >= 0 {
		neighbors = append(neighbors, OrganismReference{X: left, Y: organism.Y})
	}

	if right <= t.Cols {
		neighbors = append(neighbors, OrganismReference{X: right, Y: organism.Y})
	}

	return neighbors
}

func (t *Pond) getObliqueNeighbors(organism OrganismReference) []OrganismReference {
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
		if right <= t.Cols {
			neighbors = append(neighbors, OrganismReference{X: right, Y: above})
		}
	}

	if below <= t.Rows {
		if left >= 0 {
			neighbors = append(neighbors, OrganismReference{X: left, Y: below})
		}
		if right <= t.Cols {
			neighbors = append(neighbors, OrganismReference{X: right, Y: below})
		}
	}

	return neighbors
}

func (t *Pond) getAllNeighbors(organism OrganismReference) []OrganismReference {
	neighbors := append(t.getOrthogonalNeighbors(organism), t.getObliqueNeighbors(organism)...)
	// neighbors := make([]OrganismReference, 8)
	// neighbors = append(neighbors, getOrthogonalNeighbors(pond, organism))
	// neighbors = append(neighbors, getObliqueNeighbors(pond, organism))

	return neighbors
}

func (t *Pond) getNeighborCount(organism OrganismReference) int {
	return t.gameboard[organism.X][organism.Y]
}

func (t *Pond) setNeighborCount(organism OrganismReference, numNeighbors int) {
	// TODO: Mutex protection?
	originalNumNeighbors := t.gameboard[organism.X][organism.Y]

	t.gameboard[organism.X][organism.Y] = numNeighbors

	// Update the living count if organism changed living state
	if originalNumNeighbors < 0 && numNeighbors >= 0 {
		t.NumLiving++
		t.processQueue = append(t.processQueue, organism)
	} else if originalNumNeighbors >= 0 && numNeighbors < 0 {
		t.NumLiving--
		t.processQueue = append(t.processQueue, organism)
	}
}

func (t *Pond) incrementNeighborCount(organism OrganismReference) {
	t.setNeighborCount(organism, t.getNeighborCount(organism)+1)
}

func (t *Pond) decrementNeighborCount(organism OrganismReference) {
	t.setNeighborCount(organism, t.getNeighborCount(organism)-1)
}

// func (t *Pond) cycleProcessQueue() {
// 	front := t.processQueue[0]
// 	t.processQueue = append(t.processQueue[:0], t.processQueue[1:]...)
// 	t.ruleset(t, front)
// }

func (t *Pond) start() {
	// TODO: Kick off thread that periodically goes through the 'living' queue
}

func (t *Pond) stop() {
	// TODO: stop the processing thread
}

func (t *Pond) Display() {
	fmt.Printf("[%s] (%dx%d)\n", t.Name, t.Rows, t.Cols)
	fmt.Printf("Living organisms: %d\tStatus: %s\n", t.NumLiving, t.Status)
	for i := 0; i < t.Rows; i++ {
		for j := 0; j < t.Cols; j++ {
			if t.gameboard[i][j] >= 0 {
				fmt.Printf("%d", t.gameboard[i][j])
			} else {
				fmt.Printf("-")
			}
		}
		fmt.Printf("\n")
	}
}

func CreatePond(name string, rows int, cols int, rules func(*Pond, OrganismReference), init func(*Pond)) *Pond {
	p := new(Pond)
	p.Name = name
	p.Rows = rows
	p.Cols = cols
	p.NumLiving = 0
	p.ruleset = rules
	p.initializer = init
	p.Status = Active

	p.initializer(p)

	return p
}
