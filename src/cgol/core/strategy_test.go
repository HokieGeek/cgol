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
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyCreation",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	expected, _ := newLifeboard(LifeboardDims{3, 3})
	expected.SetValue(LifeboardLocation{X: 0, Y: 1}, 0)
	expected.SetValue(LifeboardLocation{X: 1, Y: 1}, 0)
	expected.SetValue(LifeboardLocation{X: 2, Y: 1}, 0)

	if !strategy.pond.lifeboard.Equals(expected) {
		t.Fatalf("Actual lifeboard\n%s\ndoes not match expected\n%s\n", strategy.pond.lifeboard.String(), expected.String())
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
	t.Skip("This doesn't work as expected")
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

	updates := make(chan bool)
	strategy.Start(updates)

	// go func() {
	time.Sleep(strategy.UpdateRate * 1)
	strategy.Stop()
	/*}()

	for {
		select {
		case <-updates:
			t.Log(strategy.String())
		}
	}
	*/

	processedPond := strategy.pond

	if seededPond.Equals(processedPond) {
		t.Fatal("Pond did not change after processing")
	}

	// Check statistics
	if strategy.Statistics.Generations < 2 {
		t.Errorf("Tracked %d generations when there should be two or more\n", strategy.Statistics.Generations)
	}
}

func TestStrategyGetGeneration(t *testing.T) {
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	strategy := NewStrategy("TestStrategyString",
		pond,
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	expectedNumLiving := 3
	expectedGen := 31
	gen := strategy.GetGeneration(expectedGen)

	if gen.Num != expectedGen {
		t.Error("Retrieved %d generations instead of %d\n", gen.Num, expectedGen)
	}

	if len(gen.Living) != expectedNumLiving {
		t.Fatalf("Retrieved %d living organisms instead of %d\n", len(gen.Living), expectedNumLiving)
	}
	// TODO
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
