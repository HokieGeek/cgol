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

type Life struct {
	Label            string
	Statistics       LifeStats
	UpdateRate       time.Duration
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
	} else if organismsDelta < 0 {
		t.Statistics.OrganismsKilled += (organismsDelta * -1)
	}
	// }

	// If the pond is dead, let's just stop doing things
	if t.pond.Status == Dead {
		t.Stop()
	}
}

// func (t *Life) Start(updateAlert chan bool, updateRate time.Duration) {
func (t *Life) Start(updateAlert chan bool) {
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

	// Initialize the pond and schedule the currently living organisms
	// s.initialOrganisms = append(s.initialOrganisms, initializer(s.pond.board.Dims)...)
	s.initialOrganisms = initializer(s.pond.board.Dims)
	s.pond.SetOrganisms(s.initialOrganisms)
	s.Statistics.OrganismsCreated = len(s.initialOrganisms)

	return s, nil
}
