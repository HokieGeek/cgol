package cgol

import (
	"bytes"
	"fmt"
	"strconv"
)

type gameboardReadOp struct {
	key  OrganismReference
	resp chan int
}
type gameboardWriteOp struct {
	key  OrganismReference
	val  int
	resp chan bool
}

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

func (t *OrganismReference) String() string {
	/*
		var buf bytes.Buffer
		buf.WriteString("[")
		buf.WriteString(t.X)
		buf.WriteString(",")
		buf.WriteString(t.Y)
		buf.WriteString("]")
		return buf.String()
	*/
	return fmt.Sprintf("[%d,%d]", t.X, t.Y)
}

type Pond struct {
	Rows int
	Cols int
	// gameboard [][]int
	gameboardReads  chan *gameboardReadOp
	gameboardWrites chan *gameboardWriteOp

	NumLiving         int
	Status            GameStatus
	neighborsSelector NeighborsSelector
	initialOrganisms  []OrganismReference
}

func (t *Pond) gameboard() {
	// Initialize the gameboard
	var gameboard = make([][]int, t.Rows)
	// completion := make(chan int, t.Rows)
	for i := 0; i < t.Rows; i++ {
		/*
			go func(row int, c chan int) {
				// fmt.Printf("Doing: %d\n", row)
				gameboard[row] = make([]int, t.Cols)
				for j := 0; j < t.Cols; j++ {
					gameboard[row][j] = -1
				}
				c <- row
			}(i, completion)
		*/
		gameboard[i] = make([]int, t.Cols)
		for j := 0; j < t.Cols; j++ {
			gameboard[i][j] = -1
		}
	}
	// for c := range completion {
	// fmt.Printf("%d is done\n", c)
	// }

	// Listen for requests
	for {
		select {
		case read := <-t.gameboardReads:
			// fmt.Printf("gb read: %s\n", read.key.String())
			// fmt.Printf("gb size: %dx%d\n", t.Rows, t.Cols)
			// fmt.Printf("gb len1: %d\n", len(gameboard))
			// fmt.Printf("gb len2(%d): %d\n", read.key.X, len(gameboard[read.key.X]))
			read.resp <- gameboard[read.key.X][read.key.Y]
		case write := <-t.gameboardWrites:
			gameboard[write.key.X][write.key.Y] = write.val
			write.resp <- true
		}
	}
}

func (t *Pond) getOrthogonalNeighbors(organism OrganismReference) []OrganismReference {
	neighbors := make([]OrganismReference, 0)

	// Determine the offsets
	// ROWS = X, COLS = Y
	left := organism.Y - 1
	right := organism.Y + 1
	above := organism.X - 1
	below := organism.X + 1

	if above >= 0 {
		neighbors = append(neighbors, OrganismReference{X: above, Y: organism.Y})
	}

	if below < t.Rows {
		neighbors = append(neighbors, OrganismReference{X: below, Y: organism.Y})
	}

	if left >= 0 {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: left})
	}

	if right < t.Cols {
		neighbors = append(neighbors, OrganismReference{X: organism.X, Y: right})
	}

	// fmt.Printf("getOrthogonalNeighbors(%s): %v\n", organism.String(), neighbors)
	return neighbors
}

func (t *Pond) getObliqueNeighbors(organism OrganismReference) []OrganismReference {
	neighbors := make([]OrganismReference, 0)

	// Determine the offsets
	above := organism.Y - 1
	below := organism.Y + 1
	left := organism.X - 1
	right := organism.X + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, OrganismReference{X: left, Y: above})
		}
		if right < t.Cols {
			neighbors = append(neighbors, OrganismReference{X: right, Y: above})
		}
	}

	if below < t.Rows {
		if left >= 0 {
			neighbors = append(neighbors, OrganismReference{X: left, Y: below})
		}
		if right < t.Cols {
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
	// fmt.Printf("\tgetNeighborCount(%s)\n", organism.String())
	read := &gameboardReadOp{key: organism, resp: make(chan int)}
	t.gameboardReads <- read
	val := <-read.resp
	return val
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
	fmt.Printf("\tsetNeighborCount(%s, %d)\n", organism.String(), numNeighbors)
	originalNumNeighbors := t.getNeighborCount(organism)

	// Write the value to the gameboard
	write := &gameboardWriteOp{key: organism, val: numNeighbors, resp: make(chan bool)}
	t.gameboardWrites <- write
	<-write.resp

	// Update the living count if organism changed living state
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
	// Initialize the gameboard and its channels
	t.gameboardReads = make(chan *gameboardReadOp)
	t.gameboardWrites = make(chan *gameboardWriteOp)
	go t.gameboard()

	// Initialize the first organisms and set their neighbor counts
	t.initialOrganisms = append(t.initialOrganisms, initialLiving...)
	for _, initialOrganism := range initialLiving {
		t.setNeighborCount(initialOrganism, 0)
	}
	for _, initialOrganism := range initialLiving {
		livingNeighborsCount := 0
		// fmt.Printf("Looking for living neighbors of: %s\n", initialOrganism.String())
		// fmt.Print(t.String())
		for _, neighbor := range t.GetNeighbors(initialOrganism) {
			// fmt.Printf("Checking neighbor: %s\n", neighbor.String())
			if t.isOrganismAlive(neighbor) {
				// fmt.Printf("Live neighbor: %s\n", neighbor.String())
				livingNeighborsCount++
			}
		}
		t.setNeighborCount(initialOrganism, livingNeighborsCount)
	}
}

func (t *Pond) String() string {

	s := fmt.Sprintf("Size: %dx%d, Neighbor selection: %s", t.Rows, t.Cols, t.neighborsSelector)
	s = fmt.Sprintf("%s\nLiving organisms: %d\tStatus: %s\n", s, t.NumLiving, t.Status)

	var matrix bytes.Buffer
	for i := 0; i < t.Rows; i++ {
		for j := 0; j < t.Cols; j++ {
			neighborCount := t.getNeighborCount(OrganismReference{X: i, Y: j})
			if neighborCount >= 0 {
				matrix.WriteString(strconv.Itoa(neighborCount))
				// fmt.Printf("%d", neighborCount)
			} else {
				matrix.WriteString("-")
				// fmt.Printf("-")
			}
		}
		matrix.WriteString("\n")
		//fmt.Printf("\n")
	}

	return s + matrix.String()
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
