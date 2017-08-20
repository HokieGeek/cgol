package life

import "testing"

func TestLocationString(t *testing.T) {
	loc := Location{X: 1, Y: 1}
	if len(loc.String()) == 0 {
		t.Error("The Location String() function unexpectedly returned an empty string")
	}
}

func TestDimensionsString(t *testing.T) {
	dims := Dimensions{Height: 1, Width: 1}
	if len(dims.String()) == 0 {
		t.Error("The Dimensions Strings() function unexpectedly returned an empty string")
	}
}

func TestDimensionsEquals(t *testing.T) {
	dims := Dimensions{Height: 42, Width: 42}

	// Check that they are equal
	if !dims.Equals(&dims) {
		t.Fatal("Stupid equals function failed identity check")
	}

	// Check where height is not equal
	notEqual := Dimensions{Height: 44, Width: 42}
	if dims.Equals(&notEqual) {
		t.Fatalf("Equals function says that %d is equal to %d\n", dims.Height, notEqual.Height)
	}

	// Check where width is not equal
	notEqual = Dimensions{Height: 42, Width: 22}
	if dims.Equals(&notEqual) {
		t.Fatalf("Equals function says that %d is equal to %d\n", dims.Height, notEqual.Height)
	}

	// Check where neither is equal
	notEqual = Dimensions{Height: 44, Width: 22}
	if dims.Equals(&notEqual) {
		t.Fatalf("Equals function says that %s is equal to %s\n", dims.String(), notEqual.String())
	}
}

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
	pond, err := newPond(Dimensions{Height: 3, Width: 3}, newTracker(), NEIGHBORS_ALL)
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
	pond, err := newPond(Dimensions{Height: 3, Width: 3}, newTracker(), NEIGHBORS_ORTHOGONAL)
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
	pond, err := newPond(Dimensions{Height: 3, Width: 3}, newTracker(), NEIGHBORS_OBLIQUE)
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
	pond, err := newPond(Dimensions{Height: 1, Width: 1}, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}

	_, err = pond.GetNeighbors(Location{X: 2, Y: 2})
	if err == nil {
		t.Error("Did not encounter error when retrieving neighbors using bogus selector")
	}
}

/* TODO
func TestPondOrganismValue(t *testing.T) {
	expectedVal := 2
	pos := Location{X: 0, Y: 0}
	pond, err := newPond(Dimensions{Height: 1, Width: 1}, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	pond.setOrganismValue(pos, expectedVal)

	actualVal := pond.GetOrganismValue(pos)

	if actualVal != expectedVal {
		t.Fatalf("Retrieved actual value %d instead of expected value %d\n", actualVal, expectedVal)
	}
}
*/

func TestPondString(t *testing.T) {
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newPond(dims, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}
	if len(pond.String()) <= 0 {
		t.Error("Unexpectly retrieved empty string from pond string function")
	}
}

func TestPondEquals(t *testing.T) {
	t.Skip("whoops")
	dims := Dimensions{Height: 3, Width: 3}
	pond, err := newPond(dims, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	if !pond.Equals(pond) {
		t.Fatal("Pond Equals failed identity test")
	}

	pond2, err := newPond(dims, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatal("Unable to create pond")
	}

	if pond.Equals(pond2) {
		t.Fatal("Pond Equals returned true with ponds having different neighbor selectors")
	}
}

/*
func testBoardNeighbors(t *testing.T, expected []Location, actual []Location) {
	// Check the results
	if len(actual) != len(expected) {
		t.Fatalf("Number of neighbors (%d) does not match expected (%d)\n", len(actual), len(expected))
	}

	// Check that all expected locations are in the actual list
	for _, expectedLoc := range expected {
		found := false
		for _, actualLoc := range actual {
			if expectedLoc.Equals(&actualLoc) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Did not find location %s in actual list\n", expectedLoc.String())
		}
	}

	// Check that the actual list doesn't have any unexpected locations
	for _, actualLoc := range actual {
		found := false
		for _, expectedLoc := range expected {
			if expectedLoc.Equals(&actualLoc) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Found location %s in actual list that was not expected\n", actualLoc.String())
		}
	}
}

func TestBoardGetOrthogonalNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]Location, 4)
	expected[0] = Location{X: 1, Y: 0}
	expected[1] = Location{X: 0, Y: 1}
	expected[2] = Location{X: 2, Y: 1}
	expected[3] = Location{X: 1, Y: 2}

	// Initialize a board
	dims := Dimensions{Height: 3, Width: 3}
	board, err := newBoard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetOrthogonalNeighbors(Location{X: 1, Y: 1})

	testBoardNeighbors(t, expected, actual)
}

func TestBoardGetObliqueNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]Location, 4)
	expected[0] = Location{X: 0, Y: 0}
	expected[1] = Location{X: 2, Y: 0}
	expected[2] = Location{X: 0, Y: 2}
	expected[3] = Location{X: 2, Y: 2}

	// Initialize a board
	dims := Dimensions{Height: 3, Width: 3}
	board, err := newBoard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetObliqueNeighbors(Location{X: 1, Y: 1})

	testBoardNeighbors(t, expected, actual)
}

func TestBoardGetAllNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]Location, 0)
	for i := 0; i < 4; i++ {
		switch i {
		case 0, 2:
			for j := 0; j < 3; j++ {
				expected = append(expected, Location{X: j, Y: i})
			}
		case 1:
			expected = append(expected, Location{X: 0, Y: 1})
			expected = append(expected, Location{X: 2, Y: 1})
		}
	}

	// Initialize a board
	dims := Dimensions{Height: 3, Width: 3}
	board, err := newBoard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetAllNeighbors(Location{X: 1, Y: 1})

	testBoardNeighbors(t, expected, actual)
}
*/

// vim: set foldmethod=marker:
