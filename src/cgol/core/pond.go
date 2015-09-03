package cgol

import (
	"bytes"
	// "fmt"
	// "io/ioutil"
	// "log"
	"strconv"
	// "os"
)

type PondStatus int

const (
	Active PondStatus = 1
	Stable PondStatus = 2
	Dead   PondStatus = 3
)

func (t PondStatus) String() string {
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

type livingTrackerAddOp struct {
	loc  GameboardLocation
	resp chan bool
}

type livingTrackerRemoveOp struct {
	loc  GameboardLocation
	resp chan bool
}

type livingTrackerTestOp struct {
	loc  GameboardLocation
	resp chan bool
}

type livingTrackerGetAllOp struct {
	resp chan []GameboardLocation
}

type livingTrackerCountOp struct {
	resp chan int
}

type LivingTracker struct {
	trackerAdd    chan *livingTrackerAddOp
	trackerRemove chan *livingTrackerRemoveOp
	trackerTest   chan *livingTrackerTestOp
	trackerGetAll chan *livingTrackerGetAllOp
	trackerCount  chan *livingTrackerCountOp
}

func (t *LivingTracker) living() {
	var livingMap = make(map[int]map[int]GameboardLocation)
	var count int
	// logger := log.New(os.Stderr, "LivingTracker: ", log.Ltime)
	// logger := log.New(ioutil.Discard, "LivingTracker: ", log.Ltime)

	for {
		select {
		case add := <-t.trackerAdd:
			added := true
			_, keyExists := livingMap[add.loc.Y]
			if !keyExists {
				livingMap[add.loc.Y] = make(map[int]GameboardLocation)
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
			all := make([]GameboardLocation, 0)
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

func (t *LivingTracker) Set(location GameboardLocation) {
	add := &livingTrackerAddOp{loc: location, resp: make(chan bool)}
	t.trackerAdd <- add
	<-add.resp
}

func (t *LivingTracker) Remove(location GameboardLocation) {
	remove := &livingTrackerRemoveOp{loc: location, resp: make(chan bool)}
	t.trackerRemove <- remove
	<-remove.resp
}

func (t *LivingTracker) Test(location GameboardLocation) bool {
	read := &livingTrackerTestOp{loc: location, resp: make(chan bool)}
	t.trackerTest <- read
	val := <-read.resp

	return val
}

func (t *LivingTracker) GetAll() []GameboardLocation {
	get := &livingTrackerGetAllOp{resp: make(chan []GameboardLocation)}
	t.trackerGetAll <- get
	val := <-get.resp

	return val
}

func (t *LivingTracker) GetCount() int {
	count := &livingTrackerCountOp{resp: make(chan int)}
	t.trackerCount <- count
	val := <-count.resp

	return val
}

func NewLivingTracker() *LivingTracker {
	t := new(LivingTracker)

	t.trackerAdd = make(chan *livingTrackerAddOp)
	t.trackerRemove = make(chan *livingTrackerRemoveOp)
	t.trackerTest = make(chan *livingTrackerTestOp)
	t.trackerGetAll = make(chan *livingTrackerGetAllOp)
	t.trackerCount = make(chan *livingTrackerCountOp)

	go t.living()

	return t
}

type Pond struct {
	gameboard         *Gameboard
	Status            PondStatus
	neighborsSelector NeighborsSelector
	living            *LivingTracker
}

func (t *Pond) GetNeighbors(organism GameboardLocation) []GameboardLocation {
	switch {
	case t.neighborsSelector == NEIGHBORS_ORTHOGONAL:
		return t.gameboard.GetOrthogonalNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_OBLIQUE:
		return t.gameboard.GetObliqueNeighbors(organism)
	case t.neighborsSelector == NEIGHBORS_ALL:
		return t.gameboard.GetAllNeighbors(organism)
	}

	return nil
}

func (t *Pond) isOrganismAlive(organism GameboardLocation) bool {
	return (t.GetOrganismValue(organism) >= 0)
}

func (t *Pond) GetNumLiving() int {
	return t.living.GetCount()
}

func (t *Pond) GetOrganismValue(organism GameboardLocation) int {
	// fmt.Printf("\tgetNeighborCount(%s)\n", organism.String())
	val, err := t.gameboard.GetValue(organism)

	if err != nil {
		// TODO: print the error
		return -1
	}

	return val
}

func (t *Pond) setOrganismValue(organism GameboardLocation, num int) {
	// fmt.Printf("\tsetNeighborCount(%s, %d)\n", organism.String(), num)
	originalNum := t.GetOrganismValue(organism)

	// Write the value to the gameboard
	t.gameboard.SetValue(organism, num)

	// Update the living count if organism changed living state
	if originalNum < 0 && num >= 0 {
		t.living.Set(organism)
	} else if originalNum >= 0 && num < 0 {
		t.living.Remove(organism)
	}
}

func (t *Pond) calculateNeighborCount(organism GameboardLocation) (int, []GameboardLocation) {
	numNeighbors := 0
	neighbors := t.GetNeighbors(organism)
	for _, neighbor := range neighbors {
		if t.isOrganismAlive(neighbor) {
			numNeighbors++
		}
	}
	return numNeighbors, neighbors
}

func (t *Pond) SetOrganisms(organisms []GameboardLocation) {
	// Initialize the first organisms and set their neighbor counts
	for _, organism := range organisms {
		t.setOrganismValue(organism, 0)
	}
}

func (t *Pond) GetGameboard() [][]int {
	return t.gameboard.getSnapshot()
}

func (t *Pond) Clone() *Pond {
	shadowPond := NewPond(t.gameboard.Dims.Height,
		t.gameboard.Dims.Width,
		t.neighborsSelector)

	shadowPond.Status = t.Status
	shadowPond.neighborsSelector = t.neighborsSelector

	shadowPond.SetOrganisms(t.living.GetAll())

	return shadowPond
}

func (t *Pond) Equals(rhs *Pond) bool {
	if !t.gameboard.Equals(rhs.gameboard) {
		return false
	}
	if t.Status != rhs.Status {
		return false
	}
	if t.neighborsSelector != rhs.neighborsSelector {
		return false
	}
	return true
}

func (t *Pond) String() string {
	var buf bytes.Buffer
	buf.WriteString("Neighbor selection: ")
	buf.WriteString(t.neighborsSelector.String())
	buf.WriteString("\nLiving organisms: ")
	buf.WriteString(strconv.Itoa(t.living.GetCount()))
	buf.WriteString("\tStatus: ")
	buf.WriteString(t.Status.String())
	buf.WriteString("\n")
	buf.WriteString(t.gameboard.String())

	return buf.String()
}

func NewPond(rows int, cols int, neighbors NeighborsSelector) *Pond {
	p := new(Pond)

	// Create values
	p.Status = Active
	p.living = NewLivingTracker()

	// Add the given values
	var err error
	p.gameboard, err = NewGameboard(GameboardDims{Height: rows, Width: cols})
	if err != nil {
		// TODO: return nil,errors.New(blah)
		// t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	p.neighborsSelector = neighbors

	return p
}
