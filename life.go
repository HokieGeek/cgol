package life

import (
	"bytes"
	"strconv"
	"time"
)

type Statistics struct {
	OrganismsCreated int
	OrganismsKilled  int
	Generations      int
}

func (t *Statistics) String() string {
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
	switch t {
	case Seeded:
		return "Seeded"
	case Active:
		return "Active"
	case Stable:
		return "Stable"
	case Dead:
		return "Dead"
	}

	return "Unknown"
}

type Life struct {
	Label            string
	Stats            Statistics
	Status           Status
	pond             *pond
	processor        func(pond *pond, rules func(int, bool) bool)
	ruleset          func(int, bool) bool
	initialOrganisms []Location
}

func (t *Life) process() {
	startingLivingCount := t.pond.GetNumLiving()

	// Process any organisms that need to be
	t.processor(t.pond, t.ruleset)

	// Update the pond's statistics
	t.Stats.Generations++

	// Update the statistics
	organismsDelta := t.pond.GetNumLiving() - startingLivingCount
	if organismsDelta > 0 {
		t.Stats.OrganismsCreated += organismsDelta
		t.Status = Active
	} else if organismsDelta < 0 {
		t.Stats.OrganismsKilled += (organismsDelta * -1)
		t.Status = Active
	} else {
		t.Status = Stable
	}

	// If the pond is dead, let's just stop doing things
	if t.pond.GetNumLiving() <= 0 {
		t.Status = Dead
	}
}

func (t *Life) Start(alert chan bool, rate time.Duration) func() {
	t.Status = Active

	if rate > 0 {
		ticker := time.NewTicker(rate)

		go func() {
			for {
				select {
				case <-ticker.C:
					t.process()
					alert <- true
				}
			}
		}()

		return func() {
			ticker.Stop()
		}
	} else {
		stop := false

		go func() {
			for {
				if stop {
					break
				} else {
					t.process()
					alert <- true
				}
			}
		}()

		return func() {
			stop = true
		}
	}
}

type Generation struct {
	Num    int
	Living []Location
	// stats...
}

func (t *Life) Generation(num int) *Generation {
	var p *pond
	if num == t.Stats.Generations {
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
	buf.WriteString("Status: ")
	buf.WriteString(t.Status.String())
	buf.WriteString("\tGeneration: ")
	buf.WriteString(t.Stats.String())
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

	s.Status = Seeded

	// Initialize the pond and schedule the currently living organisms
	s.initialOrganisms = initializer(s.pond.board.Dims)
	s.pond.SetOrganisms(s.initialOrganisms)
	s.Stats.OrganismsCreated = len(s.initialOrganisms)

	return s, nil
}
