package life

import (
	"math/rand"
	"testing"
	"time"
)

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

func TestLifeboardCreation(t *testing.T) {
	// Create a board of random size
	rand.Seed(time.Now().UnixNano())
	size := Dimensions{Height: rand.Intn(100) + 1, Width: rand.Intn(100) + 1}
	board, err := newLifeboard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}

	// Test that the values were stored correctly
	if board.Dims.Height != size.Height {
		t.Error("Height not stored correctly")
	}
	if board.Dims.Width != size.Width {
		t.Error("Width not stored correctly")
	}

	// Now check the size of the board itself
	snapshot := board.getSnapshot()
	if len(snapshot) != size.Height {
		t.Fatal("The board is not the correct number of rows")
	}
	for row := 0; row < size.Height; row++ {
		if len(snapshot[row]) != size.Width {
			t.Fatal("The board is not the correct number of columns")
		}
	}
}

func TestLifeboardCreateWithErrors(t *testing.T) {
	// No height
	board, err := newLifeboard(Dimensions{Height: 0, Width: 1})
	if err == nil {
		t.Error("Creating a board with 0 height did not return an error")
	}
	if board != nil {
		t.Error("Creating a board with 0 height returned a valid object")
	}

	// No width
	board, err = newLifeboard(Dimensions{Height: 1, Width: 0})
	if err == nil {
		t.Error("Creating a board with 0 width did not return an error")
	}
	if board != nil {
		t.Error("Creating a board with 0 width returned a valid object")
	}

	// Both <0
	board, err = newLifeboard(Dimensions{Height: -1, Width: -1})
	if err == nil {
		t.Error("Creating a board with width and height less than 0 did not return an error")
	}
	if board != nil {
		t.Error("Creating a board with width and height less than 0 returned a valid object")
	}
}

func TestLifeboardSetValue(t *testing.T) {
	// Create the test board
	dims := Dimensions{Height: 1, Width: 1}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}

	// Set a good value and retrieve it
	loc := Location{X: 0, Y: 0}
	testVal := 42
	board.SetValue(loc, testVal)

	realVal, err := board.GetValue(loc)
	if err != nil {
		t.Fatalf("Lifeboard.GetValue() unexpectly returned an error: %s\n", err)
	}
	if realVal != testVal {
		t.Fatalf("Found value %d but expected %d\n", realVal, testVal)
	}
}

func TestLifeboardSetValueOutOfBounds(t *testing.T) {
	// Create the test board
	dims := Dimensions{Height: 1, Width: 1}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds (X)
	loc := Location{X: 2, Y: 0}
	testVal := 42
	err = board.SetValue(loc, testVal)
	if err == nil {
		t.Error("board did not return an error when setting a value at an out-of-bounds location")
	}

	// Check for out of bounds (Y)
	loc = Location{X: 0, Y: 2}
	err = board.SetValue(loc, testVal)
	if err == nil {
		t.Error("board did not return an error when setting a value at an out-of-bounds location")
	}
}

func TestLifeboardGetValueOutOfBounds(t *testing.T) {
	// Create the test board
	dims := Dimensions{Height: 1, Width: 1}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}

	// Check for out of bounds (X)
	loc := Location{X: 2, Y: 0}
	_, err = board.GetValue(loc)
	if err == nil {
		t.Error("board did not return an error when retrieving a value at an out-of-bounds location")
	}

	// Check for out of bounds (Y)
	loc = Location{X: 0, Y: 2}
	_, err = board.GetValue(loc)
	if err == nil {
		t.Error("board did not return an error when retrieving a value at an out-of-bounds location")
	}
}

func testLifeboardNeighbors(t *testing.T, expected []Location, actual []Location) {
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
	expected := make([]Location, 4)
	expected[0] = Location{X: 1, Y: 0}
	expected[1] = Location{X: 0, Y: 1}
	expected[2] = Location{X: 2, Y: 1}
	expected[3] = Location{X: 1, Y: 2}

	// Initialize a board
	dims := Dimensions{Height: 3, Width: 3}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetOrthogonalNeighbors(Location{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetObliqueNeighbors(t *testing.T) {
	// Build list of expected locations
	expected := make([]Location, 4)
	expected[0] = Location{X: 0, Y: 0}
	expected[1] = Location{X: 2, Y: 0}
	expected[2] = Location{X: 0, Y: 2}
	expected[3] = Location{X: 2, Y: 2}

	// Initialize a board
	dims := Dimensions{Height: 3, Width: 3}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetObliqueNeighbors(Location{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetAllNeighbors(t *testing.T) {
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
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for i := 0; i < len(expected); i++ {
		board.SetValue(expected[i], 0)
	}

	// Retrieve neighbors
	actual := board.GetAllNeighbors(Location{X: 1, Y: 1})

	testLifeboardNeighbors(t, expected, actual)
}

func TestLifeboardGetSnapshot(t *testing.T) {
	// Initialize boards
	dims := Dimensions{Height: 2, Width: 2}
	locations := make([]Location, 2)
	locations[0] = Location{X: 1, Y: 0}
	locations[1] = Location{X: 0, Y: 1}

	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		board.SetValue(loc, 0)
	}

	// Test that the snapshots work?
	snapshot := board.getSnapshot()
	for row := dims.Height - 1; row >= 0; row-- {
		for col := dims.Width - 1; col >= 0; col-- {
			location := Location{X: col, Y: row}
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
	// Initialize one board
	dims := Dimensions{Height: 2, Width: 2}
	locations := make([]Location, 2)
	locations[0] = Location{X: 1, Y: 0}
	locations[1] = Location{X: 0, Y: 1}

	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for _, loc := range locations {
		board.SetValue(loc, 0)
	}

	// Test against itself
	if !board.Equals(board) {
		t.Fatal("Equality function is definitely broken")
	}
}

func TestLifeboardNotEquals(t *testing.T) {
	dims := Dimensions{Height: 2, Width: 2}

	// Create one board
	locationsLeft := make([]Location, 2)
	locationsLeft[0] = Location{X: 1, Y: 0}
	locationsLeft[1] = Location{X: 0, Y: 1}

	boardLeft, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsLeft {
		boardLeft.SetValue(loc, 0)
	}

	// Create the other board
	locationsRight := make([]Location, 2)
	locationsRight[0] = Location{X: 1, Y: 1}
	locationsRight[1] = Location{X: 0, Y: 0}

	boardRight, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	for _, loc := range locationsRight {
		boardRight.SetValue(loc, 0)
	}

	// The two should not be equal to each other
	if boardLeft.Equals(boardRight) {
		t.Fatalf("Equality function is definitely broken. This board: \n%s\nshould not be equal to this one:\n%s\n", boardLeft.String(), boardRight.String())
	}
}

func TestLifeboardString(t *testing.T) {
	// Create the test board
	dims := Dimensions{Height: 2, Width: 2}
	board, err := newLifeboard(dims)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", dims.String())
	}
	board.SetValue(Location{X: 0, Y: 0}, 0)

	// Now test the string call
	if len(board.String()) == 0 {
		t.Error("The board String() function unexpectedly returned an empty string")
	}
}
