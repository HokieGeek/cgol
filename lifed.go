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
		anaylzer, err := life.NewAnalyzer(req.Dims)
		if err != nil {
			panic(err)
		}
		// TODO: Add to the manager
		mgr.Analyzers = append(mgr.Analyzers, anaylzer)

		fmt.Printf("Id: %x\n", anaylzer.Id)

		// Respond the request with the ID of the analyzer
		resp := NewCreateAnalysisResponse(anaylzer)

		postJson(w, http.StatusCreated, resp)
	}
}

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
	Id []byte
}

type AnalysisUpdateResponse struct {
	Id      []byte
	Updates []AnalysisUpdate
	// TODO: timestamp
}

func NewAnalysisUpdateResponse(analyzer *life.Analyzer) *AnalysisUpdateResponse {
	r := new(AnalysisUpdateResponse)

	r.Id = analyzer.Id

	r.Updates = make([]AnalysisUpdate, 0)

	// TODO: only add the most recent ones. The manager should keep a pointer
	// r.Updates = append(r.Updates, *NewAnalysisUpdate(analyzer, 0))
	for i := 0; i < analyzer.NumAnalyses(); i++ {
		r.Updates = append(r.Updates, *NewAnalysisUpdate(analyzer, i))
		fmt.Printf("Num changes: %d\n", len(r.Updates[i].Changes))
	}
	fmt.Printf("Num updates: %d\n", len(r.Updates))

	return r
}

func GetAnalysisStatus(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	// func GetAnalysisStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetAnalysisStatus()\n")
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

		// TODO: Retrieve from the manager

		// fmt.Printf("ID: %x\n", analyzer.Id)

		// Respond the request with the ID of the analyzer
		resp := NewAnalysisUpdateResponse(mgr.Analyzers[0])
		postJson(w, http.StatusCreated, resp)
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
	Analyzers []*life.Analyzer
}

func main() {
	mux := http.NewServeMux()

	// TODO: create the manager here and the handlers below will take an anon func

	mgr := new(Manager)
	mgr.Analyzers = make([]*life.Analyzer, 0)

	mux.HandleFunc("/analyze",
		func(w http.ResponseWriter, r *http.Request) {
			CreateAnalysis(mgr, w, r)
		})
	mux.HandleFunc("/poll",
		func(w http.ResponseWriter, r *http.Request) {
			GetAnalysisStatus(mgr, w, r)
		})

	http.ListenAndServe(":8081", mux)
}
