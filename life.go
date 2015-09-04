package life

import (
	"bytes"
	"strconv"
	"time"
)

type LifeStats struct {
	OrganismsCreated int
	OrganismsKilled  int
	Generations      int
}

func (t *LifeStats) String() string {
	var buf bytes.Buffer
	buf.WriteString("Generation: ")
	buf.WriteString(strconv.Itoa(t.Generations))

	return buf.String()
}

type Status int

const (
	Seeded Status = 0
	Active Status = 1
	Stable Status = 2
	Dead   Status = 3
)

func (t Status) String() string {
	s := ""

	if t&Seeded == Seeded {
		s += "Seeded"
	} else if t&Active == Active {
		s += "Active"
	} else if t&Stable == Stable {
		s += "Stable"
	} else if t&Dead == Dead {
		s += "Dead"
	}

	return s
}

type Life struct {
	Label            string
	Statistics       LifeStats
	UpdateRate       time.Duration
	Status           Status
	pond             *pond
	processor        func(pond *pond, rules func(int, bool) bool)
	ruleset          func(int, bool) bool
	initialOrganisms []Location
	ticker           *time.Ticker
}

func (t *Life) process() {
	startingLivingCount := t.pond.GetNumLiving()

	// Process any organisms that need to be
	t.processor(t.pond, t.ruleset)

	// Update the pond's statistics
	// if stillProcessing {
	t.Statistics.Generations++

	// Update the statistics
	organismsDelta := t.pond.GetNumLiving() - startingLivingCount
	if organismsDelta > 0 {
		t.Statistics.OrganismsCreated += organismsDelta
		t.Status = Active
	} else if organismsDelta < 0 {
		t.Statistics.OrganismsKilled += (organismsDelta * -1)
		t.Status = Active
	} else {
		t.Status = Stable
	}

	// }

	// If the pond is dead, let's just stop doing things
	if t.pond.GetNumLiving() <= 0 {
		t.Status = Dead
		t.Stop()
	}
}

// func (t *Life) Start(updateAlert chan bool, updateRate time.Duration) {
func (t *Life) Start(updateAlert chan bool) {
	t.Status = Active
	go func() {
		t.ticker = time.NewTicker(t.UpdateRate)
		/*
			        if updateRate > 0 {
				    	for {
			                t.process()
			                updateAlert <- true
			                // TODO: need to figure out a way to stop this
				    		}
				    	}
			        } else {
		*/
		for {
			select {
			case <-t.ticker.C:
				t.process()
				updateAlert <- true
			}
		}
		// }
	}()
}

func (t *Life) Stop() {
	t.ticker.Stop()
}

/*
func (t *Life) GetLifeboard() [][]int {
	return t.pond.GetLifeboard()
}
*/

type Generation struct {
	Num    int
	Living []Location
	// stats...
}

func (t *Life) Generation(num int) *Generation {
	var p *pond
	if num == t.Statistics.Generations {
		p = t.pond
	} else {
		cloned, err := t.pond.Clone()
		if err != nil {
			// logf("Unable to clone pond: %s\n", err)
			return nil // FIXME
		}
		cloned.SetOrganisms(t.initialOrganisms)
		for i := 0; i < num; i++ {
			t.processor(cloned, t.ruleset)
		}

		p = cloned
	}

	return &Generation{Num: num, Living: p.living.GetAll()}
}

func (t *Life) String() string {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(t.Label)
	buf.WriteString("]\n")
	buf.WriteString("\tStatus: ")
	buf.WriteString(t.Status.String())
	buf.WriteString("Generation: ")
	buf.WriteString(t.Statistics.String())
	buf.WriteString("\n")
	buf.WriteString(t.pond.String())

	return buf.String()
}

func New(label string,
	dims Dimensions,
	neighbors neighborsSelector,
	initializer func(Dimensions) []Location,
	rules func(int, bool) bool,
	processor func(pond *pond, rules func(int, bool) bool)) (*Life, error) {
	s := new(Life)

	var err error
	s.pond, err = newpond(Dimensions{Height: dims.Height, Width: dims.Width}, neighbors)
	if err != nil {
		return nil, err
	}

	// Save the given values
	s.Label = label
	s.ruleset = rules
	s.processor = processor

	s.UpdateRate = time.Millisecond * 250
	s.Status = Seeded

	// Initialize the pond and schedule the currently living organisms
	// s.initialOrganisms = append(s.initialOrganisms, initializer(s.pond.board.Dims)...)
	s.initialOrganisms = initializer(s.pond.board.Dims)
	s.pond.SetOrganisms(s.initialOrganisms)
	s.Statistics.OrganismsCreated = len(s.initialOrganisms)

	return s, nil
}
