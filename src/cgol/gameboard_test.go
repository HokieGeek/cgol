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
	snapshot := gameboard.getGameboardSnapshot()
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
	gameboard.SetGameboardValue(loc, testVal)

	realVal := gameboard.GetGameboardValue(loc)
	if realVal != testVal {
		t.Fatalf("Found value %d but expected %d\n", realVal, testVal)
	}

	// Attempt to retrieve from a non-existent location
	// TODO: gameboard.GetGameboardValue(GameboardLocation{X: 1, Y: 1})
}
