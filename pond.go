package life

import (
	"bytes"
	// "fmt"
	// "io/ioutil"
	// "log"
	"strconv"
	// "os"
)

type neighborsSelector int

const (
	NEIGHBORS_ORTHOGONAL neighborsSelector = 1
	NEIGHBORS_OBLIQUE    neighborsSelector = 2
	NEIGHBORS_ALL        neighborsSelector = 3
)

func (t neighborsSelector) String() string {
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

type livingTrackerAddOp struct {
	loc  Location
	resp chan bool
}

type livingTrackerRemoveOp struct {
	loc  Location
	resp chan bool
}

type livingTrackerTestOp struct {
	loc  Location
	resp chan bool
}

type livingTrackerGetAllOp struct {
	resp chan []Location
}

type livingTrackerCountOp struct {
	resp chan int
}

type livingTracker struct {
	trackerAdd    chan *livingTrackerAddOp
	trackerRemove chan *livingTrackerRemoveOp
	trackerTest   chan *livingTrackerTestOp
	trackerGetAll chan *livingTrackerGetAllOp
	trackerCount  chan *livingTrackerCountOp
}

func (t *livingTracker) living() {
	var livingMap = make(map[int]map[int]Location)
	var count int
	// logger := log.New(os.Stderr, "livingTracker: ", log.Ltime)
	// logger := log.New(ioutil.Discard, "livingTracker: ", log.Ltime)

	for {
		select {
		case add := <-t.trackerAdd:
			added := true
			_, keyExists := livingMap[add.loc.Y]
			if !keyExists {
				livingMap[add.loc.Y] = make(map[int]Location)
			}
			_, keyExists = livingMap[add.loc.Y][add.loc.X]
			if !keyExists {
				livingMap[add.loc.Y][add.loc.X] = add.loc
				count++
			}
			add.resp <- added
		case remove := <-t.trackerRemove:
			removed := false
			_, keyExists := livingMap[remove.loc.Y]
			if keyExists {
				_, keyExists = livingMap[remove.loc.Y][remove.loc.X]
				if keyExists {
					delete(livingMap[remove.loc.Y], remove.loc.X)
					removed = true
					count--

					// TODO Delete the row if it has no children?
					// if len(livingMap[remove.loc.Y]) <= 0 {
					// 	delete(livingMap, remove.loc.Y)
					// }
				}
			}
			remove.resp <- removed
		case test := <-t.trackerTest:
			_, keyExists := livingMap[test.loc.Y]
			if keyExists {
				_, keyExists = livingMap[test.loc.Y][test.loc.X]
				if !keyExists {
					test.resp <- false
				} else {
					test.resp <- true
				}
			} else {
				test.resp <- false
			}
		case getall := <-t.trackerGetAll:
			all := make([]Location, 0)
			for rowNum := range livingMap {
				for _, col := range livingMap[rowNum] {
					all = append(all, col)
				}
			}
			getall.resp <- all
		case countOp := <-t.trackerCount:
			countOp.resp <- count
		}
	}
}

func (t *livingTracker) Set(location Location) {
	add := &livingTrackerAddOp{loc: location, resp: make(chan bool)}
	t.trackerAdd <- add
	<-add.resp
}

func (t *livingTracker) Remove(location Location) {
	remove := &livingTrackerRemoveOp{loc: location, resp: make(chan bool)}
	t.trackerRemove <- remove
	<-remove.resp
}

func (t *livingTracker) Test(location Location) bool {
	read := &livingTrackerTestOp{loc: location, resp: make(chan bool)}
	t.trackerTest <- read
	val := <-read.resp

	return val
}

func (t *livingTracker) GetAll() []Location {
	get := &livingTrackerGetAllOp{resp: make(chan []Location)}
	t.trackerGetAll <- get
	val := <-get.resp

	return val
}

func (t *livingTracker) GetCount() int {
	count := &livingTrackerCountOp{resp: make(chan int)}
	t.trackerCount <- count
	val := <-count.resp

	return val
}

func newLivingTracker() *livingTracker {
	t := new(livingTracker)

	t.trackerAdd = make(chan *livingTrackerAddOp)
	t.trackerRemove = make(chan *livingTrackerRemoveOp)
	t.trackerTest = make(chan *livingTrackerTestOp)
	t.trackerGetAll = make(chan *livingTrackerGetAllOp)
	t.trackerCount = make(chan *livingTrackerCountOp)

	go t.living()

	return t
}

type pond struct {
	board             *board
	neighborsSelector neighborsSelector
	living            *livingTracker
}

func (t *pond) GetNeighbors(organism Location) []Location {
	switch {
	case t.neighborsSelector == NEIGHBORS_ORTHOGONAL:
		return t.board.GetOrthogonalNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_OBLIQUE:
		return t.board.GetObliqueNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_ALL:
		return t.board.GetAllNeighbors(organism)
	}

	return nil
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
	t.board.SetValue(organism, num)

	// Update the living count if organism changed living state
	if originalNum < 0 && num >= 0 {
		t.living.Set(organism)
	} else if originalNum >= 0 && num < 0 {
		t.living.Remove(organism)
	}
}

func (t *pond) calculateNeighborCount(organism Location) (int, []Location) {
	numNeighbors := 0
	neighbors := t.GetNeighbors(organism)
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

func (t *pond) GetLifeboard() [][]int {
	return t.board.getSnapshot()
}

func (t *pond) Clone() (*pond, error) {
	shadowpond, err := newpond(t.board.Dims, t.neighborsSelector)
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
	buf.WriteString("Neighbor selection: ")
	buf.WriteString(t.neighborsSelector.String())
	buf.WriteString("\nLiving organisms: ")
	buf.WriteString(strconv.Itoa(t.living.GetCount()))
	buf.WriteString("\n")
	buf.WriteString(t.board.String())

	return buf.String()
}

func newpond(dims Dimensions, neighbors neighborsSelector) (*pond, error) {
	p := new(pond)

	// Create values
	p.living = newLivingTracker()

	// Add the given values
	var err error
	p.board, err = newLifeboard(dims)
	if err != nil {
		return nil, err
	}
	p.neighborsSelector = neighbors

	return p, nil
}
