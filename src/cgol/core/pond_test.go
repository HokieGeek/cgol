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

func TestLivingTrackerSetTest(t *testing.T) {
	tracker := NewLivingTracker()

	// Ok, set the value
	loc := GameboardLocation{X: 42, Y: 24}
	tracker.Set(loc)

	// Test that it exists
	if !tracker.Test(loc) {
		t.Fatal("Added location unexpectedly tested false")
	}
}

func TestLivingTrackerTestError(t *testing.T) {
	tracker := NewLivingTracker()

	if tracker.Test(GameboardLocation{X: 0, Y: 0}) {
		t.Error("Unexpectedly tested true a location that does not exist in structure")
	}
}

func TestLivingTrackerRemove(t *testing.T) {
	tracker := NewLivingTracker()

	// Ok, set the value
	loc := GameboardLocation{X: 42, Y: 24}
	tracker.Set(loc)

	tracker.Remove(loc)

	// Test that it exists
	if tracker.Test(loc) {
		t.Fatal("Unexpectedly tested true a location that was removed")
	}
}

func TestLivingTrackerRemoveError(t *testing.T) {
	t.Skip("TODO")
}

func TestLivingTrackerGetAll(t *testing.T) {
	tracker := NewLivingTracker()

	// Create the expected values
	expectedLocations := make([]GameboardLocation, 3)
	expectedLocations[0] = GameboardLocation{X: 42, Y: 42}
	expectedLocations[1] = GameboardLocation{X: 11, Y: 11}
	expectedLocations[2] = GameboardLocation{X: 12, Y: 34}

	// Ok, set the values
	for _, l := range expectedLocations {
		tracker.Set(l)
	}

	actualLocations := tracker.GetAll()

	if len(actualLocations) != len(expectedLocations) {
		t.Fatalf("Received %d locations when I expected %d\n", len(actualLocations), len(expectedLocations))
	}

	// Check that all returned locations were expected
	for _, actual := range actualLocations {
		found := false
		for _, expected := range expectedLocations {
			if actual.Equals(&expected) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Returned location %s was not in list of test locations\n", actual.String())
			break
		}
	}

	// Check that all expected locations were returned
	for _, expected := range expectedLocations {
		found := false
		for _, actual := range actualLocations {
			if expected.Equals(&actual) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected location %s was not in list of returned locations\n", expected.String())
			break
		}
	}
}

func TestLivingTrackerCount(t *testing.T) {
	tracker := NewLivingTracker()

	// Create the expected values
	expectedLocations := make([]GameboardLocation, 3)
	expectedLocations[0] = GameboardLocation{X: 42, Y: 42}
	expectedLocations[1] = GameboardLocation{X: 11, Y: 11}
	expectedLocations[2] = GameboardLocation{X: 12, Y: 34}

	// Ok, set the values
	for _, l := range expectedLocations {
		tracker.Set(l)
	}

	// Test the counter
	expectedCount := len(expectedLocations)
	count := tracker.GetCount()
	if count != expectedCount {
		t.Fatalf("Retrieved count of %d instead of expected %d\n", count, expectedCount)
	}

	// Now remove a location and try again
	tracker.Remove(expectedLocations[2])
	expectedCount--
	count = tracker.GetCount()
	if count != expectedCount {
		t.Fatalf("Retrieved count of %d instead of expected %d after remove a location\n", count, expectedCount)
	}
}

func TestPondSettingInitialPatterns(t *testing.T) {
	rows := 3
	cols := 3
	pond, err := NewPond(rows, cols, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 3)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	initialLiving[2] = GameboardLocation{X: 2, Y: 2}

	pond.SetOrganisms(initialLiving)

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
	pond, err := NewPond(1, 1, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	pond.setOrganismValue(pos, expectedVal)

	actualVal := pond.GetOrganismValue(pos)

	if actualVal != expectedVal {
		t.Fatalf("Retrieved actual value %d instead of expected value %d\n", actualVal, expectedVal)
	}
}

func TestPondGetNumLiving(t *testing.T) {
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 2)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	pond.SetOrganisms(initialLiving)

	numLiving := pond.GetNumLiving()

	if numLiving != len(initialLiving) {
		t.Fatalf("Retrieved actual %d living orgnamisms instead of expected %d\n", numLiving, len(initialLiving))
	}
}

func TestPondNeighborCountCalutation(t *testing.T) {
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 2)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	pond.SetOrganisms(initialLiving)

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
	pond, err := NewPond(3, 3, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	if len(pond.String()) <= 0 {
		t.Error("Unexpectly retrieved empty string from Pond string function")
	}
}
