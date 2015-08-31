package cgol

import (
	"math/rand"
	"testing"
	"time"
)

func TestGameboardCreation(t *testing.T) {
	// Create a gameboard of random size
	rand.Seed(time.Now().UnixNano())
	size := GameboardDims{Height: rand.Intn(100), Width: rand.Intn(100)}
	gameboard := NewGameboard(size)

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

func TestGameboardSettingValue(t *testing.T) {
	// Create the test gameboard
	gameboard := NewGameboard(GameboardDims{Height: 1, Width: 1})

	// Set a good value and retrieve it
	loc := GameboardLocation{X: 0, Y: 0}
	testVal := 42
	gameboard.SetValue(loc, testVal)

	realVal := gameboard.GetValue(loc)
	if realVal != testVal {
		t.Fatalf("Found value %d but expected %d\n", realVal, testVal)
	}

	// Attempt to retrieve from a non-existent location
	// TODO: gameboard.GetValue(GameboardLocation{X: 1, Y: 1})
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
	gameboard := NewGameboard(GameboardDims{Height: 3, Width: 3})
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
	gameboard := NewGameboard(GameboardDims{Height: 3, Width: 3})
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
	gameboard := NewGameboard(GameboardDims{Height: 3, Width: 3})
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

	gameboard := NewGameboard(dims)
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

	gameboard := NewGameboard(dims)
	for _, loc := range locations {
		gameboard.SetValue(loc, 0)
	}

	// Test against itself
	if !gameboard.Equals(gameboard) {
		t.Fatal("Equality function is definitely broken")
	}
}

func TestGameboardNotEquals(t *testing.T) {
	t.Skip("TODO")
}
