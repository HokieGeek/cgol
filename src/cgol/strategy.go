package cgol

import (
    "fmt"
    "time"
)

type StrategyStats struct {
    OrganismsCreated int
    OrganismsKilled int
    Iterations int
}

type Strategy struct {
	Label      string
    Statistics StrategyStats
	pond       *Pond
	processor  Processor
	ruleset    func(*Pond, OrganismReference) bool
    ticker     Ticker
    updateRate Duration
}

func (t *Strategy) process() {
    currentLivingCount = t.pond.NumLiving

	stillProcessing := t.processor.Process(t.pond, t.ruleset)

    // Update the pond's status
    if stillProcessing {
        t.pond.Status = Active
        t.Statistics.Iterations++

        // Update the statistics
        organismsDelta := t.pond.NumLiving - currentLivingCount
        if organismsDelta > 0
            t.Statistics.OrganismsCreated += organismsDelta
        } else if organismsDelta < 0 {
            t.Statistics.OrganismsKilled += (organismsDelta * -1) // FIXME
        }
    } else {
        if t.pond.NumLiving > 0 {
            t.pond.Status = Stable
        } else {
            t.pond.Status = Dead
            t.Stop()
        }
    }
}

func (t *Strategy) Start() {
    // ticker := time.NewTicker(t.updateRate)
    // go func() {
    //     for u := range ticker.C {
	        t.process()
    //     }
    // } ()
}

func (t *Strategy) Stop() {
	// TODO: stop the processing thread
    ticker.Stop()
}

func (t *Strategy) Display() {
	fmt.Printf("[%s]\n", t.Label)
    fmt.Printf("Ruleset: %s\n", t.ruleset) // TODO: is this a good idea?
	t.pond.Display()
}

func (t *Strategy) init(initializer func(*Pond) []OrganismReference) {
    // Initialize the pond
    initialLiving = initializer(t.pond)
    // TODO: Try several times if initialLiving == 0
	t.pond.init(initialLiving)

    // Initialize the statistics tracker
    t.Statistics.iterations = 0
    t.Statistics.OrganismsCreated = len(initialLiving)
    t.Statistics.OrganismsKilled = 0

    // Schedule the currently living organisms
    t.processor.schedule(initialLiving)
}

func CreateStrategy(label string,
	pond *Pond,
    initializer func(*Pond) []OrganismReference,
	rules func(*Pond, OrganismReference) bool,
	processor Processor,
    rate Duration) *Strategy {
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
