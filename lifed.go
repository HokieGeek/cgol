package main

import (
	"bytes"
	"encoding/json"
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

type CreateAnalysisRequest struct {
	Dims life.Dimensions
	// life.Rules
	// Seed
	// Processor
}

func (t *CreateAnalysisRequest) String() string {
	var buf bytes.Buffer

	buf.WriteString(t.Dims.String())

	return buf.String()
}

func CreateAnalysis(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	// Retrieve the necessary stuffs
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

		// Create the analyzer
		analyzer, err := life.NewAnalyzer(req.Dims)
		if err != nil {
			panic(err)
		}
		mgr.Add(analyzer)

		fmt.Printf("Id: %x\n", analyzer.Id)

		// Respond the request with the ID of the analyzer
		resp := NewCreateAnalysisResponse(analyzer)

		postJson(w, http.StatusCreated, resp)
	}
}

/////////////////////////////////// UPDATE ANALYSIS ///////////////////////////////////

type AnalysisUpdate struct {
	Id         []byte
	Status     life.Status
	Generation int
	Living     []life.Location
	Changes    []life.ChangedLocation
	// Neighbors life.NeighborSelector
}

func NewAnalysisUpdate(analyzer *life.Analyzer, generation int) *AnalysisUpdate {
	a := new(AnalysisUpdate)

	a.Id = analyzer.Id
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
}

type AnalysisUpdateResponse struct {
	Id      []byte
	Updates []AnalysisUpdate
	// TODO: timestamp
}

func NewAnalysisUpdateResponse(analyzer *life.Analyzer, startingGeneration int) *AnalysisUpdateResponse {
	r := new(AnalysisUpdateResponse)

	r.Id = analyzer.Id

	r.Updates = make([]AnalysisUpdate, 0)

	// only add the most recent ones
	for i := startingGeneration; i < analyzer.NumAnalyses(); i++ {
		r.Updates = append(r.Updates, *NewAnalysisUpdate(analyzer, i))
	}

	return r
}

func GetAnalysisStatus(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	// Retrieve the necessary stuffs
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

		resp := NewAnalysisUpdateResponse(mgr.Analyzer(req.Id), req.StartingGeneration)
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
	fmt.Printf("ControlAnalysis()\n")

	// Retrieve the necessary stuffs
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

		// TODO: Retrieve from the manager

		// fmt.Printf("ID: %x\n", analyzer.Id)

		// Respond the request with the ID of the analyzer
		// resp := NewAnalysisUpdateResponse(mgr.Analyzers[0], req.StartingGeneration)
		// postJson(w, http.StatusCreated, resp)
	}
}

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

type Manager struct {
	analyzers map[string]*life.Analyzer
}

func (t *Manager) stringId(id []byte) string {
	// n := bytes.IndexByte(id, 0)
	// return string(id[:n])
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

func main() {
	mux := http.NewServeMux()

	// TODO: create the manager here and the handlers below will take an anon func

	mgr := NewManager()
	// mgr := new(Manager)
	// mgr.Analyzers = make([]*life.Analyzer, 0)
	// mgr.Analyzers = make(map[[]byte]*life.Analyzer, 0)

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

	http.ListenAndServe(":8081", mux)
}
