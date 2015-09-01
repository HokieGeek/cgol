package cgol

import "testing"

func TestPondStatusString(t *testing.T) {
	var status PondStatus

	status = Active
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from PondStatus object")
	}

	status = Stable
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from PondStatus object")
	}

	status = Dead
	if len(status.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from PondStatus object")
	}
}

func TestNeighborSelectorString(t *testing.T) {
	var selector NeighborsSelector

	selector = NEIGHBORS_ORTHOGONAL
	if len(selector.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Pondselector object")
	}

	selector = NEIGHBORS_OBLIQUE
	if len(selector.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from Pondselector object")
	}

	selector = NEIGHBORS_ALL
	if len(selector.String()) <= 0 {
		t.Error("Unexpectedly retrieved empty string from PondStatus object")
	}
}

func TestPondSettingInitialPatterns(t *testing.T) {
	rows := 3
	cols := 3
	pond := NewPond(rows, cols, NEIGHBORS_ALL)

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 3)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	initialLiving[2] = GameboardLocation{X: 2, Y: 2}

	pond.init(initialLiving)

	// Check each expected value
	for _, loc := range initialLiving {
		if !pond.isOrganismAlive(loc) {
			t.Fatalf("Seed organism is not alive!: %s\n", loc.String())
		}
	}
}

func TestPondNeighborSelection(t *testing.T) {
	t.Skip("This will essentially just retest the gameboard tests.")
}

func TestPondNeighborSelectionError(t *testing.T) {
	t.Skip("Bad location should be the test")
	/*
		fake_selector := 999
		pond := NewPond(1, 1, fake_selector)

		neighbors, err := pond.GetNeighbors(GameboardLocation{X: 0, Y: 0})
		if err != nil {
			t.Error("Did not encounter error when retrieving neighbors using bogus selector")
		}
	*/
}

func TestPondOrganismValue(t *testing.T) {
	expectedVal := 2
	pos := GameboardLocation{X: 0, Y: 0}
	pond := NewPond(1, 1, NEIGHBORS_ALL)
	pond.setOrganismValue(pos, expectedVal)

	actualVal := pond.GetOrganismValue(pos)

	if actualVal != expectedVal {
		t.Fatalf("Retrieved actual value %d instead of expected value %d\n", actualVal, expectedVal)
	}
}

func TestPondGetNumLiving(t *testing.T) {
	pond := NewPond(3, 3, NEIGHBORS_ALL)

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 2)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	pond.init(initialLiving)

	numLiving := pond.GetNumLiving()

	if numLiving != len(initialLiving) {
		t.Fatalf("Retrieved actual %d living orgnamisms instead of expected %d\n", numLiving, len(initialLiving))
	}
}

func TestPondNeighborCountCalutation(t *testing.T) {
	pond := NewPond(3, 3, NEIGHBORS_ALL)

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 2)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	pond.init(initialLiving)

	expectedNeighbors := make([]GameboardLocation, 5)
	expectedNeighbors[0] = GameboardLocation{X: 0, Y: 0}
	expectedNeighbors[1] = GameboardLocation{X: 1, Y: 0}
	expectedNeighbors[2] = GameboardLocation{X: 1, Y: 1}
	expectedNeighbors[3] = GameboardLocation{X: 0, Y: 2}
	expectedNeighbors[4] = GameboardLocation{X: 1, Y: 2}

	expectedNeighborCount := len(initialLiving)
	actualNeighborCount, actualNeighbors := pond.calculateNeighborCount(GameboardLocation{X: 0, Y: 1})

	if actualNeighborCount != expectedNeighborCount {
		t.Fatalf("Retrieved %d neighbor count instead of expected %d\n", actualNeighborCount, expectedNeighborCount)
	}

	if len(actualNeighbors) != len(expectedNeighbors) {
		t.Fatalf("Retrieved %d neighbors instead of expected %d\n", len(actualNeighbors), expectedNeighborCount)
	}

	for _, neighbor := range actualNeighbors {
		found := false
		for _, expected := range expectedNeighbors {
			if expected.Equals(&neighbor) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Found unexpected neighbor %s\n", neighbor.String())
		}
	}
}

func TestPondString(t *testing.T) {
	pond := NewPond(3, 3, NEIGHBORS_ALL)
	if len(pond.String()) <= 0 {
		t.Error("Unexpectly retrieved empty string from Pond string function")
	}
}
