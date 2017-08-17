package life

import (
	"bytes"
	"errors"
	"fmt"
	// "io/ioutil"
	// "log"
	"strconv"
	// "os"
)

type neighborsSelector int

const (
	NEIGHBORS_ALL neighborsSelector = iota
	NEIGHBORS_ORTHOGONAL
	NEIGHBORS_OBLIQUE
)

func (t neighborsSelector) String() string {
	switch t {
	case NEIGHBORS_ALL:
		return "All"
	case NEIGHBORS_ORTHOGONAL:
		return "Orthogonal"
	case NEIGHBORS_OBLIQUE:
		return "Oblique"
	}
	return "Unknown"
}

type pond struct {
	board             *board
	neighborsSelector neighborsSelector
	living            *tracker
}

func (t *pond) GetNeighbors(organism Location) ([]Location, error) {
	if !t.board.isValidLocation(organism) {
		return nil, errors.New("Location is out of bounds")
	}

	switch {
	case t.neighborsSelector == NEIGHBORS_ORTHOGONAL:
		return t.board.GetOrthogonalNeighbors(organism), nil
	case t.neighborsSelector == NEIGHBORS_OBLIQUE:
		return t.board.GetObliqueNeighbors(organism), nil
	case t.neighborsSelector == NEIGHBORS_ALL:
		return t.board.GetAllNeighbors(organism), nil
	}

	return nil, errors.New("Did not recognize neighbor selector")
}

func (t *pond) isOrganismAlive(organism Location) bool {
	return (t.GetOrganismValue(organism) >= 0)
}

func (t *pond) GetNumLiving() int {
	return t.living.GetCount()
}

func (t *pond) GetOrganismValue(organism Location) int {
	// fmt.Printf("\tgetNeighborCount(%s)\n", organism.String())
	val, err := t.board.GetValue(organism)

	if err != nil {
		// TODO: print the error
		return -1
	}

	return val
}

func (t *pond) setOrganismValue(organism Location, num int) {
	// fmt.Printf("\tsetNeighborCount(%s, %d)\n", organism.String(), num)
	originalNum := t.GetOrganismValue(organism)

	// Write the value to the board
	err := t.board.SetValue(organism, num)
	if err != nil {
		fmt.Println("ERROR!")
		panic(err)
	}

	// Update the living count if organism changed living state
	if originalNum < 0 && num >= 0 {
		t.living.Set(organism)
	} else if originalNum >= 0 && num < 0 {
		t.living.Remove(organism)
	} else {
		fmt.Printf(">>>>> what am i doing in here?: %d\n", num)
	}
	// fmt.Printf("Living count is: %d\n", t.living.GetCount())
}

func (t *pond) calculateNeighborCount(organism Location) (int, []Location) {
	numNeighbors := 0
	neighbors, err := t.GetNeighbors(organism)
	if err != nil {
		// FIXME
	}
	for _, neighbor := range neighbors {
		if t.isOrganismAlive(neighbor) {
			numNeighbors++
		}
	}
	return numNeighbors, neighbors
}

func (t *pond) SetOrganisms(organisms []Location) {
	// Initialize the first organisms and set their neighbor counts
	for _, organism := range organisms {
		t.setOrganismValue(organism, 0)
	}
}

func (t *pond) Clone() (*pond, error) {
	shadowBoard, err := t.board.Clone()
	if err != nil {
		return nil, err
	}

	shadowpond, err := newPond(shadowBoard, t.living.Clone(), t.neighborsSelector) // FIXME: umm. board and tracker need to be cloned as well
	if err != nil {
		return nil, err
	}

	shadowpond.neighborsSelector = t.neighborsSelector

	shadowpond.SetOrganisms(t.living.GetAll())

	return shadowpond, nil
}

func (t *pond) Equals(rhs *pond) bool {
	if !t.board.Equals(rhs.board) {
		return false
	}
	if t.neighborsSelector != rhs.neighborsSelector {
		return false
	}
	return true
}

func (t *pond) String() string {
	var buf bytes.Buffer
	buf.WriteString("Neighbors: ")
	buf.WriteString(t.neighborsSelector.String())
	buf.WriteString("\tLiving cells: ")
	buf.WriteString(strconv.Itoa(t.living.GetCount()))
	buf.WriteString("\n")
	buf.WriteString(t.board.String())

	return buf.String()
}

func newPond(board *board, tracker *tracker, neighbors neighborsSelector) (*pond, error) {
	p := new(pond)

	if board == nil || tracker == nil {
		return nil, errors.New("Neither board nor tracker can be nil")
	}

	p.living = tracker
	p.board = board
	p.neighborsSelector = neighbors

	return p, nil
}

// vim: set foldmethod=marker:
