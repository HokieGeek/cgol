package cgol

import (
	"bytes"
	// "fmt"
	"strconv"
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

type LivingTracker struct {
	trackerAdd    chan *livingTrackerAddOp
	trackerRemove chan *livingTrackerRemoveOp
	trackerTest   chan *livingTrackerTestOp
	trackerGetAll chan *livingTrackerGetAllOp
}

func (t *LivingTracker) living() {
	var livingMap = make(map[int]map[int]GameboardLocation)

	for {
		select {
		case add := <-t.trackerAdd:
			_, keyExists := livingMap[add.loc.Y]
			if !keyExists {
				livingMap[add.loc.Y] = make(map[int]GameboardLocation)
			}
			livingMap[add.loc.Y][add.loc.X] = add.loc
			add.resp <- true
		case remove := <-t.trackerRemove:
			_, keyExists := livingMap[remove.loc.Y]
			if keyExists {
				_, keyExists = livingMap[remove.loc.Y][remove.loc.X]
				if keyExists {
					delete(livingMap[remove.loc.Y], remove.loc.X)
					if len(livingMap[remove.loc.Y]) <= 0 {
						delete(livingMap, remove.loc.Y)
					}
				}
			}
			remove.resp <- true
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

func NewLivingTracker() *LivingTracker {
	t := new(LivingTracker)

	t.trackerAdd = make(chan *livingTrackerAddOp)
	t.trackerRemove = make(chan *livingTrackerRemoveOp)
	t.trackerTest = make(chan *livingTrackerTestOp)
	t.trackerGetAll = make(chan *livingTrackerGetAllOp)

	go t.living()

	return t
}

type Pond struct {
	gameboard         *Gameboard
	NumLiving         int
	Status            PondStatus
	neighborsSelector NeighborsSelector
	Living            map[int]map[int]GameboardLocation
	// living            *LivingTracker
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
	return len(t.Living)
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
		// TODO: add to 'living'
		t.NumLiving++
	} else if originalNum >= 0 && num < 0 {
		t.NumLiving--
		// TODO: remove from 'living'
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

func (t *Pond) init(initialLiving []GameboardLocation) {
	// Initialize the first organisms and set their neighbor counts
	t.Living = make(map[int]map[int]GameboardLocation)
	for _, organism := range initialLiving {
		// TODO: this logic needs to move into its own place function with channel accessors
		_, keyExists := t.Living[organism.Y]
		if !keyExists {
			t.Living[organism.Y] = make(map[int]GameboardLocation)
		}
		t.Living[organism.Y][organism.X] = organism
		t.setOrganismValue(organism, 0)
	}
}

func (t *Pond) Clone() *Pond {
	shadowPond := NewPond(t.gameboard.Dims.Height,
		t.gameboard.Dims.Width,
		t.neighborsSelector)

	shadowPond.NumLiving = t.NumLiving
	shadowPond.Status = t.Status

	// TODO
	// shadowPond.init(t.Living)
	// Living            map[int]map[int]GameboardLocation

	return shadowPond
}

func (t *Pond) String() string {
	var buf bytes.Buffer
	buf.WriteString("Neighbor selection: ")
	buf.WriteString(t.neighborsSelector.String())
	buf.WriteString("\nLiving organisms: ")
	buf.WriteString(strconv.Itoa(t.NumLiving))
	buf.WriteString("\tStatus: ")
	buf.WriteString(t.Status.String())
	buf.WriteString("\n")
	buf.WriteString(t.gameboard.String())

	return buf.String()
}

func NewPond(rows int, cols int, neighbors NeighborsSelector) *Pond {
	p := new(Pond)

	// Create values
	p.NumLiving = 0
	p.Status = Active

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
