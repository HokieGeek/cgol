package life

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type Location struct {
	X int
	Y int
}

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

type Dimensions struct {
	Width  int
	Height int
}

func (t *Dimensions) Capacity() int {
	return t.Width * t.Height
}

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

type boardReadOp struct {
	loc  Location
	resp chan int
}

type boardSnapshotOp struct {
	resp chan [][]int
}

type boardWriteOp struct {
	loc  Location
	val  int
	resp chan bool
}

type board struct {
	Dims           Dimensions
	boardReads     chan *boardReadOp
	boardWrites    chan *boardWriteOp
	boardSnapshots chan *boardSnapshotOp
}

func (t *board) board() {
	// Initialize the board
	var board = make([][]int, t.Dims.Height)
	completion := make(chan bool, t.Dims.Height)
	for i := 0; i < t.Dims.Height; i++ {
		go func(row int, c chan bool) {
			board[row] = make([]int, t.Dims.Width)
			for j := 0; j < t.Dims.Width; j++ {
				board[row][j] = -1
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
		case read := <-t.boardReads:
			// FIXME: what if there is no value?
			read.resp <- board[read.loc.Y][read.loc.X]
		case write := <-t.boardWrites:
			board[write.loc.Y][write.loc.X] = write.val
			write.resp <- true
		case snapshot := <-t.boardSnapshots:
			snapshot.resp <- board
		}
	}
}

func (t *board) isValidLocation(location Location) bool {
	if location.X < 0 || location.X > t.Dims.Width {
		return false
	}
	if location.Y < 0 || location.Y > t.Dims.Height {
		return false
	}
	return true
}

func (t *board) GetValue(location Location) (int, error) {
	// Check that the given location is valid
	if !t.isValidLocation(location) {
		return -1, errors.New(fmt.Sprintf("Given location is out of bounds: %s", location.String()))
	}

	read := &boardReadOp{loc: location, resp: make(chan int)}
	t.boardReads <- read
	val := <-read.resp

	return val, nil
}

func (t *board) getSnapshot() [][]int {
	snapshot := &boardSnapshotOp{resp: make(chan [][]int)}
	t.boardSnapshots <- snapshot
	val := <-snapshot.resp
	return val
}

func (t *board) SetValue(location Location, val int) error {
	// Check that the given location is valid
	if !t.isValidLocation(location) {
		return errors.New(fmt.Sprintf("Given location is out of bounds: %s", location.String()))
	}

	// Write the value to the board
	write := &boardWriteOp{loc: location, val: val, resp: make(chan bool)}
	t.boardWrites <- write
	<-write.resp

	return nil
}

func (t *board) GetOrthogonalNeighbors(location Location) []Location {
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

func (t *board) GetObliqueNeighbors(location Location) []Location {
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

func (t *board) GetAllNeighbors(location Location) []Location {
	neighbors := append(t.GetOrthogonalNeighbors(location), t.GetObliqueNeighbors(location)...)

	return neighbors
}

func (t *board) Clone() (*board, error) {
	shadow, err := newBoard(t.Dims)
	if err != nil {
		return nil, err
	}

	// FIXME: copy the board using t.getSnapshot()

	return shadow, nil
}

func (t *board) Equals(rhs *board) bool {
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

func (t *board) String() string {
	var buf bytes.Buffer

	buf.WriteString("Board size: ")
	buf.WriteString(t.Dims.String())
	buf.WriteString("\n")

	// Draw the top border
	buf.WriteString("┌")
	for j := t.Dims.Width; j > 0; j-- {
		buf.WriteString("─")
	}
	buf.WriteString("┐\n")

	// Draw out the matrix
	snapshot := t.getSnapshot()
	for i := 0; i < t.Dims.Height; i++ {
		buf.WriteString("│") // Left border
		for j := 0; j < t.Dims.Width; j++ {
			if val := snapshot[i][j]; val >= 0 {
				buf.WriteString(strconv.Itoa(val))
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

func newBoard(dims Dimensions) (*board, error) {
	if dims.Height <= 0 || dims.Width <= 0 {
		return nil, errors.New("Dimensions must be greater than 0")
	}

	g := new(board)
	g.Dims = dims

	// Initialize the board and its channels
	g.boardReads = make(chan *boardReadOp)
	g.boardWrites = make(chan *boardWriteOp)
	g.boardSnapshots = make(chan *boardSnapshotOp)
	go g.board()

	return g, nil
}

// vim: set foldmethod=marker:
