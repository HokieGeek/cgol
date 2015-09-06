package life

import "testing"

func TestTrackerSetTest(t *testing.T) {
	tracker := newTracker()

	// Ok, set the value
	loc := Location{X: 42, Y: 24}
	tracker.Set(loc)

	// Test that it exists
	if !tracker.Test(loc) {
		t.Fatal("Added location unexpectedly tested false")
	}
}

func TestTrackerTestError(t *testing.T) {
	tracker := newTracker()

	if tracker.Test(Location{X: 0, Y: 0}) {
		t.Error("Unexpectedly tested true a location that does not exist in structure")
	}
}

func TestTrackerRemove(t *testing.T) {
	tracker := newTracker()

	// Ok, set the value
	loc := Location{X: 42, Y: 24}
	tracker.Set(loc)

	tracker.Remove(loc)

	// Test that it exists
	if tracker.Test(loc) {
		t.Fatal("Unexpectedly tested true a location that was removed")
	}
}

func TestTrackerRemoveError(t *testing.T) {
	t.Skip("TODO")
}

func TestTrackerGetAll(t *testing.T) {
	tracker := newTracker()

	// Create the expected values
	expectedLocations := make([]Location, 3)
	expectedLocations[0] = Location{X: 42, Y: 42}
	expectedLocations[1] = Location{X: 11, Y: 11}
	expectedLocations[2] = Location{X: 12, Y: 34}

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

func TestTrackerCount(t *testing.T) {
	tracker := newTracker()

	// Create the expected values
	expectedLocations := make([]Location, 3)
	expectedLocations[0] = Location{X: 42, Y: 42}
	expectedLocations[1] = Location{X: 11, Y: 11}
	expectedLocations[2] = Location{X: 12, Y: 34}

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
