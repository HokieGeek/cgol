package main

import (
	"encoding/json"
	"fmt"
	"github.com/hokiegeek/life/core"
	"net/http"
)

type AnalysisResp struct {
	Name string
	Age  int
	Locs []life.Location
}

func CreateAnalysis(w http.ResponseWriter, r *http.Request) {
	list := make([]life.Location, 0)
	list = append(list, life.Location{0, 1})
	list = append(list, life.Location{2, 3})
	list = append(list, life.Location{4, 5})
	list = append(list, life.Location{6, 7})
	list = append(list, life.Location{8, 9})

	test := &AnalysisResp{Name: "Hugh F. Kares", Age: 42, Locs: list}

	fmt.Println(test)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(test); err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create", CreateAnalysis)
	// mux.HandleFunc("/poll", ReturnAnalysisStatus)

	http.ListenAndServe(":8080", mux)
}
