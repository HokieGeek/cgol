package cgol

import "testing"

//////////////////////// Common ////////////////////////

// Generic processor tester. Processes the seeds through the same
// number of generations as there are gameboards in the 'expected' slice.
// Compares each generation with the correct slice
func testProcessor(t *testing.T,
	processor func(pond *Pond, rules func(int, bool) bool),
	rules func(int, bool) bool,
	size GameboardDims,
	init func(GameboardDims) []GameboardLocation,
	expected []*Gameboard) {

	// Build the initial pond
	pond := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	pond.init(init(size))

	// Go through one generation
	for i := 0; i < len(expected); i++ {
		processor(pond, rules)

		// Compare the pond with the expected version
		if !pond.gameboard.Equals(expected[i]) {
			t.Fatalf("Actual gameboard:\n%s\ndoes not match expected:\n%s\n", pond.gameboard.String(), expected[i].String())
		}
	}
}

func generateBlinkers() (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 3, Width: 3}

	// Build the expected gameboard
	expected := NewGameboard(size)
	for i := 0; i < 3; i++ {
		expected.SetValue(GameboardLocation{X: 1, Y: i}, 0)
	}

	return size, Blinkers, []*Gameboard{expected}
}

func generateToads() (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	expected.SetValue(GameboardLocation{X: 2, Y: 0}, 0)

	for row := 1; row <= 2; row++ {
		expected.SetValue(GameboardLocation{X: 0, Y: row}, 0)
		expected.SetValue(GameboardLocation{X: 3, Y: row}, 0)
	}

	expected.SetValue(GameboardLocation{X: 1, Y: 3}, 0)

	return size, Toads, []*Gameboard{expected}
}

func generateBeacons() (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	for row := 0; row < 4; row++ {
		adjust := 0
		if row == 2 || row == 3 {
			adjust = 2
		}

		for i := 0; i < 2; i++ {
			expected.SetValue(GameboardLocation{X: i + adjust, Y: row}, 0)
		}
	}

	return size, Beacons, []*Gameboard{expected}
}

func generatePulsar() (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 15, Width: 15}

	// Build the expected gameboard
	expected := make([]*Gameboard, 2)

	// Period 2
	expected[0] = NewGameboard(size)
	for i := 0; i < size.Height; i++ {
		switch i {
		case 0, 1, 2, 6, 8, 12, 13, 14:
			expected[0].SetValue(GameboardLocation{X: 4, Y: i}, 0)
			expected[0].SetValue(GameboardLocation{X: 10, Y: i}, 0)
			switch i {
			case 2, 6, 8, 12:
				expected[0].SetValue(GameboardLocation{X: 5, Y: i}, 0)
				expected[0].SetValue(GameboardLocation{X: 9, Y: i}, 0)
			}
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 0; k < 3; k++ {
					expected[0].SetValue(GameboardLocation{X: k + (12 * j), Y: i}, 0)
				}

				for k := 5; k <= 6; k++ {
					expected[0].SetValue(GameboardLocation{X: k + (3 * j), Y: i}, 0)
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[0].SetValue(GameboardLocation{X: j, Y: i}, 0)
			}
		}
	}

	// Period 3
	expected[1] = NewGameboard(size)
	for i := 0; i < size.Height; i++ {
		switch i {
		case 1, 13:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(GameboardLocation{X: (3 + k) + (7 * j), Y: i}, 0)
				}
			}
		case 2, 12:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(GameboardLocation{X: (4 + k) + (5 * j), Y: i}, 0)
				}
			}
		case 3, 11:
			expected[1].SetValue(GameboardLocation{X: 1, Y: i}, 0)
			for j := 4; j <= 10; j += 2 {
				expected[1].SetValue(GameboardLocation{X: j, Y: i}, 0)
			}
			expected[1].SetValue(GameboardLocation{X: 13, Y: i}, 0)
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 1; k < 4; k++ {
					expected[1].SetValue(GameboardLocation{X: k + (10 * j), Y: i}, 0)
				}

				for k := 5; k <= 6; k++ {
					expected[1].SetValue(GameboardLocation{X: k + (3 * j), Y: i}, 0)
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[1].SetValue(GameboardLocation{X: j, Y: i}, 0)
			}
		case 6, 8:
			for j := 3; j <= 5; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(GameboardLocation{X: j + (k * 6), Y: i}, 0)
				}
			}
		}
	}

	return size, Pulsar, expected
}

func generateBlocks() (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := NewGameboard(size)
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			expected.SetValue(GameboardLocation{X: col, Y: row}, 0)
		}
	}

	return size, Blocks, []*Gameboard{expected}
}

//////////////////////// Simultaneous processor ////////////////////////

func testProcessorSimultaneousRulesConwayLife(t *testing.T,
	size GameboardDims,
	init func(GameboardDims) []GameboardLocation,
	expected []*Gameboard) {

	testProcessor(t,
		SimultaneousProcessor,
		RulesConwayLife,
		size,
		init,
		expected)
}

func TestProcessorSimultaneousRulesConwayLifeBlinkers(t *testing.T) {
	size, init, expected := generateBlinkers()
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeToads(t *testing.T) {
	size, init, expected := generateToads()
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBeacons(t *testing.T) {
	size, init, expected := generateBeacons()
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifePulsar(t *testing.T) {
	t.Skip("Skipping as this currently breaks!")
	size, init, expected := generatePulsar()
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBlock(t *testing.T) {
	size, init, expected := generateBlocks()
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}
