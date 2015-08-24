package cgol

import "fmt"

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

type OrganismReference struct {
	X int
	Y int
	// Z int
}

type Pond struct {
	Name        string
	Rows        int
	Cols        int
	NumLiving   int
	Status      GameStatus
	gameboard   [][]int
	living      []OrganismReference // Would be best if this was a tuple :-/...
	ruleset     func(*Pond, OrganismReference)
	initializer func(*Pond)
}

func (t *Pond) getOrganism(cell OrganismReference) *int {
	return &t.gameboard[cell.X][cell.Y]
}

func (t *Pond) updateNeighborCount(organism OrganismReference, delta int) {
	cell := t.getOrganism(organism)
	(*cell) += delta
}

func (t *Pond) updateOrganismLivingState(organism OrganismReference, live bool) {
	// cell := t.getOrganism(organism)
	if live {
		t.gameboard[organism.X][organism.Y] = 0
		t.NumLiving++
		t.living = append(t.living, organism)
	} else {
		t.gameboard[organism.X][organism.Y] = -1
		t.NumLiving--
	}
}

// func (t* Pond) applyRuleset() {
// }

func (t *Pond) start() {
	// TODO: Kick off thread that periodically goes through the 'living' queue
}

func (t *Pond) stop() {
	// TODO: stop the processing thread
}

func (t *Pond) Display() {
	fmt.Printf("[%s] (%dx%d)\n", t.Name, t.Rows, t.Cols)
	fmt.Printf("Living organisms: %d\tStatus: %s\n", t.NumLiving, t.Status)
	for i := 0; i < t.Rows; i++ {
		for j := 0; j < t.Cols; j++ {
			if t.gameboard[i][j] >= 0 {
				fmt.Printf("%d", t.gameboard[i][j])
			} else {
				fmt.Printf("-")
			}
		}
		fmt.Printf("\n")
	}
}

func CreatePond(name string, rows int, cols int, rules func(*Pond, OrganismReference), init func(*Pond)) *Pond {
	p := new(Pond)
	p.Name = name
	p.Rows = rows
	p.Cols = cols
	p.ruleset = rules
	p.initializer = init
	p.Status = Active

	p.initializer(p)
	// p.blah blah blah

	return p
}
