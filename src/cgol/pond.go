package cgol

import "fmt"

/*
type gameboardReadOp struct {
    key  OrganismReference
    resp chan int
}
type gameboardWriteOp struct {
    key  OrganismReference
    val  int
    resp chan bool
}

type Gameboard struct {
    Rows int
    Cols int
    reads chan *gameboardReadOp
    writes chan *gameboardWriteOp
}

func (t *Gameboard) board() {
     var gameboard = make([][]int, t.Rows)
        for i := 0; i < t.Rows; i++ {
		gameboard[0] = make([]int, t.Cols)
	 }

        for {
            select {
            case read := <-t.reads:
                read.resp <- gameboard[read.key.X][read.key.Y]
            case write := <-t.writes:
                gameboard[write.key.X][write.key.Y] = write.val
                write.resp <- true
            }
        }
}

func (t *Gameboard) SetOrganism(ref OrganismReference, v int) {
    write := &gameboardWriteOp {key: ref, val: v, resp: make(chan bool)}
    t.writes <- write
    <-write.resp
}

func (t *Gameboard) GetOrganism(ref OrganismReference) int {
    read := &gameboardReadOp{key: ref, resp: make(chan int)}
    t.reads <- read
    val := <-read.resp
    return val
}

func (t *Gameboard) init() {
    // Setup the gameboard
    t.Rows = 2
    t.Cols = 5
    t.reads = make(chan *gameboardReadOp)
    t.writes = make(chan *gameboardWriteOp)
    go t.board()
}

func CreateGameboard(rows int, cols int) {
}

func testGameboard() {
	testing := new(Gameboard)
	testing.init()

	ref := OrganismReference {X: 0, Y: 0}
    testing.SetOrganism(ref, 738)
    fmt.Printf("Testing: %d\n", testing.GetOrganism(ref))
}
*/

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

type NeighborsSelector int

const (
	NEIGHBORS_ORTHOGONAL NeighborsSelector = 1
	NEIGHBORS_OBLIQUE    NeighborsSelector = 2
	NEIGHBORS_ALL        NeighborsSelector = 3
)

func (t NeighborsSelector) String() string {
	s := ""

	if t&NEIGHBORS_ORTHOGONAL == NEIGHBORS_ORTHOGONAL {
		s += "NEIGHBORS_ORTHOGONAL"
	} else if t&NEIGHBORS_OBLIQUE == NEIGHBORS_OBLIQUE {
		s += "NEIGHBORS_OBLIQUE"
	} else if t&NEIGHBORS_ALL == NEIGHBORS_ALL {
		s += "NEIGHBORS_ALL"
	}

	return s
}

type OrganismReference struct {
	X int
	Y int
}

type Pond struct {
	Rows      int
	Cols      int
	gameboard [][]int

	NumLiving         int
	Status            GameStatus
	neighborsSelector NeighborsSelector
	initialOrganisms  []OrganismReference
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

	return neighbors
}

func (t *Pond) GetNeighbors(organism OrganismReference) []OrganismReference {
	switch {
	case t.neighborsSelector == NEIGHBORS_ORTHOGONAL:
		return t.getOrthogonalNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_OBLIQUE:
		return t.getObliqueNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_ALL:
		return t.getAllNeighbors(organism)
	}

	return make([]OrganismReference, 0)
}

func (t *Pond) isOrganismAlive(organism OrganismReference) bool {
	return (t.getNeighborCount(organism) >= 0)
}

func (t *Pond) getNeighborCount(organism OrganismReference) int {
	return t.gameboard[organism.X][organism.Y]
}

func (t *Pond) calculateNeighborCount(organism OrganismReference) int {
	numLivingNeighbors := 0
	for _, neighbor := range t.GetNeighbors(organism) {
		if t.isOrganismAlive(neighbor) {
			numLivingNeighbors++
		}
	}
	return numLivingNeighbors
}

func (t *Pond) setNeighborCount(organism OrganismReference, numNeighbors int) {
	// TODO: Mutex protection
	originalNumNeighbors := t.getNeighborCount(organism)

	t.gameboard[organism.X][organism.Y] = numNeighbors

	// Update the living count if organism changed living state
	// fmt.Printf("Original: %d vs New: %d\n", originalNumNeighbors, numNeighbors)
	if originalNumNeighbors < 0 && numNeighbors >= 0 {
		t.NumLiving++
	} else if originalNumNeighbors >= 0 && numNeighbors < 0 {
		t.NumLiving--
	}
}

func (t *Pond) incrementNeighborCount(organism OrganismReference) {
	t.setNeighborCount(organism, t.getNeighborCount(organism)+1)
}

func (t *Pond) decrementNeighborCount(organism OrganismReference) {
	t.setNeighborCount(organism, t.getNeighborCount(organism)-1)
}

func (t *Pond) init(initialLiving []OrganismReference) {
	t.initialOrganisms = append(t.initialOrganisms, initialLiving...)
	t.gameboard = make([][]int, t.Rows)

	// completion := make(chan int, pond.Rows)
	for i := 0; i < t.Rows; i++ {
		// go func(row int, c chan int) {
		// 	fmt.Printf("Doing: %d\n", row)
		// 	c <- row
		// }(i, completion)
		t.gameboard[i] = make([]int, t.Cols)
		for j := 0; j < t.Cols; j++ {
			t.gameboard[i][j] = -1
			// t.setNeighborCount(OrganismReference{X: i, Y: j}, -1)
		}
	}

	for _, initialOrganism := range initialLiving {
		t.setNeighborCount(initialOrganism, 0) // TODO: hmm...
	}
	// for c := range completion {
	// 	fmt.Printf("%d is done\n", c)
	// }
}

func (t *Pond) Display() {
	fmt.Printf("Size: %dx%d, Neighbor selection: %s\n", t.Rows, t.Cols, t.neighborsSelector)
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

func CreatePond(rows int, cols int, neighbors NeighborsSelector) *Pond {
	p := new(Pond)

	// Default values
	p.NumLiving = 0
	p.Status = Active

	// Add the given values
	p.Rows = rows
	p.Cols = cols
	p.neighborsSelector = neighbors

	return p
}
