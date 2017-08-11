package main

import (
	"gitlab.com/hokiegeek/biologist"
	"gitlab.com/hokiegeek/life"
	"testing"
)

func TestNewCreateAnalysisResponse(t *testing.T) {
	size := life.Dimensions{Width: 3, Height: 3}
	analyzer, err := biologist.NewAnalyzer(size, life.Blinkers, life.ConwayTester())
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	resp := NewCreateAnalysisResponse(analyzer)

	if !resp.Dims.Equals(&size) {
		t.Fatal("Expected size %s but received %s\n", size.String(), resp.Dims.String())
	}
}

/*
func TestNewAnalysisUpdateResponse(t *testing.T) {
	size := life.Dimensions{Width: 3, Height: 3}
	analyzer, err := life.NewAnalyzer(size, life.Blinkers, life.ConwayTester())
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	resp := NewAnalysisUpdateResponse(analyzer, 0, 1)
}
*/

// vim: set foldmethod=marker:
