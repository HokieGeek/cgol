package life

import (
	"bytes"
	// "crypto/sha1"
	"fmt"
	// "strconv"
)

type ChangeType int

const (
	Born ChangeType = 0
	Died ChangeType = 1
)

type ChangedLocation struct {
	Location
	Change ChangeType
	// PatternGroup ...
	// Classificaiton ...
}

type Analysis struct {
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
	Id       []byte
	Life     *Life
	analyses []Analysis // Each index is a generation
}

func (t *Analyzer) Analysis(generation int) *Analysis {
	// TODO: input validation
	// fmt.Printf("Analysis(%d)\n", generation)
	// fmt.Printf("Num analyses: %d\n", len(t.analyses))
	// fmt.Printf("Num analyses.Living: %d\n", len(t.analyses[generation].Living))
	// fmt.Printf("Num analyses.Changes: %d\n", len(t.analyses[generation].Changes))
	return &t.analyses[generation]
}

func (t *Analyzer) analyze(cells []Location, generation int) {
	var analysis Analysis

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
func (t *Analyzer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%x", t.Id))
	buf.WriteString("\n")
	buf.WriteString(t.Life.String())

	return buf.String()
}

func NewAnalyzer(dims Dimensions) (*Analyzer, error) {
	// fmt.Println("NewAnalyzer")
	a := new(Analyzer)

	var err error
	a.Life, err = New("HTTP REQUEST",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	// fmt.Println("Creating unique id")
	a.Id = uniqueId()

	// Generate first analysis (for generation 0 / the seed)
	a.analyze(a.Life.Seed, 0)

	/*
		// Generate first analysis (for generation 0 / the seed)
		var seedAnalysis Analysis
		seedAnalysis.Living = make([]Location, len(a.Life.Seed))
		copy(seedAnalysis.Living, a.Life.Seed)

		seedAnalysis.Changes = make([]ChangedLocation, 0)
		for _, loc := range a.Life.Seed {
			seedAnalysis.Changes = append(seedAnalysis.Changes, ChangedLocation{Location: loc, Change: Born})
		}

		a.analyses = make([]Analysis, 0)
		a.analyses = append(a.analyses, seedAnalysis)
	*/

	gen := a.Life.Generation(1)
	a.analyze(gen.Living, 1)

	return a, nil
}
