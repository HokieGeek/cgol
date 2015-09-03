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
	pond, err := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pond.SetOrganisms(init(size))

	// Go through one generation
	for i := 0; i < len(expected); i++ {
		processor(pond, rules)

		// Compare the pond with the expected version
		if !pond.gameboard.Equals(expected[i]) {
			t.Fatalf("At iteration %d, actual gameboard\n%s\ndoes not match expected\n%s\n", len(expected), pond.gameboard.String(), expected[i].String())
		}
	}
}

func generateBlinkers(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 3, Width: 3}

	// Build the expected gameboard
	expected := make([]*Gameboard, 2)

	// Period 2
	var err error
	expected[0], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for i := 0; i < 3; i++ {
		expected[0].SetValue(GameboardLocation{X: 1, Y: i}, 0)
	}

	// Period 1
	expected[1], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for _, val := range Blinkers(size) {
		expected[1].SetValue(val, 0)
	}

	return size, Blinkers, expected
}

func generateToads(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 2)

	// Period 2
	var err error
	expected[0], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	expected[0].SetValue(GameboardLocation{X: 2, Y: 0}, 0)

	for row := 1; row <= 2; row++ {
		expected[0].SetValue(GameboardLocation{X: 0, Y: row}, 0)
		expected[0].SetValue(GameboardLocation{X: 3, Y: row}, 0)
	}

	expected[0].SetValue(GameboardLocation{X: 1, Y: 3}, 0)

	// Period 1
	expected[1], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for _, val := range Toads(size) {
		expected[1].SetValue(val, 0)
	}

	return size, Toads, expected
}

func generateBeacons(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 2)

	// Period 2
	var err error
	expected[0], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for row := 0; row < 4; row++ {
		adjust := 0
		if row == 2 || row == 3 {
			adjust = 2
		}

		for i := 0; i < 2; i++ {
			expected[0].SetValue(GameboardLocation{X: i + adjust, Y: row}, 0)
		}
	}

	// Period 1
	expected[1], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for _, val := range Beacons(size) {
		expected[1].SetValue(val, 0)
	}

	return size, Beacons, expected
}

func generatePulsar(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 15, Width: 15}

	// Build the expected gameboard
	expected := make([]*Gameboard, 3)

	// Period 2
	var err error
	expected[0], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
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
	expected[1], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
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

	// Period 1
	expected[2], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	for _, val := range Pulsar(size) {
		expected[2].SetValue(val, 0)
	}

	return size, Pulsar, expected
}

func generateGlider(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 5, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 5)

	var err error
	////// Period 2
	expected[0], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[0].SetValue(GameboardLocation{X: 0, Y: 1}, 0)
	expected[0].SetValue(GameboardLocation{X: 2, Y: 1}, 0)
	// Row 2
	expected[0].SetValue(GameboardLocation{X: 1, Y: 2}, 0)
	expected[0].SetValue(GameboardLocation{X: 2, Y: 2}, 0)
	// Row 3
	expected[0].SetValue(GameboardLocation{X: 1, Y: 3}, 0)

	////// Period 3
	expected[1], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[1].SetValue(GameboardLocation{X: 2, Y: 1}, 0)
	// Row 2
	expected[1].SetValue(GameboardLocation{X: 0, Y: 2}, 0)
	expected[1].SetValue(GameboardLocation{X: 2, Y: 2}, 0)
	// Row 3
	expected[1].SetValue(GameboardLocation{X: 1, Y: 3}, 0)
	expected[1].SetValue(GameboardLocation{X: 2, Y: 3}, 0)

	////// Period 4
	expected[2], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[2].SetValue(GameboardLocation{X: 1, Y: 1}, 0)
	// Row 2
	expected[2].SetValue(GameboardLocation{X: 2, Y: 2}, 0)
	expected[2].SetValue(GameboardLocation{X: 3, Y: 2}, 0)
	// Row 3
	expected[2].SetValue(GameboardLocation{X: 1, Y: 3}, 0)
	expected[2].SetValue(GameboardLocation{X: 2, Y: 3}, 0)

	////// Period 5
	// This is the same seed except shifted over
	expected[3], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[3].SetValue(GameboardLocation{X: 2, Y: 1}, 0)
	// Row 2
	expected[3].SetValue(GameboardLocation{X: 3, Y: 2}, 0)
	// Row 3
	for i := 1; i < 4; i++ {
		expected[3].SetValue(GameboardLocation{X: i, Y: 3}, 0)
	}

	////// Period 6
	// -0--       ----       ----       ----      ----      ----
	// --0-       0-0-       --0-       -0--      --0-      ----
	// 000-       -00-       0-0-       --00      ---0      -0-0
	// ----       -0--       -00-       -00-      -000      --00
	// ----       ----       ----       ----      ----      --0-
	expected[4], err = NewGameboard(size)
	if err != nil {
		t.Fatalf("Gameboard of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[4].SetValue(GameboardLocation{X: 1, Y: 2}, 0)
	expected[4].SetValue(GameboardLocation{X: 3, Y: 2}, 0)
	// Row 2
	expected[4].SetValue(GameboardLocation{X: 2, Y: 3}, 0)
	expected[4].SetValue(GameboardLocation{X: 3, Y: 3}, 0)
	// Row 3
	expected[4].SetValue(GameboardLocation{X: 2, Y: 4}, 0)

	return size, Gliders, expected
}

func generateBlock(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = NewGameboard(size)
		if err != nil {
			t.Fatalf("Gameboard of size %s could not be created\n", size.String())
		}
		for row := 0; row < 2; row++ {
			for col := 0; col < 2; col++ {
				expected[period].SetValue(GameboardLocation{X: col, Y: row}, 0)
			}
		}
	}

	return size, Blocks, expected
}

func generateBeehive(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = NewGameboard(size)
		if err != nil {
			t.Fatalf("Gameboard of size %s could not be created\n", size.String())
		}
		for row := 0; row < 3; row++ {
			switch row {
			case 0, 2:
				expected[period].SetValue(GameboardLocation{X: 1, Y: row}, 0)
				expected[period].SetValue(GameboardLocation{X: 2, Y: row}, 0)
			case 1:
				expected[period].SetValue(GameboardLocation{X: 0, Y: row}, 0)
				expected[period].SetValue(GameboardLocation{X: 3, Y: row}, 0)
			}
		}
	}

	return size, Beehive, expected
}

func generateLoaf(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = NewGameboard(size)
		if err != nil {
			t.Fatalf("Gameboard of size %s could not be created\n", size.String())
		}
		// ROW 1
		expected[period].SetValue(GameboardLocation{X: 1, Y: 0}, 0)
		expected[period].SetValue(GameboardLocation{X: 2, Y: 0}, 0)
		// ROW 2
		expected[period].SetValue(GameboardLocation{X: 0, Y: 1}, 0)
		expected[period].SetValue(GameboardLocation{X: 3, Y: 1}, 0)
		// ROW 3
		expected[period].SetValue(GameboardLocation{X: 1, Y: 2}, 0)
		expected[period].SetValue(GameboardLocation{X: 3, Y: 2}, 0)
		// ROW 4
		expected[period].SetValue(GameboardLocation{X: 2, Y: 3}, 0)
	}

	return size, Loaf, expected
}

func generateBoat(t *testing.T) (GameboardDims, func(GameboardDims) []GameboardLocation, []*Gameboard) {
	size := GameboardDims{Height: 4, Width: 4}

	// Build the expected gameboard
	expected := make([]*Gameboard, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = NewGameboard(size)
		if err != nil {
			t.Fatalf("Gameboard of size %s could not be created\n", size.String())
		}
		// ROW 1
		expected[period].SetValue(GameboardLocation{X: 0, Y: 0}, 0)
		expected[period].SetValue(GameboardLocation{X: 1, Y: 0}, 0)
		// ROW 2
		expected[period].SetValue(GameboardLocation{X: 0, Y: 1}, 0)
		expected[period].SetValue(GameboardLocation{X: 2, Y: 1}, 0)
		// ROW 3
		expected[period].SetValue(GameboardLocation{X: 1, Y: 2}, 0)
	}

	return size, Boat, expected
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

func TestProcessorSimultaneousRulesConwayLifeRandom(t *testing.T) {
	// Build the initial pond
	size := GameboardDims{Height: 16, Width: 16}
	initialLocations := Random(size, 80)

	pondInitialSnapshot, err := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondInitialSnapshot.SetOrganisms(initialLocations)

	pondWorker, err := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondWorker.SetOrganisms(initialLocations)

	// Go through one generation
	SimultaneousProcessor(pondWorker, RulesConwayLife)

	// Compare the pond with the expected version
	if pondWorker.gameboard.Equals(pondInitialSnapshot.gameboard) {
		t.Error("Gameboard did not change after one generation of random intialization")
	}
}

func TestProcessorSimultaneousRulesConwayLifeBlinker(t *testing.T) {
	size, init, expected := generateBlinkers(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeToad(t *testing.T) {
	size, init, expected := generateToads(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBeacon(t *testing.T) {
	size, init, expected := generateBeacons(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBlock(t *testing.T) {
	size, init, expected := generateBlock(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBeehive(t *testing.T) {
	size, init, expected := generateBeehive(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeLoaf(t *testing.T) {
	size, init, expected := generateLoaf(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeBoat(t *testing.T) {
	size, init, expected := generateBoat(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifePulsar(t *testing.T) {
	size, init, expected := generatePulsar(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLifeGliders(t *testing.T) {
	size, init, expected := generateGlider(t)
	testProcessorSimultaneousRulesConwayLife(t, size, init, expected)
}

func BenchmarkProcessorSimultaneousRulesConwayLifePulsar(b *testing.B) {
	// Build the initial pond
	size := GameboardDims{Height: 33, Width: 33}
	pond, err := NewPond(size.Height, size.Width, NEIGHBORS_ALL)
	if err != nil {
		b.Fatalf("Unable to create pond: %s\n", err)
	}
	pond.SetOrganisms(Pulsar(size))

	// Ok, do the benchmark now
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SimultaneousProcessor(pond, RulesConwayLife)
	}
}
