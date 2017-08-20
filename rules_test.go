package life

import "testing"

func TestRulesString(t *testing.T) {
	rules := &Rules{[]int{1, 2, 3}, []int{4, 5}}

	if len(rules.String()) <= 0 {
		t.Error("Rules unexpectedly returned empty string")
	}
}

func TestRuleRulesConway(t *testing.T) {
	const (
		stdUnderpopulation = 2
		stdOvercrowding    = 3
		stdRevive          = 3
	)

	RulesConwayLife := ConwayTester()

	// Rule #1: Kill a cell that has few neighbors
	if RulesConwayLife(stdUnderpopulation-1, true) {
		t.Fatal("Cell that should be dead due to underpopulation is still alive")
	}

	// Rule #2: Keep cell alive
	if !RulesConwayLife(stdUnderpopulation, true) {
		t.Fatalf("Killed cell that should still be alive with %d neighbors\n", stdUnderpopulation)
	}
	if !RulesConwayLife(stdOvercrowding, true) {
		t.Fatalf("Killed cell that should still be alive with %d neighbors\n", stdOvercrowding)
	}

	// Rule #3: Kill a cell that has too many neighbors
	if RulesConwayLife(stdOvercrowding+1, true) {
		t.Fatal("Cell that should be dead due to overcrowding is still alive")
	}

	// Rule #4: Dead cell being revived
	if !RulesConwayLife(stdRevive, false) {
		t.Fatal("Did not revive dead cell")
	}

	if RulesConwayLife(stdRevive+1, false) {
		t.Fatal("Revived dead cell that shouldn't have been")
	}
}

// vim: set foldmethod=marker:
