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
}

func TestStrategyProcess(t *testing.T) {
	t.Skip("TODO")
}

func TestStrategyStartStop(t *testing.T) {
	t.Skip("TODO")
}

func TestStrategyString(t *testing.T) {
	t.Skip("TODO")
	/*
		strategy := NewStrategy()

		if len(strategy.String()) <= 0 {
			t.Error("String function unexpectly returned an empty string")
		}
	*/
}
