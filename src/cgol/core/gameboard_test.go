package cgol

import (
	"math/rand"
	"testing"
	"time"
)

func TestGameboardLocationString(t *testing.T) {
	loc := GameboardLocation{X: 1, Y: 1}
	if len(loc.String()) == 0 {
		t.Error("The GameboardLocation Strings() function unexpectedly returned an empty string")
	}
}

func TestGameboardDimsString(t *testing.T) {
	dims := GameboardDims{Height: 1, Width: 1}
	if len(dims.String()) == 0 {
		t.Error("The GameboardDims Strings() function unexpectedly returned an empty string")
	}
}

func TestGameboardCreation(t *testing.T) {
	// Create a gameboard of random size
	rand.Seed(time.Now().UnixNano())
	size := GameboardDims{Height: rand.Intn(100) + 1, Width: rand.Intn(100) + 1}
	gameboard, err := NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}

	// Test that the values were stored correctly
	if gameboard.Dims.Height != size.Height {
		t.Error("Height not stored correctly")
	}
	if gameboard.Dims.Width != size.Width {
		t.Error("Width not stored correctly")
	}

	// Now check the size of the gameboard itself
	snapshot := gameboard.getSnapshot()
	if len(snapshot) != size.Height {
		t.Fatal("The gameboard is not the correct number of rows")
	}
	for row := 0; row < size.Height; row++ {
		if len(snapshot[row]) != size.Width {
			t.Fatal("The gameboard is not the correct number of columns")
		}
	}
}

func TestGameboardCreateWithErrors(t *testing.T) {
	// No height
	gameboard, err := NewGameboard(GameboardDims{Height: 0, Width: 1})
	if err == nil {
		t.Error("Creating a gameboard with 0 height did not return an error")
	}
	if gameboard != nil {
		t.Error("Creating a gameboard with 0 height returned a valid object")
	}

	// No width
	gameboard, err = NewGameboard(GameboardDims{Height: 1, Width: 0})
	if err == nil {
		t.Error("Creating a gameboard with 0 width did not return an error")
	}
	if gameboard != nil {
		t.Error("Creating a gameboard with 0 width returned a valid object")
	}

	// Both <0
	gameboard, err = NewGameboard(GameboardDims{Height: -1, Width: -1})
	if err == nil {
		t.Error("Creating a gameboard with width and height less than 0 did not return an error")
	}
	if gameboard != nil {
		t.Error("Creating a gameboard with width and height less than 0 returned a valid object")
	}
}

func TestGameboardSetValue(t *testing.T) {
	// Create the test gameboard
	dims := GameboardDims{Height: 1, Width: 1}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}

	// Set a good value and retrieve it
	loc := GameboardLocation{X: 0, Y: 0}
	testVal := 42
	gameboard.SetValue(loc, testVal)

	realVal, err := gameboard.GetValue(loc)
	if err != nil {
		t.Fatalf("Gameboard.GetValue() unexpectly returned an error: %s\n", err)
	}
	if realVal != testVal {
		t.Fatalf("Found value %d but expected %d\n", realVal, testVal)
	}
}

func TestGameboardSetValueOutOfBounds(t *testing.T) {
	// Create the test gameboard
	dims := GameboardDims{Height: 1, Width: 1}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds
	loc := GameboardLocation{X: 2, Y: 0}
	testVal := 42
	err = gameboard.SetValue(loc, testVal)

	if err == nil {
		t.Error("Gameboard did not return an error when setting a value at an out-of-bounds location")
	}
}

func TestGameboardGetValueOutOfBounds(t *testing.T) {
	// Create the test gameboard
	dims := GameboardDims{Height: 1, Width: 1}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds
	loc := GameboardLocation{X: 2, Y: 0}
	_, err = gameboard.GetValue(loc)

	if err == nil {
		t.Error("Gameboard did not return an error when retrieving a value at an out-of-bounds location")
	}
}

func testGameboardNeighbors(t *testing.T, expected []GameboardLocation, actual []GameboardLocation) {
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

func TestGameboardGetOrthogonalNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]GameboardLocation, 4)
	expected[0] = GameboardLocation{X: 1, Y: 0}
	expected[1] = GameboardLocation{X: 0, Y: 1}
	expected[2] = GameboardLocation{X: 2, Y: 1}
	expected[3] = GameboardLocation{X: 1, Y: 2}

	// Initialize a gameboard
	dims := GameboardDims{Height: 3, Width: 3}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		gameboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := gameboard.GetOrthogonalNeighbors(GameboardLocation{X: 1, Y: 1})

	testGameboardNeighbors(t, expected, actual)
}

func TestGameboardGetObliqueNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]GameboardLocation, 4)
	expected[0] = GameboardLocation{X: 0, Y: 0}
	expected[1] = GameboardLocation{X: 2, Y: 0}
	expected[2] = GameboardLocation{X: 0, Y: 2}
	expected[3] = GameboardLocation{X: 2, Y: 2}

	// Initialize a gameboard
	dims := GameboardDims{Height: 3, Width: 3}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		gameboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := gameboard.GetObliqueNeighbors(GameboardLocation{X: 1, Y: 1})

	testGameboardNeighbors(t, expected, actual)
}

func TestGameboardGetAllNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]GameboardLocation, 0)
	for i := 0; i < 4; i++ {
		switch i {
		case 0, 2:
			for j := 0; j < 3; j++ {
				expected = append(expected, GameboardLocation{X: j, Y: i})
			}
		case 1:
			expected = append(expected, GameboardLocation{X: 0, Y: 1})
			expected = append(expected, GameboardLocation{X: 2, Y: 1})
		}
	}

	// Initialize a gameboard
	dims := GameboardDims{Height: 3, Width: 3}
	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		gameboard.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := gameboard.GetAllNeighbors(GameboardLocation{X: 1, Y: 1})

	testGameboardNeighbors(t, expected, actual)
}

func TestGameboardGetSnapshot(t *testing.T) {
	// Initialize gameboards
	dims := GameboardDims{Height: 2, Width: 2}
	locations := make([]GameboardLocation, 2)
	locations[0] = GameboardLocation{X: 1, Y: 0}
	locations[1] = GameboardLocation{X: 0, Y: 1}

	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		gameboard.SetValue(loc, 0)
	}

	// Test that the snapshots work?
	snapshot := gameboard.getSnapshot()
	for row := dims.Height - 1; row >= 0; row-- {
		for col := dims.Width - 1; col >= 0; col-- {
			location := GameboardLocation{X: col, Y: row}
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

func TestGameboardEquals(t *testing.T) {
	// Initialize one gameboard
	dims := GameboardDims{Height: 2, Width: 2}
	locations := make([]GameboardLocation, 2)
	locations[0] = GameboardLocation{X: 1, Y: 0}
	locations[1] = GameboardLocation{X: 0, Y: 1}

	gameboard, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		gameboard.SetValue(loc, 0)
	}

	// Test against itself
	if !gameboard.Equals(gameboard) {
		t.Fatal("Equality function is definitely broken")
	}
}

func TestGameboardNotEquals(t *testing.T) {
	dims := GameboardDims{Height: 2, Width: 2}

	// Create one gameboard
	locationsLeft := make([]GameboardLocation, 2)
	locationsLeft[0] = GameboardLocation{X: 1, Y: 0}
	locationsLeft[1] = GameboardLocation{X: 0, Y: 1}

	gameboardLeft, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsLeft {
		gameboardLeft.SetValue(loc, 0)
	}

	// Create the other gameboard
	locationsRight := make([]GameboardLocation, 2)
	locationsRight[0] = GameboardLocation{X: 1, Y: 1}
	locationsRight[1] = GameboardLocation{X: 0, Y: 0}

	gameboardRight, err := NewGameboard(dims)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsRight {
		gameboardRight.SetValue(loc, 0)
	}

	// The two should not be equal to each other
	if gameboardLeft.Equals(gameboardRight) {
		t.Fatalf("Equality function is definitely broken. This gameboard: \n%s\nshould not be equal to this one:\n%s\n", gameboardLeft.String(), gameboardRight.String())
	}
}
