package cgol

import (
	"bytes"
	"strconv"
	"time"
)

type StrategyStats struct {
	OrganismsCreated int
	OrganismsKilled  int
	Generations      int
}

func (t *StrategyStats) String() string {
	var buf bytes.Buffer
	buf.WriteString("Generation: ")
	buf.WriteString(strconv.Itoa(t.Generations))

	return buf.String()
}

type Strategy struct {
	Label            string
	Statistics       StrategyStats
	UpdateRate       time.Duration
	pond             *Pond
	processor        func(pond *Pond, rules func(int, bool) bool)
	ruleset          func(int, bool) bool
	initialOrganisms []GameboardLocation
	ticker           *time.Ticker
}

func (t *Strategy) process() {
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

// func (t *Strategy) Start(updateAlert chan bool, updateRate time.Duration) {
func (t *Strategy) Start(updateAlert chan bool) {
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

func (t *Strategy) Stop() {
	t.ticker.Stop()
}

/*
func (t *Strategy) GetGameboard() [][]int {
	return t.pond.GetGameboard()
}
*/

type Generation struct {
	Num    int
	Living []GameboardLocation
	// stats...
}

func (t *Strategy) GetGeneration(num int) *Generation {
	var p *Pond
	if num == t.Statistics.Generations {
		p = t.pond
	} else {
		cloned := t.pond.Clone()
		cloned.SetOrganisms(t.initialOrganisms)
		for i := 0; i < num; i++ {
			t.processor(cloned, t.ruleset)
		}

		p = cloned
	}

	return &Generation{Num: num, Living: p.living.GetAll()}
}

func (t *Strategy) String() string {
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

func NewStrategy(label string,
	pond *Pond,
	initializer func(GameboardDims) []GameboardLocation,
	rules func(int, bool) bool,
	processor func(pond *Pond, rules func(int, bool) bool)) *Strategy {
	s := new(Strategy)

	// Save the given values
	s.Label = label
	s.pond = pond
	s.ruleset = rules
	s.processor = processor

	s.UpdateRate = time.Millisecond * 250

	// Initialize the pond and schedule the currently living organisms
	// s.initialOrganisms = append(s.initialOrganisms, initializer(s.pond.gameboard.Dims)...)
	s.initialOrganisms = initializer(s.pond.gameboard.Dims)
	s.pond.SetOrganisms(s.initialOrganisms)
	s.Statistics.OrganismsCreated = len(s.initialOrganisms)

	return s
}
