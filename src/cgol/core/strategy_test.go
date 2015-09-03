package cgol

import (
	"testing"
	"time"
)

func TestStrategyStatsString(t *testing.T) {
	stats := new(StrategyStats)
	stats.Generations++

	if len(stats.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}

func TestStrategyCreation(t *testing.T) {
	t.Skip("whoops")
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyCreation",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	expected, _ := NewGameboard(GameboardDims{3, 3})
	expected.SetValue(GameboardLocation{X: 1, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 2}, 0)

	if !strategy.pond.gameboard.Equals(expected) {
		t.Fatalf("Actual gameboard\n%s\ndoes not match expected\n%s\n", strategy.pond.gameboard.String(), expected.String())
	}
}

func TestStrategyProcess(t *testing.T) {
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyProcess",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	seededPond := strategy.pond.Clone()

	strategy.process()

	processedPond := strategy.pond

	if seededPond.Equals(processedPond) {
		t.Fatal("Pond did not change after processing")
	}

	// Check statistics
	if strategy.Statistics.Generations != 1 {
		t.Errorf("Tracked %d generations when there should only be one\n", strategy.Statistics.Generations)
	}
}

func TestStrategyStartStop(t *testing.T) {
	t.Skip("not working")
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyStartStop",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	seededPond := strategy.pond.Clone()

	strategy.Start(make(chan bool))

	time.Sleep(strategy.UpdateRate * 4)
	strategy.Stop()

	processedPond := strategy.pond

	if seededPond.Equals(processedPond) {
		t.Fatal("Pond did not change after processing")
	}

	// Check statistics
	if strategy.Statistics.Generations < 2 {
		t.Errorf("Tracked %d generations when there should only be two or more\n", strategy.Statistics.Generations)
	}
}

func TestStrategyGetGeneration(t *testing.T) {
	t.Skip("TODO")
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyString",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	expectedGen := 31
	gen := strategy.GetGeneration(expectedGen)

	if gen.num != expectedGen {
		t.Error("WTF")
	}

	// TODO: Check that the living count makes sense
}

func TestStrategyString(t *testing.T) {
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyString",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	if len(strategy.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}
