package life

import (
	"bytes"
	// "crypto/sha1"
	"fmt"
	// "strconv"
)

type ChangeType int

const (
	Born ChangeType = iota
	Died
)

type ChangedLocation struct {
	Location
	Change ChangeType
	// PatternGroup ...
	// Classificaiton ...
}

type Analysis struct {
	Status  Status
	Living  []Location
	Changes []ChangedLocation
	// TODO: checksum []byte
}

// type (t *Analysis) Checksum() [sha1.Size]byte {
// var str bytes.Buffer
// str.WriteString(strconv.Itoa(t.Generations))

// h := sha1.New()
// buf := make([]byte, sha1.Size)
// h.Write(buf)
// return h.Sum(nil)
// }

type Analyzer struct {
	Id           []byte
	Life         *Life
	analyses     []Analysis // Each index is a generation
	stopAnalysis func()
}

func (t *Analyzer) Analysis(generation int) *Analysis {
	if generation < 0 || generation >= len(t.analyses) {
		// TODO: maybe an error
		return nil
	}
	return &t.analyses[generation]
}

func (t *Analyzer) analyze(cells []Location, generation int) {
	if len(cells) <= 0 {
		fmt.Println("WTF")
	}
	var analysis Analysis

	// Record the status
	// analysis.Status =

	// Copy the living cells
	analysis.Living = make([]Location, len(cells))
	copy(analysis.Living, cells)

	// Initialize and start processing the living cells
	analysis.Changes = make([]ChangedLocation, 0)

	if generation <= 0 { // Special case to reduce code duplication
		for _, loc := range cells {
			analysis.Changes = append(analysis.Changes, ChangedLocation{Location: loc, Change: Born})
		}
	} else {
		// Add any new cells
		previousLiving := t.analyses[generation-1].Living
		for _, newCell := range cells {
			found := false
			for _, oldCell := range previousLiving {
				if oldCell.Equals(&newCell) {
					found = true
					break
				}
			}

			if !found {
				analysis.Changes = append(analysis.Changes, ChangedLocation{Location: newCell, Change: Born})
			}
		}

		// Add any cells which died
		for _, oldCell := range previousLiving {
			found := false
			for _, newCell := range cells {
				if newCell.Equals(&oldCell) {
					found = true
					break
				}
			}

			if !found {
				analysis.Changes = append(analysis.Changes, ChangedLocation{Location: oldCell, Change: Died})
			}
		}

	}

	t.analyses = append(t.analyses, analysis)
}

func (t *Analyzer) NumAnalyses() int {
	return len(t.analyses)
}

func (t *Analyzer) Start() {
	updates := make(chan bool)
	t.stopAnalysis = t.Life.Start(updates, -1)

	go func() {
		for {
			select {
			case <-updates:
				nextGen := len(t.analyses)
				gen := t.Life.Generation(nextGen)
				fmt.Printf("Generation %d\n", gen.Num)
				fmt.Println(t.Life)
				t.analyze(gen.Living, gen.Num)
			}
		}
	}()
}

func (t *Analyzer) Stop() {
	t.stopAnalysis()
}

func (t *Analyzer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%x", t.Id))
	buf.WriteString("\n")
	buf.WriteString(t.Life.String())

	return buf.String()
}

func NewAnalyzer(dims Dimensions, pattern func(Dimensions, Location) []Location, rulesTester func(int, bool) bool) (*Analyzer, error) {
	// fmt.Printf("NewAnalyzer: %v\n", pattern(dims, Location{X: 0, Y: 0}))
	a := new(Analyzer)

	var err error
	a.Life, err = New("HTTP REQUEST",
		dims,
		NEIGHBORS_ALL,
		pattern,
		rulesTester,
		SimultaneousProcessor)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	// fmt.Println("Creating unique id")
	a.Id = uniqueId()

	// Generate first analysis (for generation 0 / the seed)
	a.analyze(a.Life.Seed, 0)

	return a, nil
}
