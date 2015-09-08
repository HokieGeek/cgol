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

var currentAnalyzer *life.Analyzer // FIXME: holy crap!!

func CreateAnalysis(w http.ResponseWriter, r *http.Request) {
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
		currentAnalyzer, err := life.NewAnalyzer(req.Dims)
		if err != nil {
			panic(err)
		}
		// TODO: Add to the manager

		fmt.Printf("Id: %x\n", currentAnalyzer.Id)

		// Respond the request with the ID of the analyzer
		resp := NewCreateAnalysisResponse(currentAnalyzer)

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

	analysis := analyzer.Analysis(generation)
	copy(a.Living, analysis.Living)
	copy(a.Changes, analysis.Changes)

	return a
}

type AnalysisUpdateRequest struct {
	Id []byte
}

type AnalysisUpdateResponse struct {
	Id      []byte
	Updates []AnalysisUpdate
}

func NewAnalysisUpdateResponse(analyzer *life.Analyzer) *AnalysisUpdateResponse {
	r := new(AnalysisUpdateResponse)

	r.Id = analyzer.Id

	r.Updates = make([]AnalysisUpdate, 0)
	r.Updates = append(r.Updates, *NewAnalysisUpdate(analyzer, 0))

	return r
}

func GetAnalysisStatus(w http.ResponseWriter, r *http.Request) {
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

		// TODO: Add to the manager

		// fmt.Printf("ID: %x\n", analyzer.Id)

		// Respond the request with the ID of the analyzer
		/*
			living := []life.Location{life.Location{X: 10, Y: 10},
				life.Location{X: 11, Y: 10},
				life.Location{X: 12, Y: 10}}
			update := &AnalysisUpdate{Id: req.Id, Status: life.Seeded, Generation: 1, Living: living}
			updates := []AnalysisUpdate{*update}
			resp := &AnalysisUpdateResponse{Id: req.Id, Updates: updates}
		*/
		resp := NewAnalysisUpdateResponse(currentAnalyzer)

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

func main() {
	mux := http.NewServeMux()

	// TODO: create the manager here and the handlers below will take an anon func

	mux.HandleFunc("/create", CreateAnalysis)
	mux.HandleFunc("/poll", GetAnalysisStatus)

	http.ListenAndServe(":8081", mux)
}
