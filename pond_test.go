package life

import "testing"

func TestNeighborSelectorString(t *testing.T) {
	var selector neighborsSelector

	selector = NEIGHBORS_ALL
	if len(selector.String()) <= 0 || selector.String() != "All" {
		t.Error("Unexpectedly retrieved empty string from Status object")
	}

	selector = NEIGHBORS_ORTHOGONAL
	if len(selector.String()) <= 0 || selector.String() != "Orthogonal" {
		t.Error("Unexpectedly retrieved empty string from pondselector object")
	}

	selector = NEIGHBORS_OBLIQUE
	if len(selector.String()) <= 0 || selector.String() != "Oblique" {
		t.Error("Unexpectedly retrieved empty string from pondselector object")
	}
}

func TestPondSettingInitialPatterns(t *testing.T) {
	pond, err := newpond(Dimensions{Height: 3, Width: 3}, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]Location, 3)
	initialLiving[0] = Location{X: 0, Y: 0}
	initialLiving[1] = Location{X: 1, Y: 1}
	initialLiving[2] = Location{X: 2, Y: 2}

	pond.SetOrganisms(initialLiving)

	// Check each expected value
	for _, loc := range initialLiving {
		if !pond.isOrganismAlive(loc) {
			t.Fatalf("Seed organism is not alive!: %s\n", loc.String())
		}
	}
}

func TestPondNeighborSelectionOrthogonal(t *testing.T) {
	pond, err := newpond(Dimensions{Height: 3, Width: 3}, NEIGHBORS_ORTHOGONAL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}

	expected := make([]Location, 4)
	expected[0] = Location{X: 1, Y: 0}
	expected[1] = Location{X: 0, Y: 1}
	expected[2] = Location{X: 2, Y: 1}
	expected[3] = Location{X: 1, Y: 2}

	var actual []Location
	actual, err = pond.GetNeighbors(Location{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("Unable to retrieve neighbors: %s\n", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("Retrieved %d neighbors but expected %d\n", len(actual), len(expected))
	}
}

func TestPondNeighborSelectionOblique(t *testing.T) {
	pond, err := newpond(Dimensions{Height: 3, Width: 3}, NEIGHBORS_OBLIQUE)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}

	expected := make([]Location, 4)
	expected[0] = Location{X: 0, Y: 0}
	expected[1] = Location{X: 2, Y: 0}
	expected[2] = Location{X: 0, Y: 2}
	expected[3] = Location{X: 2, Y: 2}

	var actual []Location
	actual, err = pond.GetNeighbors(Location{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("Unable to retrieve neighbors: %s\n", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("Retrieved %d neighbors but expected %d\n", len(actual), len(expected))
	}
}

func TestPondNeighborSelectionError(t *testing.T) {
	// t.Skip("Bad location should be the test")
	pond, err := newpond(Dimensions{Height: 1, Width: 1}, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}

	_, err = pond.GetNeighbors(Location{X: 2, Y: 2})
	if err == nil {
		t.Error("Did not encounter error when retrieving neighbors using bogus selector")
	}
}

func TestPondOrganismValue(t *testing.T) {
	expectedVal := 2
	pos := Location{X: 0, Y: 0}
	pond, err := newpond(Dimensions{Height: 1, Width: 1}, NEIGHBORS_ALL)
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
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newpond(dims, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]Location, 2)
	initialLiving[0] = Location{X: 0, Y: 0}
	initialLiving[1] = Location{X: 1, Y: 1}
	pond.SetOrganisms(initialLiving)

	numLiving := pond.GetNumLiving()

	if numLiving != len(initialLiving) {
		t.Fatalf("Retrieved actual %d living orgnamisms instead of expected %d\n", numLiving, len(initialLiving))
	}
}

func TestPondNeighborCountCalutation(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newpond(dims, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	// Create a pattern and call the pond's init function
	initialLiving := make([]Location, 2)
	initialLiving[0] = Location{X: 0, Y: 0}
	initialLiving[1] = Location{X: 1, Y: 1}
	pond.SetOrganisms(initialLiving)

	expectedNeighbors := make([]Location, 5)
	expectedNeighbors[0] = Location{X: 0, Y: 0}
	expectedNeighbors[1] = Location{X: 1, Y: 0}
	expectedNeighbors[2] = Location{X: 1, Y: 1}
	expectedNeighbors[3] = Location{X: 0, Y: 2}
	expectedNeighbors[4] = Location{X: 1, Y: 2}

	expectedNeighborCount := len(initialLiving)
	actualNeighborCount, actualNeighbors := pond.calculateNeighborCount(Location{X: 0, Y: 1})

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
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newpond(dims, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	if len(pond.String()) <= 0 {
		t.Error("Unexpectly retrieved empty string from pond string function")
	}
}

func TestPondEquals(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newpond(dims, NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	if !pond.Equals(pond) {
		t.Fatal("Pond Equals failed identity test")
	}

	pond2, err := newpond(dims, NEIGHBORS_OBLIQUE)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	if pond.Equals(pond2) {
		t.Fatal("Pond Equals returned true with ponds having different neighbor selectors")
	}
}
