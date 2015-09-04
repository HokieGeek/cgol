package cgol

import "testing"

func TestRuleRulesConway(t *testing.T) {
	const (
		STD_UNDERPOPULATION = 2
		STD_OVERCROWDING    = 3
		STD_REVIVE          = 3
	)

	RulesConwayLife := GetConwayTester()

	// Rule #1: Kill a cell that has few neighbors
	if RulesConwayLife(STD_UNDERPOPULATION-1, true) {
		t.Fatal("Cell that should be dead due to underpopulation is still alive")
	}

	// Rule #2: Keep cell alive
	if !RulesConwayLife(STD_UNDERPOPULATION, true) {
		t.Fatalf("Killed cell that should still be alive with %d neighbors\n", STD_UNDERPOPULATION)
	}
	if !RulesConwayLife(STD_OVERCROWDING, true) {
		t.Fatalf("Killed cell that should still be alive with %d neighbors\n", STD_OVERCROWDING)
	}

	// Rule #3: Kill a cell that has too many neighbors
	if RulesConwayLife(STD_OVERCROWDING+1, true) {
		t.Fatal("Cell that should be dead due to overcrowding is still alive")
	}

	// Rule #4: Dead cell being revived
	if !RulesConwayLife(STD_REVIVE, false) {
		t.Fatal("Did not revive dead cell")
	}

	if RulesConwayLife(STD_REVIVE+1, false) {
		t.Fatal("Revived dead cell that shouldn't have been")
	}
}
