package life

import (
	"testing"
	"time"
)

func TestLifeCreation(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New(
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	expected, _ := newPond(Dimensions{Height: 3, Width: 3}, newTracker(), NEIGHBORS_ALL)
	expected.SetOrganisms([]Location{Location{X: 0, Y: 1}, Location{X: 1, Y: 1}, Location{X: 2, Y: 1}})

	if !strategy.pond.Equals(expected) {
		t.Fatalf("Actual board\n%s\ndoes not match expected\n%s\n", strategy.pond.String(), expected.String())
	}
}

func TestLifeProcess(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New(
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	seededpond, err := strategy.pond.Clone()
	if err != nil {
		t.Fatalf("Unable to clone pond: %s\n", err)
	}

	strategy.process()

	if seededpond.Equals(strategy.pond) {
		t.Fatal("pond did not change after processing")
	}

	// Check statistics
	if strategy.Generations != 1 {
		t.Errorf("Tracked %d generations when there should only be one\n", strategy.Generations)
	}
}

func TestLifeStart(t *testing.T) {
	t.Skip("whoops")
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New(
		dims,
		NEIGHBORS_ALL,
		func(dimensions Dimensions, offset Location) []Location {
			return Random(dimensions, offset, 85)
		},
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	seededpond, err := strategy.pond.Clone()
	if err != nil {
		t.Fatalf("Unable to clone pond: %s\n", err)
	}

	stop := strategy.Start(nil)

	// t.Log(strategy.String())

	time.Sleep(time.Millisecond * 10)
	stop()
	// t.Log(strategy.String())

	processedpond := strategy.pond

	if seededpond.Equals(processedpond) {
		t.Fatal("pond did not change after processing")
	}

	// Check statistics
	if strategy.Generations < 2 {
		t.Errorf("Tracked %d generations when there should be two or more\n", strategy.Generations)
	}
}

func TestLifeGeneration(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New(
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	expectedNumLiving := 3
	expectedGen := 31
	gen := strategy.Generation(expectedGen)

	if gen.Num != expectedGen {
		t.Error("Retrieved %d generations instead of %d\n", gen.Num, expectedGen)
	}

	if len(gen.Living) != expectedNumLiving {
		t.Fatalf("Retrieved %d living organisms instead of %d\n", len(gen.Living), expectedNumLiving)
	}
}

func TestLifeString(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New(
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	if len(strategy.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}

func TestLifeDimensions(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	life, err := New(
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		ConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	retrievedDims := life.Dimensions()

	if !retrievedDims.Equals(&dims) {
		t.Error("Retrieved dimensions did not match expected dimensions")
	}
}

// vim: set foldmethod=marker:
