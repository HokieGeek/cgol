package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hokiegeek/life/core"
	"io"
	"io/ioutil"
	"net/http"
)

/////////////////////////////////// CREATE ANALYSIS ///////////////////////////////////

type CreateAnalysisResponse struct {
	Id   []byte
	Dims life.Dimensions
	// Rule string
	// Neighbors  life.NeighborsSelector
}

func NewCreateAnalysisResponse(analyzer *life.Analyzer) *CreateAnalysisResponse {
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
	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		// FIXME: this should be sent to a logger
		fmt.Printf("Received create request: %s\n", req.String())

		// Determine the pattern to use for seeding the board
		var patternFunc func(life.Dimensions, life.Location) []life.Location
		switch req.Pattern {
		case USER:
			patternFunc = func(dims life.Dimensions, offset life.Location) []life.Location {
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
		analyzer, err := life.NewAnalyzer(req.Dims, patternFunc, life.ConwayTester())
		if err != nil {
			panic(err)
		}
		mgr.Add(analyzer)

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
	Status     life.Status
	Generation int
	Living     []life.Location
	Changes    []life.ChangedLocation
	// Neighbors life.NeighborSelector
}

func NewAnalysisUpdate(analyzer *life.Analyzer, generation int) *AnalysisUpdate {
	fmt.Printf(" NewAnalysisUpdate(%d)\n", generation)
	a := new(AnalysisUpdate)

	a.Id = analyzer.Id
	a.Dims = analyzer.Life.Dimensions()
	a.Status = analyzer.Life.Status
	a.Generation = generation

	analyzer.Analysis(generation)
	analysis := analyzer.Analysis(generation)

	a.Living = make([]life.Location, len(analysis.Living))
	copy(a.Living, analysis.Living)

	a.Changes = make([]life.ChangedLocation, len(analysis.Changes))
	copy(a.Changes, analysis.Changes)

	return a
}

type AnalysisUpdateRequest struct {
	Id                 []byte
	StartingGeneration int
	NumMaxGenerations  int
}

type AnalysisUpdateResponse struct {
	Id      []byte
	Updates []AnalysisUpdate
	// TODO: timestamp
}

func NewAnalysisUpdateResponse(analyzer *life.Analyzer, startingGeneration int, maxGenerations int) *AnalysisUpdateResponse {
	// fmt.Printf("NewAnalysisUpdateResponse(%d, %d)\n", startingGeneration, maxGenerations)
	r := new(AnalysisUpdateResponse)

	r.Id = analyzer.Id

	r.Updates = make([]AnalysisUpdate, 0)

	endGen := startingGeneration + maxGenerations
	if analyzer.NumAnalyses() < endGen {
		endGen = analyzer.NumAnalyses()
	} else {
		endGen = maxGenerations
	}

	// only add the most recent ones
	for i := startingGeneration; i < endGen; i++ {
		r.Updates = append(r.Updates, *NewAnalysisUpdate(analyzer, i))
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

	fmt.Println(string(body))
	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		fmt.Printf("Received poll request: %x\n", req.Id)

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

func ControlAnalysis(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var req ControlRequest

	fmt.Println(string(body))
	if err := json.Unmarshal(body, &req); err != nil {
		postJson(w, 422, err)
	} else {
		fmt.Printf("Received control request: %x\n", req.Id)

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
	analyzers map[string]*life.Analyzer
}

func (t *Manager) stringId(id []byte) string {
	return fmt.Sprintf("%x", id)
}

func (t *Manager) Analyzer(id []byte) *life.Analyzer {
	// TODO: validate the input
	return t.analyzers[t.stringId(id)]
}

func (t *Manager) Add(analyzer *life.Analyzer) {
	// TODO: validate the input
	t.analyzers[t.stringId(analyzer.Id)] = analyzer
}

func (t *Manager) Remove(id []byte) {
	// TODO: validate the input
	delete(t.analyzers, t.stringId(id))
}

func NewManager() *Manager {
	m := new(Manager)

	m.analyzers = make(map[string]*life.Analyzer, 0)

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
