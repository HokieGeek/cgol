package cgol

import "testing"

func testProcessorSimultaneousRulesConwayLife(t *testing.T,
	size GameboardDims,
	init func(GameboardDims) []GameboardLocation,
	expected *Gameboard,
	numGenerations int) {

	// Build the initial pond
	pond := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	pond.init(init(size))

	// Go through one generation
	for i := 0; i < numGenerations; i++ {
		SimultaneousProcessor(pond, RulesConwayLife)
	}

	// Compare the pond with the expected version
	if !pond.gameboard.Equals(expected) {
		t.Fatalf("Actual gameboard:\n%s\ndoes not match expected:\n%s\n", pond.gameboard.String(), expected.String())
	}
}

func TestProcessorSimultaneousRulesConwayLifeBlinker(t *testing.T) {
	size := GameboardDims{Height: 3, Width: 3}

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 2, Y: 1}, 0)

	testProcessorSimultaneousRulesConwayLife(t, size, Blinkers, expected, 1)
}

func TestProcessorSimultaneousRulesConwayLifeToad(t *testing.T) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetValue(GameboardLocation{X: 2, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 3, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 0, Y: 2}, 0)
	expected.SetValue(GameboardLocation{X: 3, Y: 2}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 3}, 0)

	testProcessorSimultaneousRulesConwayLife(t, size, Toads, expected, 1)
}

func TestProcessorSimultaneousRulesConwayLifeBeacon(t *testing.T) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetValue(GameboardLocation{X: 0, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 2, Y: 2}, 0)
	expected.SetValue(GameboardLocation{X: 3, Y: 2}, 0)
	expected.SetValue(GameboardLocation{X: 2, Y: 3}, 0)
	expected.SetValue(GameboardLocation{X: 3, Y: 3}, 0)

	testProcessorSimultaneousRulesConwayLife(t, size, Beacons, expected, 1)
}

func TestProcessorSimultaneousRulesConwayLifeBlock(t *testing.T) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetValue(GameboardLocation{X: 0, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 0}, 0)
	expected.SetValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetValue(GameboardLocation{X: 1, Y: 1}, 0)

	testProcessorSimultaneousRulesConwayLife(t, size, Blocks, expected, 1)
}
