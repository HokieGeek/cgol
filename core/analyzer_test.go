package life

import (
	"testing"
	"time"
)

func TestAnalyzerCreate(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	if len(analyzer.Life.Seed) <= 0 {
		t.Error("Created analyzer with an empty seed")
	}
}

func TestAnalyzerCreateError(t *testing.T) {
	size := Dimensions{Width: 0, Height: 0}
	_, err := NewAnalyzer(size, Blinkers)
	if err == nil {
		t.Fatal("Unexpectedly successful at creating analyzer with board of 0 size")
	}
}

func TestAnalyzerString(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	if len(analyzer.String()) <= 0 {
		t.Error("Analyzer String function returned empty string")
	}
}

func TestAnalyzerStart(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	analyzer.Start()
	waitTime := time.Millisecond * 50
	time.Sleep(waitTime)
	analyzer.Stop()

	if analyzer.NumAnalyses() <= 0 {
		t.Fatalf("No analyses found after %s of running\n", waitTime.String())
	}
}

func TestAnalyzerStop(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	analyzer.Start()
	time.Sleep(time.Millisecond * 3)
	analyzer.Stop()

	time.Sleep(time.Millisecond * 1)
	stoppedCount := analyzer.NumAnalyses()

	time.Sleep(time.Millisecond * 10)

	waitedCount := analyzer.NumAnalyses()
	if stoppedCount != waitedCount {
		t.Fatalf("Analyses increased after stopped. Expected %d and got %d\n", stoppedCount, waitedCount)
	}
}

func TestAnalyzerAnalysis(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	analyzer.Start()
	time.Sleep(time.Millisecond * 10)
	analyzer.Stop()

	if analyzer.Analysis(0) == nil {
		t.Fatal("Could not retrieve seed")
	}

	for i := analyzer.NumAnalyses() - 1; i >= 0; i-- {
		if analyzer.Analysis(i) == nil {
			t.Fatalf("Analysis for generation %d is nil\n", i)
		}
	}
}

func TestAnalyzerAnalysisError(t *testing.T) {
	size := Dimensions{Width: 3, Height: 3}
	analyzer, err := NewAnalyzer(size, Blinkers)
	if err != nil {
		t.Fatalf("Unable to create analyzer: %s\n", err)
	}

	analyzer.Start()
	time.Sleep(time.Millisecond * 10)
	analyzer.Stop()

	if analyzer.Analysis(-1) != nil {
		t.Fatal("Analyzer returned to me analysis at generation -1")
	}

	if analyzer.Analysis(analyzer.NumAnalyses()) != nil {
		t.Fatal("Analyzer returned to me analysis at generation greater than the number of generations analyzed")
	}
}
