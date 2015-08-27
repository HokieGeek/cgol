package cgol

import (
	"bytes"
	// "fmt"
	"strconv"
)

type GameboardLocation struct {
	X int
	Y int
}

func (t *GameboardLocation) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(strconv.Itoa(t.X))
	buf.WriteString(",")
	buf.WriteString(strconv.Itoa(t.Y))
	buf.WriteString("]")
	return buf.String()
}

// TODO: make use of this
type GameboardDims struct {
	Height int
	Width  int
}

func (t *GameboardDims) String() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(t.Height))
	buf.WriteString("x")
	buf.WriteString(strconv.Itoa(t.Width))
	return buf.String()
}

type gameboardReadOp struct {
	loc  GameboardLocation
	resp chan int
}

type gameboardSnapshotOp struct {
	resp chan [][]int
}

type gameboardWriteOp struct {
	loc  GameboardLocation
	val  int
	resp chan bool
}

type Gameboard struct {
	Rows               int
	Cols               int
	Dims               GameboardDims
	gameboardReads     chan *gameboardReadOp
	gameboardWrites    chan *gameboardWriteOp
	gameboardSnapshots chan *gameboardSnapshotOp
}

func (t *Gameboard) gameboard() {
	// Initialize the gameboard
	var gameboard = make([][]int, t.Rows)
	completion := make(chan bool, t.Rows)
	for i := 0; i < t.Rows; i++ {
		go func(row int, c chan bool) {
			gameboard[row] = make([]int, t.Cols)
			for j := 0; j < t.Cols; j++ {
				gameboard[row][j] = -1
			}
			c <- true
		}(i, completion)
	}
	completed := 0
	for c := range completion {
		if c {
			completed++
			if completed >= t.Rows {
				close(completion)
			}
		}
	}

	// Listen for requests
	for {
		select {
		case read := <-t.gameboardReads:
			read.resp <- gameboard[read.loc.X][read.loc.Y]
		case write := <-t.gameboardWrites:
			gameboard[write.loc.X][write.loc.Y] = write.val
			write.resp <- true
		case snapshot := <-t.gameboardSnapshots:
			snapshot.resp <- gameboard
		}
	}
}

func (t *Gameboard) GetGameboardValue(location GameboardLocation) int {
	// fmt.Printf("\tGetGameboardValue(%s)\n", location.String())
	read := &gameboardReadOp{loc: location, resp: make(chan int)}
	t.gameboardReads <- read
	val := <-read.resp
	return val
}

func (t *Gameboard) GetGameboardSnapshot() [][]int {
	snapshot := &gameboardSnapshotOp{resp: make(chan [][]int)}
	t.gameboardSnapshots <- snapshot
	val := <-snapshot.resp
	return val
}

func (t *Gameboard) SetGameboardValue(location GameboardLocation, val int) {
	// fmt.Printf("\tSetGameboardValue(%s, %d)\n", location.String(), val)
	// Write the value to the gameboard
	write := &gameboardWriteOp{loc: location, val: val, resp: make(chan bool)}
	t.gameboardWrites <- write
	<-write.resp
}

// FIXME: using X for row and Y for col is idiotic
func (t *Gameboard) GetOrthogonalNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := make([]GameboardLocation, 0)

	// Determine the offsets
	// ROWS = X, COLS = Y
	left := location.Y - 1
	right := location.Y + 1
	above := location.X - 1
	below := location.X + 1

	if above >= 0 {
		neighbors = append(neighbors, GameboardLocation{X: above, Y: location.Y})
	}

	if below < t.Rows {
		neighbors = append(neighbors, GameboardLocation{X: below, Y: location.Y})
	}

	if left >= 0 {
		neighbors = append(neighbors, GameboardLocation{X: location.X, Y: left})
	}

	if right < t.Cols {
		neighbors = append(neighbors, GameboardLocation{X: location.X, Y: right})
	}

	// fmt.Printf("GetOrthogonalNeighbors(%s): %v\n", location.String(), neighbors)
	return neighbors
}

func (t *Gameboard) GetObliqueNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := make([]GameboardLocation, 0)

	// Determine the offsets
	above := location.X - 1
	below := location.X + 1
	left := location.Y - 1
	right := location.Y + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, GameboardLocation{X: left, Y: above})
		}
		if right < t.Cols {
			neighbors = append(neighbors, GameboardLocation{X: right, Y: above})
		}
	}

	if below < t.Rows {
		if left >= 0 {
			neighbors = append(neighbors, GameboardLocation{X: left, Y: below})
		}
		if right < t.Cols {
			neighbors = append(neighbors, GameboardLocation{X: right, Y: below})
		}
	}

	return neighbors
}

func (t *Gameboard) GetAllNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := append(t.GetOrthogonalNeighbors(location), t.GetObliqueNeighbors(location)...)

	return neighbors
}

func (t *Gameboard) String() string {
	var buf bytes.Buffer

	buf.WriteString("Gameboard size: ")
	buf.WriteString(strconv.Itoa(t.Rows))
	buf.WriteString("x")
	buf.WriteString(strconv.Itoa(t.Cols))
	buf.WriteString("\n")

	// Draw out the matrix
	snapshot := t.GetGameboardSnapshot()
	for i := 0; i < t.Rows; i++ {
		for j := 0; j < t.Cols; j++ {
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

func NewGameboard(dims GameboardDims) *Gameboard {
	g := new(Gameboard)
	g.Dims = dims
	g.Rows = dims.Height
	g.Cols = dims.Width

	// Initialize the gameboard and its channels
	g.gameboardReads = make(chan *gameboardReadOp)
	g.gameboardWrites = make(chan *gameboardWriteOp)
	g.gameboardSnapshots = make(chan *gameboardSnapshotOp)
	go g.gameboard()

	return g
}
