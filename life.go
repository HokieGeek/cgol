package life

import (
	"bytes"
)

type Generation struct {
	Num    int
	Living []Location
}

type Life struct {
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

	return &Generation{Num: t.Generations, Living: t.pond.living.GetAll()}
}

func (t *Life) Start(listener chan *Generation) func() {
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

	// Initialize the pond and schedule the currently living organisms
	s.Seed = initializer(s.pond.board.Dims, Location{})
	s.pond.SetOrganisms(s.Seed)

	return s, nil
}

// vim: set foldmethod=marker:
