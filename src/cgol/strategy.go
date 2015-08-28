package cgol

import (
	"bytes"
	"fmt"
	"time"
)

type StrategyStats struct {
	OrganismsCreated int
	OrganismsKilled  int
	Iterations       int
}

type Strategy struct {
	Label      string
	Statistics StrategyStats
	pond       *Pond
	processor  func(pond *Pond, rules func(int, bool) bool)
	ruleset    func(int, bool) bool
	ticker     *time.Ticker
	updateRate time.Duration
}

// FIXME: this method shouldn't exist at all, really
func (t *Strategy) process() {
	// TODO: if have been stable for a while stop processing
	// TODO: does this stuff belong here or in the processor?
	// startingLivingCount := t.pond.NumLiving

	// Process any organisms that need to be
	t.processor(t.pond, t.ruleset)

	// Update the pond's statistics
	// if stillProcessing {
	// 	t.Statistics.Iterations++

	// 	// Update the statistics
	// 	organismsDelta := t.pond.NumLiving - startingLivingCount
	// 	if organismsDelta > 0 {
	// 		t.Statistics.OrganismsCreated += organismsDelta
	// 	} else if organismsDelta < 0 {
	// 		t.Statistics.OrganismsKilled += (organismsDelta * -1) // FIXME
	// 	}
	// }

	// If the pond is dead, let's just stop doing things
	if t.pond.Status == Dead {
		t.Stop()
	}
}

func (t *Strategy) Start() {
	// t.processor.Schedule(t.pond.initialOrganisms)
	// t.processor.Process(t.pond, t.ruleset)

	t.ticker = time.NewTicker(t.updateRate)
	for {
		select {
		case <-t.ticker.C:
			t.process()
			fmt.Println(t) // TODO: remove
		}
	}
}

func (t *Strategy) Stop() {
	// t.processor.Stop()
	t.ticker.Stop()
}

func (t *Strategy) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(t.Label)
	buf.WriteString("]\n")
	// TODO: fmt.Printf("Ruleset: %s\n", t.ruleset)
	buf.WriteString(t.pond.String())

	return buf.String()
}

func NewStrategy(label string,
	pond *Pond,
	initializer func(*Pond) []GameboardLocation,
	rules func(int, bool) bool,
	processor func(pond *Pond, rules func(int, bool) bool)) *Strategy {
	s := new(Strategy)

	// Save the given values
	s.Label = label
	s.pond = pond
	s.ruleset = rules
	s.processor = processor

	s.updateRate = time.Millisecond * 250

	// Initialize the pond and schedule the currently living organisms
	initialLiving := initializer(s.pond)
	s.pond.init(initialLiving)
	s.Statistics.OrganismsCreated = len(initialLiving)

	return s
}
