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
	processor  Processor
	ruleset    func(*Pond, OrganismReference) bool
	ticker     *time.Ticker
	updateRate time.Duration
}

func (t *Strategy) process() {
	currentLivingCount := t.pond.NumLiving

	stillProcessing := false
	if currentLivingCount > 3 { // Only process the current queue if it makes sense to do so
		stillProcessing = t.processor.Process(t.pond, t.ruleset)
		// } else {
		// TODO: t.processor.clear()
	}

	// Update the pond's status
	if stillProcessing {
		t.pond.Status = Active
		t.Statistics.Iterations++

		// Update the statistics
		organismsDelta := t.pond.NumLiving - currentLivingCount
		if organismsDelta > 0 {
			t.Statistics.OrganismsCreated += organismsDelta
		} else if organismsDelta < 0 {
			t.Statistics.OrganismsKilled += (organismsDelta * -1) // FIXME
		}
	} else {
		if t.pond.NumLiving > 0 {
			t.pond.Status = Stable
			// TODO: if have been stable for a while stop processing
		} else {
			t.pond.Status = Dead
			t.Stop()
		}
	}
}

func (t *Strategy) Start() {
	t.ticker = time.NewTicker(t.updateRate)
	for {
		select {
		case <-t.ticker.C:
			t.process()
			fmt.Println(t)
		}
	}
}

func (t *Strategy) Stop() {
	t.ticker.Stop()
}

func (t *Strategy) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(t.Label)
	buf.WriteString("]\n")
	// fmt.Printf("Ruleset: %s\n", t.ruleset)
	buf.WriteString(t.pond.String())

	return buf.String()
}

func (t *Strategy) init(initializer func(*Pond) []OrganismReference) {
	// Initialize the pond
	initialLiving := initializer(t.pond)
	// TODO: Try several times if initialLiving == 0
	t.pond.init(initialLiving)

	// Initialize the statistics tracker
	t.Statistics.Iterations = 0
	t.Statistics.OrganismsCreated = len(initialLiving)
	t.Statistics.OrganismsKilled = 0

	// Schedule the currently living organisms
	t.processor.schedule(initialLiving)
}

func CreateStrategy(label string,
	pond *Pond,
	initializer func(*Pond) []OrganismReference,
	rules func(*Pond, OrganismReference) bool,
	processor Processor) *Strategy {
	s := new(Strategy)

	// Save the given values
	s.Label = label
	s.pond = pond
	s.ruleset = rules
	s.processor = processor

	s.updateRate = time.Millisecond * 500

	// Prime the pump
	s.init(initializer)

	return s
}
