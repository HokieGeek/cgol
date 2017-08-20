package life

import (
	"bytes"
	"errors"
	"strconv"
)

// Location is a simple coordinate structure
type Location struct {
	X int
	Y int
}

// Equals is a basic equality test for the given object
func (t *Location) Equals(rhs *Location) bool {
	if t.X != rhs.X {
		return false
	}
	if t.Y != rhs.Y {
		return false
	}
	return true
}

func (t *Location) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(strconv.Itoa(t.X))
	buf.WriteString(",")
	buf.WriteString(strconv.Itoa(t.Y))
	buf.WriteString("]")
	return buf.String()
}

// Dimensions is a simple structure capturing the given dimensions
type Dimensions struct {
	Width  int
	Height int
}

// Capacity returns the number of cells possible with the given dimensions
func (t *Dimensions) Capacity() int {
	return t.Width * t.Height
}

// Equals is a basic equality test for the given object
func (t *Dimensions) Equals(rhs *Dimensions) bool {
	if t.Width != rhs.Width {
		return false
	}

	if t.Height != rhs.Height {
		return false
	}

	return true
}

func (t *Dimensions) String() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(t.Width))
	buf.WriteString("x")
	buf.WriteString(strconv.Itoa(t.Height))
	return buf.String()
}

type neighborsSelector int

// Enumeration of the method to select neighbors
const (
	NeighborsAll neighborsSelector = iota
	NeighborsOrthogonal
	NeighborsOblique
)

func (t neighborsSelector) String() string {
	switch t {
	case NeighborsAll:
		return "All"
	case NeighborsOrthogonal:
		return "Orthogonal"
	case NeighborsOblique:
		return "Oblique"
	}
	return "Unknown"
}

type pond struct {
	Dims              Dimensions
	neighborsSelector neighborsSelector
	living            *tracker
}

func (t *pond) getOrthogonalNeighbors(location Location) []Location {
	neighbors := make([]Location, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		neighbors = append(neighbors, Location{X: location.X, Y: above})
	}

	if below < t.Dims.Height {
		neighbors = append(neighbors, Location{X: location.X, Y: below})
	}

	if left >= 0 {
		neighbors = append(neighbors, Location{X: left, Y: location.Y})
	}

	if right < t.Dims.Width {
		neighbors = append(neighbors, Location{X: right, Y: location.Y})
	}

	// fmt.Printf("GetOrthogonalNeighbors(%s): %v\n", location.String(), neighbors)
	return neighbors
}

func (t *pond) getObliqueNeighbors(location Location) []Location {
	neighbors := make([]Location, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, Location{X: left, Y: above})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, Location{X: right, Y: above})
		}
	}

	if below < t.Dims.Height {
		if left >= 0 {
			neighbors = append(neighbors, Location{X: left, Y: below})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, Location{X: right, Y: below})
		}
	}

	return neighbors
}

func (t *pond) getAllNeighbors(location Location) []Location {
	neighbors := append(t.getOrthogonalNeighbors(location), t.getObliqueNeighbors(location)...)

	return neighbors
}

func (t *pond) GetNeighbors(organism Location) ([]Location, error) {
	if !t.isValidLocation(organism) {
		return nil, errors.New("Location is out of bounds")
	}

	switch {
	case t.neighborsSelector == NeighborsOrthogonal:
		return t.getOrthogonalNeighbors(organism), nil
	case t.neighborsSelector == NeighborsOblique:
		return t.getObliqueNeighbors(organism), nil
	case t.neighborsSelector == NeighborsAll:
		return t.getAllNeighbors(organism), nil
	}

	return nil, errors.New("Did not recognize neighbor selector")
}

func (t *pond) isValidLocation(location Location) bool {
	if location.X < 0 || location.X > t.Dims.Width {
		return false
	}
	if location.Y < 0 || location.Y > t.Dims.Height {
		return false
	}
	return true
}

func (t *pond) isOrganismAlive(organism Location) bool {
	// return (t.GetOrganismValue(organism) >= 0)
	return t.living.Test(organism)
}

func (t *pond) setOrganismState(organism Location, alive bool) {
	// fmt.Printf("\tsetNeighborCount(%s, %d)\n", organism.String(), num)
	originalState := t.isOrganismAlive(organism)

	// Only do the deed if something has changed TODO: is this a stupid optimization?
	if originalState != alive {
		if alive {
			t.living.Set(organism)
		} else {
			t.living.Remove(organism)
		}
		// fmt.Printf("Living count is: %d\n", t.living.Count())
	}
}

func (t *pond) SetOrganisms(organisms []Location) {
	// Initialize the first organisms and set their neighbor counts
	for _, organism := range organisms {
		t.setOrganismState(organism, true)
	}
}

func (t *pond) Clone() (*pond, error) {
	shadowpond, err := newPond(t.Dims, t.living.Clone(), t.neighborsSelector)
	if err != nil {
		return nil, err
	}

	shadowpond.neighborsSelector = t.neighborsSelector

	shadowpond.SetOrganisms(t.living.GetAll())

	return shadowpond, nil
}

func (t *pond) Equals(rhs *pond) bool {
	if !t.living.Equals(rhs.living) {
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
	buf.WriteString(strconv.Itoa(t.living.Count()))
	buf.WriteString("\n")

	//// DRAW THE BOARD ////
	buf.WriteString("Size: ")
	buf.WriteString(t.Dims.String())
	buf.WriteString("\n")

	// Draw the top border
	buf.WriteString("┌")
	for j := t.Dims.Width; j > 0; j-- {
		buf.WriteString("─")
	}
	buf.WriteString("┐\n")

	// Draw out the matrix
	for y := 0; y < t.Dims.Height; y++ {
		buf.WriteString("│") // Left border
		for x := 0; x < t.Dims.Width; x++ {
			if t.isOrganismAlive(Location{X: x, Y: y}) {
				buf.WriteString("0")
			} else {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("│\n") // Right border
	}

	// Draw the bottom border
	buf.WriteString("└")
	for j := t.Dims.Width; j > 0; j-- {
		buf.WriteString("─")
	}
	buf.WriteString("┘\n")

	return buf.String()
}

func newPond(dims Dimensions, tracker *tracker, neighbors neighborsSelector) (*pond, error) {
	if dims.Capacity() == 0 {
		return nil, errors.New("Cannot create pond of zero capacity")
	}

	p := new(pond)

	if tracker == nil {
		return nil, errors.New("tracker cannot be nil")
	}

	p.living = tracker
	p.neighborsSelector = neighbors

	p.Dims = dims

	return p, nil
}

// vim: set foldmethod=marker:
