package life

import (
	"testing"
	"time"
)

func TestStatisticsString(t *testing.T) {
	stats := new(Statistics)
	stats.Generations++

	if len(stats.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}

func TestStatusString(t *testing.T) {
	var status Status

	status = Seeded
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Status object")
	}

	status = Active
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Status object")
	}

	status = Stable
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Status object")
	}

	status = Dead
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Status object")
	}
}

func TestLifeCreation(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeCreation",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		GetConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	expected, _ := newBoard(Dimensions{3, 3})
	expected.SetValue(Location{X: 0, Y: 1}, 0)
	expected.SetValue(Location{X: 1, Y: 1}, 0)
	expected.SetValue(Location{X: 2, Y: 1}, 0)

	if !strategy.pond.board.Equals(expected) {
		t.Fatalf("Actual board\n%s\ndoes not match expected\n%s\n", strategy.pond.board.String(), expected.String())
	}
}

func TestLifeProcess(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeProcess",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		GetConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	seededpond, err := strategy.pond.Clone()
	if err != nil {
		t.Fatalf("Unable to clone pond: %s\n", err)
	}

	strategy.process()

	processedpond := strategy.pond

	if seededpond.Equals(processedpond) {
		t.Fatal("pond did not change after processing")
	}

	// Check statistics
	if strategy.Stats.Generations != 1 {
		t.Errorf("Tracked %d generations when there should only be one\n", strategy.Stats.Generations)
	}
}

func TestLifeStartRated(t *testing.T) {
	t.Skip("Unpredictable. Needs refactoring")
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeStartStop",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		GetConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	seededpond, err := strategy.pond.Clone()
	if err != nil {
		t.Fatalf("Unable to clone pond: %s\n", err)
	}

	testRate := time.Second

	stop := strategy.Start(nil, testRate)

	// t.Log(strategy.String())

	time.Sleep(testRate * 4)
	stop()
	// t.Log(strategy.String())

	processedpond := strategy.pond

	if seededpond.Equals(processedpond) {
		t.Fatal("pond did not change after processing")
	}

	// Check statistics
	if strategy.Stats.Generations < 2 {
		t.Errorf("Tracked %d generations when there should be two or more\n", strategy.Stats.Generations)
	}
}

func TestLifeStartFullTilt(t *testing.T) {
	t.Skip("Unpredictable. Needs refactoring")
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeStartStop",
		dims,
		NEIGHBORS_ALL,
		func(dimensions Dimensions, offset Location) []Location {
			return Random(dimensions, offset, 85)
		},
		GetConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	seededpond, err := strategy.pond.Clone()
	if err != nil {
		t.Fatalf("Unable to clone pond: %s\n", err)
	}

	stop := strategy.Start(nil, 0)

	// t.Log(strategy.String())

	time.Sleep(time.Second * 2)
	stop()
	// t.Log(strategy.String())

	processedpond := strategy.pond

	if seededpond.Equals(processedpond) {
		t.Fatal("pond did not change after processing")
	}

	// Check statistics
	if strategy.Stats.Generations < 2 {
		t.Errorf("Tracked %d generations when there should be two or more\n", strategy.Stats.Generations)
	}
}

func TestLifeGeneration(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeString",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		GetConwayTester(),
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
	// TODO
}

func TestLifeString(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	strategy, err := New("TestLifeString",
		dims,
		NEIGHBORS_ALL,
		Blinkers,
		GetConwayTester(),
		SimultaneousProcessor)
	if err != nil {
		t.Fatalf("Unable to create strategy: %s\n", err)
	}

	if len(strategy.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}
