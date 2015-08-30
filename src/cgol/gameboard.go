package cgol

import (
	"bytes"
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

type GameboardDims struct {
	Height int
	Width  int
}

func (t *GameboardDims) GetCapacity() int {
	return t.Height * t.Width
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
	Dims               GameboardDims
	gameboardReads     chan *gameboardReadOp
	gameboardWrites    chan *gameboardWriteOp
	gameboardSnapshots chan *gameboardSnapshotOp
}

func (t *Gameboard) gameboard() {
	// Initialize the gameboard
	var gameboard = make([][]int, t.Dims.Height)
	completion := make(chan bool, t.Dims.Height)
	for i := 0; i < t.Dims.Height; i++ {
		go func(row int, c chan bool) {
			gameboard[row] = make([]int, t.Dims.Width)
			for j := 0; j < t.Dims.Width; j++ {
				gameboard[row][j] = -1
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
		case read := <-t.gameboardReads:
			read.resp <- gameboard[read.loc.Y][read.loc.X]
		case write := <-t.gameboardWrites:
			gameboard[write.loc.Y][write.loc.X] = write.val
			write.resp <- true
		case snapshot := <-t.gameboardSnapshots:
			snapshot.resp <- gameboard
		}
	}
}

func (t *Gameboard) GetValue(location GameboardLocation) int {
	read := &gameboardReadOp{loc: location, resp: make(chan int)}
	t.gameboardReads <- read
	val := <-read.resp
	return val
}

func (t *Gameboard) getSnapshot() [][]int {
	snapshot := &gameboardSnapshotOp{resp: make(chan [][]int)}
	t.gameboardSnapshots <- snapshot
	val := <-snapshot.resp
	return val
}

func (t *Gameboard) SetValue(location GameboardLocation, val int) {
	// Write the value to the gameboard
	write := &gameboardWriteOp{loc: location, val: val, resp: make(chan bool)}
	t.gameboardWrites <- write
	<-write.resp
}

func (t *Gameboard) GetOrthogonalNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := make([]GameboardLocation, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		neighbors = append(neighbors, GameboardLocation{X: location.X, Y: above})
	}

	if below < t.Dims.Height {
		neighbors = append(neighbors, GameboardLocation{X: location.X, Y: below})
	}

	if left >= 0 {
		neighbors = append(neighbors, GameboardLocation{X: left, Y: location.Y})
	}

	if right < t.Dims.Width {
		neighbors = append(neighbors, GameboardLocation{X: right, Y: location.Y})
	}

	// fmt.Printf("GetOrthogonalNeighbors(%s): %v\n", location.String(), neighbors)
	return neighbors
}

func (t *Gameboard) GetObliqueNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := make([]GameboardLocation, 0)

	// Determine the offsets
	left := location.X - 1
	right := location.X + 1
	above := location.Y - 1
	below := location.Y + 1

	if above >= 0 {
		if left >= 0 {
			neighbors = append(neighbors, GameboardLocation{X: left, Y: above})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, GameboardLocation{X: right, Y: above})
		}
	}

	if below < t.Dims.Height {
		if left >= 0 {
			neighbors = append(neighbors, GameboardLocation{X: left, Y: below})
		}
		if right < t.Dims.Width {
			neighbors = append(neighbors, GameboardLocation{X: right, Y: below})
		}
	}

	return neighbors
}

func (t *Gameboard) GetAllNeighbors(location GameboardLocation) []GameboardLocation {
	neighbors := append(t.GetOrthogonalNeighbors(location), t.GetObliqueNeighbors(location)...)

	return neighbors
}

func (t *Gameboard) Equals(rhs *Gameboard) bool {
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

func (t *Gameboard) String() string {
	var buf bytes.Buffer

	buf.WriteString("Gameboard size: ")
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

func NewGameboard(dims GameboardDims) *Gameboard {
	g := new(Gameboard)
	g.Dims = dims

	// Initialize the gameboard and its channels
	g.gameboardReads = make(chan *gameboardReadOp)
	g.gameboardWrites = make(chan *gameboardWriteOp)
	g.gameboardSnapshots = make(chan *gameboardSnapshotOp)
	go g.gameboard()

	return g
}
