package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.com/hokiegeek/biologist"
	"gitlab.com/hokiegeek/life"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

/////////////////////////////////// CREATE ANALYSIS ///////////////////////////////////

type CreateAnalysisResponse struct {
	Id   []byte
	Dims life.Dimensions
	// Rule string
	// Neighbors  life.NeighborsSelector
}

func NewCreateAnalysisResponse(analyzer *biologist.Analyzer) *CreateAnalysisResponse {
	resp := new(CreateAnalysisResponse)

	resp.Id = analyzer.Id
	resp.Dims = analyzer.Life.Dimensions()
	// resp.Rule = analyzer.Generation()

	return resp
}

type PatternType int

const (
	USER PatternType = iota
	RANDOM
	BLINKERS
	TOADS
	BEACONS
	PULSARS
	GLIDERS
	BLOCKS
	BEEHIVES
	LOAVES
	BOATS
)

type CreateAnalysisRequest struct {
	Dims    life.Dimensions
	Pattern PatternType
	Seed    []life.Location
	// life.Rules
	// Processor
}

func (t *CreateAnalysisRequest) String() string {
	var buf bytes.Buffer

	buf.WriteString(t.Dims.String())

	return buf.String()
}

func CreateAnalysis(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var req CreateAnalysisRequest
	fmt.Printf("REQ: %s\n", body)
	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		// FIXME: this should be sent to a logger
		fmt.Printf("Received create request: %s\n", req.String())

		// Determine the pattern to use for seeding the board
		var patternFunc func(life.Dimensions, life.Location) []life.Location
		switch req.Pattern {
		case USER:
			// fmt.Printf("Created USER pattern func: %v\n", req.Seed)
			patternFunc = func(dims life.Dimensions, offset life.Location) []life.Location {
				// fmt.Println("HERE I AM")
				return req.Seed
			}
		case RANDOM:
			patternFunc = func(dims life.Dimensions, offset life.Location) []life.Location {
				return life.Random(dims, offset, 35)
			}
		case BLINKERS:
			patternFunc = life.Blinkers
		case PULSARS:
			patternFunc = life.Pulsar
		case GLIDERS:
			patternFunc = life.Gliders
		case BLOCKS:
			patternFunc = life.Blocks
		}

		// Create the analyzer
		// fmt.Printf("Creating new analyzer with pattern: %v\n", patternFunc(req.Dims, life.Location{X: 0, Y: 0}))
		analyzer, err := biologist.NewAnalyzer(req.Dims, patternFunc, life.ConwayTester())
		if err != nil {
			panic(err)
		}
		mgr.Add(analyzer)
		fmt.Println(analyzer)

		// fmt.Printf("Id: %x\n", analyzer.Id)

		// Respond the request with the ID of the analyzer
		resp := NewCreateAnalysisResponse(analyzer)

		postJson(w, http.StatusCreated, resp)
	}
}

/////////////////////////////////// UPDATE ANALYSIS ///////////////////////////////////

type AnalysisUpdate struct {
	Id         []byte
	Dims       life.Dimensions
	Generation int
	Living     []life.Location
	Changes    []biologist.ChangedLocation
	// Neighbors life.NeighborSelector
}

func NewAnalysisUpdate(analyzer *biologist.Analyzer, generation int) *AnalysisUpdate {
	fmt.Printf(" NewAnalysisUpdate(%d)\n", generation)
	a := new(AnalysisUpdate)

	a.Id = analyzer.Id
	a.Dims = analyzer.Life.Dimensions()
	a.Generation = generation

	analyzer.Analysis(generation)
	analysis := analyzer.Analysis(generation)

	a.Living = make([]life.Location, len(analysis.Living))
	copy(a.Living, analysis.Living)

	a.Changes = make([]biologist.ChangedLocation, len(analysis.Changes))
	copy(a.Changes, analysis.Changes)

	return a
}

type AnalysisUpdateRequest struct {
	Id                 []byte
	StartingGeneration int
	NumMaxGenerations  int
}

func (t *AnalysisUpdateRequest) String() string {
	var buf bytes.Buffer

	buf.WriteString("Id: ")
	buf.WriteString(fmt.Sprintf("%x", t.Id))
	buf.WriteString("\nStarting Generation: ")
	buf.WriteString(strconv.Itoa(t.StartingGeneration))
	buf.WriteString("\nMax: ")
	buf.WriteString(strconv.Itoa(t.NumMaxGenerations))

	return buf.String()
}

type AnalysisUpdateResponse struct {
	Id      []byte
	Updates []AnalysisUpdate
	// TODO: timestamp
}

func NewAnalysisUpdateResponse(analyzer *biologist.Analyzer, startingGeneration int, maxGenerations int) *AnalysisUpdateResponse {
	// fmt.Printf("NewAnalysisUpdateResponse(%d, %d)\n", startingGeneration, maxGenerations)
	r := new(AnalysisUpdateResponse)

	r.Id = analyzer.Id

	r.Updates = make([]AnalysisUpdate, 0)

	endGen := startingGeneration + maxGenerations
	if analyzer.NumAnalyses() < endGen {
		endGen = analyzer.NumAnalyses()
	} else {
		endGen = analyzer.NumAnalyses() - startingGeneration
	}

	// only add the most recent ones
	for i := startingGeneration; i < endGen; i++ {
		// fmt.Printf(">> Generation %d living <<\n", i)
		update := *NewAnalysisUpdate(analyzer, i)
		// for j, change := range update.Changes {
		// fmt.Printf("  Change[%d] = %s\n", j, change.String())
		// }
		r.Updates = append(r.Updates, update)
	}

	return r
}

func GetAnalysisStatus(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var req AnalysisUpdateRequest

	// fmt.Println(string(body))
	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		fmt.Printf("Received poll request: %s\n", req.String())

		resp := NewAnalysisUpdateResponse(mgr.Analyzer(req.Id), req.StartingGeneration, req.NumMaxGenerations)
		postJson(w, http.StatusCreated, resp)
	}
}

/////////////////////////////////// CONTROL ANALYSIS ///////////////////////////////////

type ControlOrder int

const (
	Start ControlOrder = 0
	Stop  ControlOrder = 1
)

type ControlRequest struct {
	Id    []byte
	Order ControlOrder
}

func (t *ControlRequest) String() string {
	var buf bytes.Buffer
	buf.WriteString("Id: ")
	buf.WriteString(fmt.Sprintf("%x", t.Id))
	buf.WriteString("\nOrder: ")
	switch t.Order {
	case 0:
		buf.WriteString("Start")
	case 1:
		buf.WriteString("Stop")
	}

	return buf.String()
}

func ControlAnalysis(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var req ControlRequest

	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		fmt.Printf("Received control request: %s\n", req.String())

		analyzer := mgr.Analyzer(req.Id)

		switch req.Order {
		case Start:
			analyzer.Start()
		case Stop:
			analyzer.Stop()
		}
	}
}

/////////////////////////////////// MANAGER ///////////////////////////////////

type Manager struct {
	analyzers map[string]*biologist.Analyzer
}

func (t *Manager) stringId(id []byte) string {
	return fmt.Sprintf("%x", id)
}

func (t *Manager) Analyzer(id []byte) *biologist.Analyzer {
	// TODO: validate the input
	return t.analyzers[t.stringId(id)]
}

func (t *Manager) Add(analyzer *biologist.Analyzer) {
	// TODO: validate the input
	t.analyzers[t.stringId(analyzer.Id)] = analyzer
}

func (t *Manager) Remove(id []byte) {
	// TODO: validate the input
	delete(t.analyzers, t.stringId(id))
}

func NewManager() *Manager {
	m := new(Manager)

	m.analyzers = make(map[string]*biologist.Analyzer, 0)

	return m
}

/////////////////////////////////// OTHER ///////////////////////////////////

func postJson(w http.ResponseWriter, httpStatus int, send interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(send); err != nil {
		panic(err)
	}
}

func main() {
	portPtr := flag.Int("port", 8081, "Specify the port to use")
	flag.Parse()

	mux := http.NewServeMux()

	mgr := NewManager()

	mux.HandleFunc("/analyze",
		func(w http.ResponseWriter, r *http.Request) {
			CreateAnalysis(mgr, w, r)
		})
	mux.HandleFunc("/poll",
		func(w http.ResponseWriter, r *http.Request) {
			GetAnalysisStatus(mgr, w, r)
		})
	mux.HandleFunc("/control",
		func(w http.ResponseWriter, r *http.Request) {
			ControlAnalysis(mgr, w, r)
		})

	http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), mux)
}

// vim: set foldmethod=marker:
