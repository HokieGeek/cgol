package life

import "testing"

func TestGetCountsForDimensions(t *testing.T) {
	// Board size
	boardDims := Dimensions{Width: 7, Height: 3}

	// Pattern size
	patternDims := Dimensions{Height: 2, Width: 2}

	// Expected values
	expectedNumPerRow := 3
	expectedNumPerCol := 1

	// Execute the function being tested
	actualNumPerRow, actualNumPerCol := getCountsForDimensions(boardDims, patternDims)

	// Now test the results
	if expectedNumPerRow != actualNumPerRow {
		t.Fatalf("Expected %d patterns per row but got %d instead\n", expectedNumPerRow, actualNumPerRow)
	}
	if expectedNumPerCol != actualNumPerCol {
		t.Fatalf("Expected %d patterns per column but got %d instead\n", expectedNumPerCol, actualNumPerCol)
	}
}

func TestGetRepeatingPattern(t *testing.T) {
	// Board size
	boardDims := Dimensions{Width: 4, Height: 4}

	// Pattern size
	patternDims := Dimensions{Height: 1, Width: 1}

	// What I expect (a checker board)
	expectedLocations := make([]Location, 0)
	for row := 0; row < boardDims.Height; row++ {
		switch row {
		case 0, 2:
			expectedLocations = append(expectedLocations, Location{X: 0, Y: row})
			expectedLocations = append(expectedLocations, Location{X: 2, Y: row})
		case 1, 3:
			expectedLocations = append(expectedLocations, Location{X: 1, Y: row})
			expectedLocations = append(expectedLocations, Location{X: 3, Y: row})
		}
	}

	// Run the function being tested
	actualLocations := getRepeatingPattern(boardDims, patternDims, Location{},
		func(initialLiving *[]Location, currentX int, currentY int) {
			if (currentY%2 == 0 && currentX%2 == 0) || (currentY%2 != 0 && currentX%2 != 0) {
				*initialLiving = append(*initialLiving, Location{X: currentX, Y: currentY})
			}
		})

	// Now test the results
	if len(actualLocations) != len(expectedLocations) {
		t.Fatalf("Function returned %d locations but expected %d\n", len(actualLocations), len(expectedLocations))
	}

	// Check that all expected locations are in the actual list
	for _, expectedLoc := range expectedLocations {
		found := false
		for _, actualLoc := range actualLocations {
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
	for _, actualLoc := range actualLocations {
		found := false
		for _, expectedLoc := range expectedLocations {
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

func TestPatternOffsetting(t *testing.T) {
	// Create the expected list
	expected := make([]Location, 3)
	expected[0] = Location{X: 2, Y: 3}
	expected[1] = Location{X: 3, Y: 3}
	expected[2] = Location{X: 4, Y: 3}

	// Call blinker with an offset
	actual := Blinkers(Dimensions{Height: 6, Width: 6}, Location{X: 2, Y: 2})

	if len(actual) != len(expected) {
		t.Fatalf("Length of seed %d does not match expected %d\n", len(actual), len(expected))
	}

	// t.Logf("Actual: %v\n", actual)
	// t.Logf("Expected: %v\n", expected)

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
