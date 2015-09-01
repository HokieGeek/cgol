package cgol

import "testing"

func TestStrategyStatsString(t *testing.T) {
	stats := new(StrategyStats)
	stats.Generations++

	if len(stats.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}

func TestStrategyCreation(t *testing.T) {
	t.Skip("TODO")
	// strategy := NewStrategy("TestStrategyCreation",
	// 	NewPond(3, 3, NEIGHBORS_ALL),
	// 	Blinkers,
	// 	RulesConwayLife,
	// 	SimultaneousProcessor)
	// TODO: Check that the pond was initialized correctly
	// TODO: Check statistics
}

func TestStrategyProcess(t *testing.T) {
	t.Skip("TODO: Check that process actually changes the pond")
	// strategy := NewStrategy("TestStrategyProcess",
	// 	NewPond(3, 3, NEIGHBORS_ALL),
	// 	Blinkers,
	// 	RulesConwayLife,
	// 	SimultaneousProcessor)
	// strategy.process()
}

func TestStrategyStartStop(t *testing.T) {
	t.Skip("TODO: Check that after a couple of generations it has changed. check statistics too")
	// strategy := NewStrategy("TestStrategyStartStop",
	// 	NewPond(3, 3, NEIGHBORS_ALL),
	// 	Blinkers,
	// 	RulesConwayLife,
	// 	SimultaneousProcessor)
	// strategy.Start(make(chan bool))
	/*
		go func() {
			time.Sleep(strategy.UpdateRate * 2)
			strategy.Stop()
		}()
	*/
}

func TestStrategyString(t *testing.T) {
	strategy := NewStrategy("TestStrategyString",
		NewPond(3, 3, NEIGHBORS_ALL),
		Blinkers,
		RulesConwayLife,
		SimultaneousProcessor)

	if len(strategy.String()) <= 0 {
		t.Error("String function unexpectly returned an empty string")
	}
}
