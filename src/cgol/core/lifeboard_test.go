package cgol

import (
	"math/rand"
	"testing"
	"time"
)

func TestLifeboardLocationString(t *testing.T) {
	loc := LifeboardLocation{X: 1, Y: 1}
	if len(loc.String()) == 0 {
		t.Error("The LifeboardLocation String() function unexpectedly returned an empty string")
	}
}

func TestLifeboardDimsString(t *testing.T) {
	dims := LifeboardDims{Height: 1, Width: 1}
	if len(dims.String()) == 0 {
		t.Error("The LifeboardDims Strings() function unexpectedly returned an empty string")
	}
}

func TestLifeboardCreation(t *testing.T) {
	// Create a lifeboard of random size
	rand.Seed(time.Now().UnixNano())
	size := LifeboardDims{Height: rand.Intn(100) + 1, Width: rand.Intn(100) + 1}
	lifeboard, err := newLifeboard(size)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", size.String())
	}

	// Test that the values were stored correctly
	if lifeboard.Dims.Height != size.Height {
		t.Error("Height not stored correctly")
	}
	if lifeboard.Dims.Width != size.Width {
		t.Error("Width not stored correctly")
	}

	// Now check the size of the lifeboard itself
	snapshot := lifeboard.getSnapshot()
	if len(snapshot) != size.Height {
		t.Fatal("The lifeboard is not the correct number of rows")
	}
	for row := 0; row < size.Height; row++ {
		if len(snapshot[row]) != size.Width {
			t.Fatal("The lifeboard is not the correct number of columns")
		}
	}
}

func TestLifeboardCreateWithErrors(t *testing.T) {
	// No height
	lifeboard, err := newLifeboard(LifeboardDims{Height: 0, Width: 1})
	if err == nil {
		t.Error("Creating a lifeboard with 0 height did not return an error")
	}
	if lifeboard != nil {
		t.Error("Creating a lifeboard with 0 height returned a valid object")
	}

	// No width
	lifeboard, err = newLifeboard(LifeboardDims{Height: 1, Width: 0})
	if err == nil {
		t.Error("Creating a lifeboard with 0 width did not return an error")
	}
	if lifeboard != nil {
		t.Error("Creating a lifeboard with 0 width returned a valid object")
	}

	// Both <0
	lifeboard, err = newLifeboard(LifeboardDims{Height: -1, Width: -1})
	if err == nil {
		t.Error("Creating a lifeboard with width and height less than 0 did not return an error")
	}
	if lifeboard != nil {
		t.Error("Creating a lifeboard with width and height less than 0 returned a valid object")
	}
}

func TestLifeboardSetValue(t *testing.T) {
	// Create the test lifeboard
	dims := LifeboardDims{Height: 1, Width: 1}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}

	// Set a good value and retrieve it
	loc := LifeboardLocation{X: 0, Y: 0}
	testVal := 42
	lifeboard.SetValue(loc, testVal)

	realVal, err := lifeboard.GetValue(loc)
	if err != nil {
		t.Fatalf("Lifeboard.GetValue() unexpectly returned an error: %s\n", err)
	}
	if realVal != testVal {
		t.Fatalf("Found value %d but expected %d\n", realVal, testVal)
	}
}

func TestLifeboardSetValueOutOfBounds(t *testing.T) {
	// Create the test lifeboard
	dims := LifeboardDims{Height: 1, Width: 1}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds (X)
	loc := LifeboardLocation{X: 2, Y: 0}
	testVal := 42
	err = lifeboard.SetValue(loc, testVal)
	if err == nil {
		t.Error("lifeboard did not return an error when setting a value at an out-of-bounds location")
	}

	// Check for out of bounds (Y)
	loc = LifeboardLocation{X: 0, Y: 2}
	err = lifeboard.SetValue(loc, testVal)
	if err == nil {
		t.Error("lifeboard did not return an error when setting a value at an out-of-bounds location")
	}
}

func TestLifeboardGetValueOutOfBounds(t *testing.T) {
	// Create the test lifeboard
	dims := LifeboardDims{Height: 1, Width: 1}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds (X)
	loc := LifeboardLocation{X: 2, Y: 0}
	_, err = lifeboard.GetValue(loc)
	if err == nil {
		t.Error("lifeboard did not return an error when retrieving a value at an out-of-bounds location")
	}

	// Check for out of bounds (Y)
	loc = LifeboardLocation{X: 0, Y: 2}
	_, err = lifeboard.GetValue(loc)
	if err == nil {
		t.Error("lifeboard did not return an error when retrieving a value at an out-of-bounds location")
	}
}

func testLifeboardNeighbors(t *testing.T, expected []LifeboardLocation, actual []LifeboardLocation) {
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

func TestLifeboardGetOrthogonalNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]LifeboardLocation, 4)
	expected[0] = LifeboardLocation{X: 1, Y: 0}
	expected[1] = LifeboardLocation{X: 0, Y: 1}
	expected[2] = LifeboardLocation{X: 2, Y: 1}
	expected[3] = LifeboardLocation{X: 1, Y: 2}

	// Initialize a lifeboard
	dims := LifeboardDims{Height: 3, Width: 3}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		lifeboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := lifeboard.GetOrthogonalNeighbors(LifeboardLocation{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetObliqueNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]LifeboardLocation, 4)
	expected[0] = LifeboardLocation{X: 0, Y: 0}
	expected[1] = LifeboardLocation{X: 2, Y: 0}
	expected[2] = LifeboardLocation{X: 0, Y: 2}
	expected[3] = LifeboardLocation{X: 2, Y: 2}

	// Initialize a lifeboard
	dims := LifeboardDims{Height: 3, Width: 3}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		lifeboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := lifeboard.GetObliqueNeighbors(LifeboardLocation{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetAllNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]LifeboardLocation, 0)
	for i := 0; i < 4; i++ {
		switch i {
		case 0, 2:
			for j := 0; j < 3; j++ {
				expected = append(expected, LifeboardLocation{X: j, Y: i})
			}
		case 1:
			expected = append(expected, LifeboardLocation{X: 0, Y: 1})
			expected = append(expected, LifeboardLocation{X: 2, Y: 1})
		}
	}

	// Initialize a lifeboard
	dims := LifeboardDims{Height: 3, Width: 3}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		lifeboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := lifeboard.GetAllNeighbors(LifeboardLocation{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetSnapshot(t *testing.T) {
	// Initialize lifeboards
	dims := LifeboardDims{Height: 2, Width: 2}
	locations := make([]LifeboardLocation, 2)
	locations[0] = LifeboardLocation{X: 1, Y: 0}
	locations[1] = LifeboardLocation{X: 0, Y: 1}

	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		lifeboard.SetValue(loc, 0)
	}

	// Test that the snapshots work?
	snapshot := lifeboard.getSnapshot()
	for row := dims.Height - 1; row >= 0; row-- {
		for col := dims.Width - 1; col >= 0; col-- {
			location := LifeboardLocation{X: col, Y: row}
			found := false
			for _, loc := range locations {
				if loc.Equals(&location) {
					found = true
					break
				}
			}
			if snapshot[row][col] == 0 && !found {
				t.Fatalf("Did not find location %s in initial list\n", location.String())
			} else if snapshot[row][col] != 0 && found {
				t.Fatalf("Found location %s in snapshot with a value that was not expected\n", location.String())
			}
		}
	}
}

func TestLifeboardEquals(t *testing.T) {
	// Initialize one lifeboard
	dims := LifeboardDims{Height: 2, Width: 2}
	locations := make([]LifeboardLocation, 2)
	locations[0] = LifeboardLocation{X: 1, Y: 0}
	locations[1] = LifeboardLocation{X: 0, Y: 1}

	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		lifeboard.SetValue(loc, 0)
	}

	// Test against itself
	if !lifeboard.Equals(lifeboard) {
		t.Fatal("Equality function is definitely broken")
	}
}

func TestLifeboardNotEquals(t *testing.T) {
	dims := LifeboardDims{Height: 2, Width: 2}

	// Create one lifeboard
	locationsLeft := make([]LifeboardLocation, 2)
	locationsLeft[0] = LifeboardLocation{X: 1, Y: 0}
	locationsLeft[1] = LifeboardLocation{X: 0, Y: 1}

	lifeboardLeft, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsLeft {
		lifeboardLeft.SetValue(loc, 0)
	}

	// Create the other lifeboard
	locationsRight := make([]LifeboardLocation, 2)
	locationsRight[0] = LifeboardLocation{X: 1, Y: 1}
	locationsRight[1] = LifeboardLocation{X: 0, Y: 0}

	lifeboardRight, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsRight {
		lifeboardRight.SetValue(loc, 0)
	}

	// The two should not be equal to each other
	if lifeboardLeft.Equals(lifeboardRight) {
		t.Fatalf("Equality function is definitely broken. This lifeboard: \n%s\nshould not be equal to this one:\n%s\n", lifeboardLeft.String(), lifeboardRight.String())
	}
}

func TestLifeboardString(t *testing.T) {
	// Create the test lifeboard
	dims := LifeboardDims{Height: 2, Width: 2}
	lifeboard, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("lifeboard of size %s could not be created\n", dims.String())
	}
	lifeboard.SetValue(LifeboardLocation{X: 0, Y: 0}, 0)

	// Now test the string call
	if len(lifeboard.String()) == 0 {
		t.Error("The lifeboard String() function unexpectedly returned an empty string")
	}
}
