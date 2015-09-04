package cgol

import (
	"bytes"
	"errors"
	"strconv"
)

type LifeboardLocation struct {
	X int
	Y int
}

func (t *LifeboardLocation) Equals(rhs *LifeboardLocation) bool {
	if t.X != rhs.X {
		return false
	}
	if t.Y != rhs.Y {
		return false
	}
	return true
}

func (t *LifeboardLocation) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(strconv.Itoa(t.X))
	buf.WriteString(",")
	buf.WriteString(strconv.Itoa(t.Y))
	buf.WriteString("]")
	return buf.String()
}

type LifeboardDims struct {
	Height int
	Width  int
}

func (t *LifeboardDims) GetCapacity() int {
	return t.Height * t.Width
}

func (t *LifeboardDims) String() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(t.Height))
	buf.WriteString("x")
	buf.WriteString(strconv.Itoa(t.Width))
	return buf.String()
}

type lifeboardReadOp struct {
	loc  LifeboardLocation
	resp chan int
}

type lifeboardSnapshotOp struct {
	resp chan [][]int
}

type lifeboardWriteOp struct {
	loc  LifeboardLocation
	val  int
	resp chan bool
}

type lifeboard struct {
	Dims               LifeboardDims
	lifeboardReads     chan *lifeboardReadOp
	lifeboardWrites    chan *lifeboardWriteOp
	lifeboardSnapshots chan *lifeboardSnapshotOp
}

func (t *lifeboard) lifeboard() {
	// Initialize the lifeboard
	var lifeboard = make([][]int, t.Dims.Height)
	completion := make(chan bool, t.Dims.Height)
	for i := 0; i < t.Dims.Height; i++ {
		go func(row int, c chan bool) {
			lifeboard[row] = make([]int, t.Dims.Width)
			for j := 0; j < t.Dims.Width; j++ {
				lifeboard[row][j] = -1
			}
			c <- true
		}(i, completion)
	}
	completed := 0
	for c := range completion {
		if c {
			completed++
			if completed >= t.Dims.Height {
				close(completion)
			}
		}
	}

	// Listen for requests
	for {
		select {
		case read := <-t.lifeboardReads:
			// FIXME: what if there is no value?
			read.resp <- lifeboard[read.loc.Y][read.loc.X]
		case write := <-t.lifeboardWrites:
			lifeboard[write.loc.Y][write.loc.X] = write.val
			write.resp <- true
		case snapshot := <-t.lifeboardSnapshots:
			snapshot.resp <- lifeboard
		}
	}
}

func (t *lifeboard) GetValue(location LifeboardLocation) (int, error) {
	// Check that the given location is valid
	if location.X < 0 || location.X > t.Dims.Width {
		return -1, errors.New("Given location is out of bounds")
	}
	if location.Y < 0 || location.Y > t.Dims.Height {
		return -1, errors.New("Given location is out of bounds")
	}

	read := &lifeboardReadOp{loc: location, resp: make(chan int)}
	t.lifeboardReads <- read
	val := <-read.resp

	return val, nil
}

func (t *lifeboard) getSnapshot() [][]int {
	snapshot := &lifeboardSnapshotOp{resp: make(chan [][]int)}
	t.lifeboardSnapshots <- snapshot
	val := <-snapshot.resp
	return val
}

func (t *lifeboard) SetValue(location LifeboardLocation, val int) error {
	// Check that the given location is valid
	if location.X < 0 || location.X > t.Dims.Width {
		return errors.New("Given location is out of bounds")
	}
	if location.Y < 0 || location.Y > t.Dims.Height {
		return errors.New("Given location is out of bounds")
	}

	// Write the value to the lifeboard
	write := &lifeboardWriteOp{loc: location, val: val, resp: make(chan bool)}
	t.lifeboardWrites <- write
	<-write.resp

	return nil
}

func (t *lifeboard) GetOrthogonalNeighbors(location LifeboardLocation) []LifeboardLocation {
	neighbors := make([]LifeboardLocation, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		neighbors = append(neighbors, LifeboardLocation{X: location.X, Y: above})
	}

	if below < t.Dims.Height {
		neighbors = append(neighbors, LifeboardLocation{X: location.X, Y: below})
	}

	if left >= 0 {
		neighbors = append(neighbors, LifeboardLocation{X: left, Y: location.Y})
	}

	if right < t.Dims.Width {
		neighbors = append(neighbors, LifeboardLocation{X: right, Y: location.Y})
	}

	// fmt.Printf("GetOrthogonalNeighbors(%s): %v\n", location.String(), neighbors)
	return neighbors
}

func (t *lifeboard) GetObliqueNeighbors(location LifeboardLocation) []LifeboardLocation {
	neighbors := make([]LifeboardLocation, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, LifeboardLocation{X: left, Y: above})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, LifeboardLocation{X: right, Y: above})
		}
	}

	if below < t.Dims.Height {
		if left >= 0 {
			neighbors = append(neighbors, LifeboardLocation{X: left, Y: below})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, LifeboardLocation{X: right, Y: below})
		}
	}

	return neighbors
}

func (t *lifeboard) GetAllNeighbors(location LifeboardLocation) []LifeboardLocation {
	neighbors := append(t.GetOrthogonalNeighbors(location), t.GetObliqueNeighbors(location)...)

	return neighbors
}

func (t *lifeboard) Equals(rhs *lifeboard) bool {
	rhsSnapshot := rhs.getSnapshot()
	thisSnapshot := t.getSnapshot()
	for row := t.Dims.Height - 1; row >= 0; row-- {
		for col := t.Dims.Width - 1; col >= 0; col-- {
			if thisSnapshot[row][col] != rhsSnapshot[row][col] {
				return false
			}
		}
	}
	return true
}

func (t *lifeboard) String() string {
	var buf bytes.Buffer

	buf.WriteString("lifeboard size: ")
	buf.WriteString(strconv.Itoa(t.Dims.Height))
	buf.WriteString("x")
	buf.WriteString(strconv.Itoa(t.Dims.Width))
	buf.WriteString("\n")

	// Draw out the matrix
	snapshot := t.getSnapshot()
	for i := 0; i < t.Dims.Height; i++ {
		for j := 0; j < t.Dims.Width; j++ {
			val := snapshot[i][j]
			if val >= 0 {
				buf.WriteString(strconv.Itoa(val))
			} else {
				buf.WriteString("-")
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func newLifeboard(dims LifeboardDims) (*lifeboard, error) {
	if dims.Height <= 0 || dims.Width <= 0 {
		return nil, errors.New("Dimensions must be greater than 0")
	}

	g := new(lifeboard)
	g.Dims = dims

	// Initialize the lifeboard and its channels
	g.lifeboardReads = make(chan *lifeboardReadOp)
	g.lifeboardWrites = make(chan *lifeboardWriteOp)
	g.lifeboardSnapshots = make(chan *lifeboardSnapshotOp)
	go g.lifeboard()

	return g, nil
}
