package life

import (
	"bytes"
	"fmt"
)

type Analyzer struct {
	Id   []byte
	Life *Life
	// Stats
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

	return a, nil
}
