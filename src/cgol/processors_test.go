package cgol

import "testing"

func TestProcessorSimultaneousStandardBlinker(t *testing.T) {
	size := GameboardDims{Height: 3, Width: 3}

	// Build the initial pond
	pond := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	initial := make([]GameboardLocation, 3)
	initial[0] = GameboardLocation{X: 1, Y: 0}
	initial[1] = GameboardLocation{X: 1, Y: 1}
	initial[2] = GameboardLocation{X: 1, Y: 2}
	pond.init(initial)

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetGameboardValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 1, Y: 1}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 2, Y: 1}, 0)

	// Go through one generation
	SimultaneousProcessor(pond, Standard)

	// Compare the pond with the expected version
	if !pond.gameboard.Equals(expected) {
		t.FailNow()
	}
}

func TestProcessorSimultaneousStandardToad(t *testing.T) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the initial pond
	pond := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	initial := make([]GameboardLocation, 6)
	initial[0] = GameboardLocation{X: 1, Y: 1}
	initial[1] = GameboardLocation{X: 2, Y: 1}
	initial[2] = GameboardLocation{X: 3, Y: 1}
	initial[3] = GameboardLocation{X: 0, Y: 2}
	initial[4] = GameboardLocation{X: 1, Y: 2}
	initial[5] = GameboardLocation{X: 2, Y: 2}
	pond.init(initial)

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetGameboardValue(GameboardLocation{X: 2, Y: 0}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 3, Y: 1}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 0, Y: 2}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 3, Y: 2}, 0)
	expected.SetGameboardValue(GameboardLocation{X: 1, Y: 3}, 0)

	// Go through one generation
	SimultaneousProcessor(pond, Standard)

	// Compare the pond with the expected version
	if !pond.gameboard.Equals(expected) {
		t.FailNow()
	}
}
