package life

import "testing"

func TestGetCountsForDimensions(t *testing.T) {
	// Board size
	dims := LifeboardDims{Width: 7, Height: 3}

	// Pattern size
	patternWidth := 2
	patternHeight := 2

	// Expected values
	expectedNumPerRow := 3
	expectedNumPerCol := 1

	// Execute the function being tested
	actualNumPerRow, actualNumPerCol := getCountsForDimensions(dims, patternWidth, patternHeight)

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
	dims := LifeboardDims{Width: 4, Height: 4}

	// Pattern size
	height := 1
	width := 1

	// What I expect (a checker board)
	expectedLocations := make([]LifeboardLocation, 0)
	for row := 0; row < dims.Height; row++ {
		switch row {
		case 0, 2:
			expectedLocations = append(expectedLocations, LifeboardLocation{X: 0, Y: row})
			expectedLocations = append(expectedLocations, LifeboardLocation{X: 2, Y: row})
		case 1, 3:
			expectedLocations = append(expectedLocations, LifeboardLocation{X: 1, Y: row})
			expectedLocations = append(expectedLocations, LifeboardLocation{X: 3, Y: row})
		}
	}

	// Run the function being tested
	actualLocations := getRepeatingPattern(dims, height, width,
		func(initialLiving *[]LifeboardLocation, currentX int, currentY int) {
			if (currentY%2 == 0 && currentX%2 == 0) || (currentY%2 != 0 && currentX%2 != 0) {
				*initialLiving = append(*initialLiving, LifeboardLocation{X: currentX, Y: currentY})
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
