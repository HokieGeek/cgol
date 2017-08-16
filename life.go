package life

import (
	"bytes"
)

type Status int

const (
	Seeded Status = iota
	Active
	Stable
	Dead
)

func (t Status) String() string {
	switch t {
	case Seeded:
		return "Seeded"
	case Active:
		return "Active"
	case Dead:
		return "Dead"
	}

	return "Unknown"
}

type Life struct {
	Status      Status
	pond        *pond
	processor   func(pond *pond, rules func(int, bool) bool)
	ruleset     func(int, bool) bool
	Seed        []Location
	Generations int
}

func (t *Life) process() *Generation {
	// Process any organisms that need to be
	t.processor(t.pond, t.ruleset)

	// Update the pond's statistics
	t.Generations++

	// If we don't have living cells, then we are done.
	if t.pond.GetNumLiving() > 0 {
		t.Status = Active
	} else {
		t.Status = Dead
	}

	return &Generation{Num: t.Generations, Living: t.pond.living.GetAll()}
}

func (t *Life) Start(listener chan *Generation) func() {
	t.Status = Active

	stop := false

	go func() {
		for {
			if stop {
				break
			} else {
				if listener != nil {
					listener <- t.process()
				}
			}
		}
	}()

	return func() {
		stop = true
	}
}

type Generation struct {
	Num    int
	Living []Location
}

func (t *Life) Generation(num int) *Generation {
	var p *pond
	if num == t.Generations {
		p = t.pond
	} else {
		cloned, err := t.pond.Clone()
		if err != nil {
			// logf("Unable to clone pond: %s\n", err)
			return nil // FIXME
		}
		cloned.SetOrganisms(t.Seed)
		for i := 0; i < num; i++ {
			t.processor(cloned, t.ruleset)
		}

		p = cloned
	}

	return &Generation{Num: num, Living: p.living.GetAll()}
}

func (t *Life) Dimensions() Dimensions {
	return t.pond.board.Dims
}

func (t *Life) String() string {
	var buf bytes.Buffer

	// buf.WriteString("Status: ")
	// buf.WriteString(t.Status.String())
	buf.WriteString("\n")
	buf.WriteString(t.pond.String())

	return buf.String()
}

func New(dims Dimensions,
	neighbors neighborsSelector,
	initializer func(Dimensions, Location) []Location,
	rules func(int, bool) bool,
	processor func(pond *pond, rules func(int, bool) bool)) (*Life, error) {
	s := new(Life)

	var err error
	board, err := newBoard(dims)
	if err != nil {
		return nil, err
	}
	livingTracker := newTracker()
	s.pond, err = newPond(board, livingTracker, neighbors)
	if err != nil {
		return nil, err
	}

	// Save the given values
	s.ruleset = rules
	s.processor = processor

	s.Status = Seeded

	// Initialize the pond and schedule the currently living organisms
	s.Seed = initializer(s.pond.board.Dims, Location{})
	s.pond.SetOrganisms(s.Seed)

	return s, nil
}

// vim: set foldmethod=marker:
