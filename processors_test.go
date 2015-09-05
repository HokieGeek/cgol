package life

import "testing"

//////////////////////// Common ////////////////////////

// Generic processor tester. Processes the seeds through the same
// number of generations as there are boards in the 'expected' slice.
// Compares each generation with the correct slice
func testProcessor(t *testing.T,
	processor func(pond *pond, rules func(int, bool) bool),
	rules func(int, bool) bool,
	size Dimensions,
	init func(Dimensions, Location) []Location,
	expected []*board) {

	// Build the initial pond
	pond, err := newpond(size, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pond.SetOrganisms(init(size, Location{}))

	// Go through one generation
	for i := 0; i < len(expected); i++ {
		processor(pond, rules)

		// Compare the pond with the expected version
		if !pond.board.Equals(expected[i]) {
			t.Fatalf("At iteration %d, actual board\n%s\ndoes not match expected\n%s\n", len(expected), pond.board.String(), expected[i].String())
		}
	}
}

func generateBlinkers(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 3, Width: 3}

	// Build the expected board
	expected := make([]*board, 2)

	// Period 2
	var err error
	expected[0], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for i := 0; i < 3; i++ {
		expected[0].SetValue(Location{X: 1, Y: i}, 0)
	}

	// Period 1
	expected[1], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for _, val := range Blinkers(size, Location{}) {
		expected[1].SetValue(val, 0)
	}

	return size, Blinkers, expected
}

func generateToads(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 2)

	// Period 2
	var err error
	expected[0], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	expected[0].SetValue(Location{X: 2, Y: 0}, 0)

	for row := 1; row <= 2; row++ {
		expected[0].SetValue(Location{X: 0, Y: row}, 0)
		expected[0].SetValue(Location{X: 3, Y: row}, 0)
	}

	expected[0].SetValue(Location{X: 1, Y: 3}, 0)

	// Period 1
	expected[1], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for _, val := range Toads(size, Location{}) {
		expected[1].SetValue(val, 0)
	}

	return size, Toads, expected
}

func generateBeacons(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 2)

	// Period 2
	var err error
	expected[0], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for row := 0; row < 4; row++ {
		adjust := 0
		if row == 2 || row == 3 {
			adjust = 2
		}

		for i := 0; i < 2; i++ {
			expected[0].SetValue(Location{X: i + adjust, Y: row}, 0)
		}
	}

	// Period 1
	expected[1], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for _, val := range Beacons(size, Location{}) {
		expected[1].SetValue(val, 0)
	}

	return size, Beacons, expected
}

func generatePulsar(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 15, Width: 15}

	// Build the expected board
	expected := make([]*board, 3)

	// Period 2
	var err error
	expected[0], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for i := 0; i < size.Height; i++ {
		switch i {
		case 0, 1, 2, 6, 8, 12, 13, 14:
			expected[0].SetValue(Location{X: 4, Y: i}, 0)
			expected[0].SetValue(Location{X: 10, Y: i}, 0)
			switch i {
			case 2, 6, 8, 12:
				expected[0].SetValue(Location{X: 5, Y: i}, 0)
				expected[0].SetValue(Location{X: 9, Y: i}, 0)
			}
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 0; k < 3; k++ {
					expected[0].SetValue(Location{X: k + (12 * j), Y: i}, 0)
				}

				for k := 5; k <= 6; k++ {
					expected[0].SetValue(Location{X: k + (3 * j), Y: i}, 0)
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[0].SetValue(Location{X: j, Y: i}, 0)
			}
		}
	}

	// Period 3
	expected[1], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for i := 0; i < size.Height; i++ {
		switch i {
		case 1, 13:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(Location{X: (3 + k) + (7 * j), Y: i}, 0)
				}
			}
		case 2, 12:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(Location{X: (4 + k) + (5 * j), Y: i}, 0)
				}
			}
		case 3, 11:
			expected[1].SetValue(Location{X: 1, Y: i}, 0)
			for j := 4; j <= 10; j += 2 {
				expected[1].SetValue(Location{X: j, Y: i}, 0)
			}
			expected[1].SetValue(Location{X: 13, Y: i}, 0)
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 1; k < 4; k++ {
					expected[1].SetValue(Location{X: k + (10 * j), Y: i}, 0)
				}

				for k := 5; k <= 6; k++ {
					expected[1].SetValue(Location{X: k + (3 * j), Y: i}, 0)
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[1].SetValue(Location{X: j, Y: i}, 0)
			}
		case 6, 8:
			for j := 3; j <= 5; j++ {
				for k := 0; k < 2; k++ {
					expected[1].SetValue(Location{X: j + (k * 6), Y: i}, 0)
				}
			}
		}
	}

	// Period 1
	expected[2], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	for _, val := range Pulsar(size, Location{}) {
		expected[2].SetValue(val, 0)
	}

	return size, Pulsar, expected
}

func generateGlider(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 5, Width: 4}

	// Build the expected board
	expected := make([]*board, 5)

	var err error
	////// Period 2
	expected[0], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[0].SetValue(Location{X: 0, Y: 1}, 0)
	expected[0].SetValue(Location{X: 2, Y: 1}, 0)
	// Row 2
	expected[0].SetValue(Location{X: 1, Y: 2}, 0)
	expected[0].SetValue(Location{X: 2, Y: 2}, 0)
	// Row 3
	expected[0].SetValue(Location{X: 1, Y: 3}, 0)

	////// Period 3
	expected[1], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[1].SetValue(Location{X: 2, Y: 1}, 0)
	// Row 2
	expected[1].SetValue(Location{X: 0, Y: 2}, 0)
	expected[1].SetValue(Location{X: 2, Y: 2}, 0)
	// Row 3
	expected[1].SetValue(Location{X: 1, Y: 3}, 0)
	expected[1].SetValue(Location{X: 2, Y: 3}, 0)

	////// Period 4
	expected[2], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[2].SetValue(Location{X: 1, Y: 1}, 0)
	// Row 2
	expected[2].SetValue(Location{X: 2, Y: 2}, 0)
	expected[2].SetValue(Location{X: 3, Y: 2}, 0)
	// Row 3
	expected[2].SetValue(Location{X: 1, Y: 3}, 0)
	expected[2].SetValue(Location{X: 2, Y: 3}, 0)

	////// Period 5
	// This is the same seed except shifted over
	expected[3], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[3].SetValue(Location{X: 2, Y: 1}, 0)
	// Row 2
	expected[3].SetValue(Location{X: 3, Y: 2}, 0)
	// Row 3
	for i := 1; i < 4; i++ {
		expected[3].SetValue(Location{X: i, Y: 3}, 0)
	}

	////// Period 6
	// -0--       ----       ----       ----      ----      ----
	// --0-       0-0-       --0-       -0--      --0-      ----
	// 000-       -00-       0-0-       --00      ---0      -0-0
	// ----       -0--       -00-       -00-      -000      --00
	// ----       ----       ----       ----      ----      --0-
	expected[4], err = newBoard(size)
	if err != nil {
		t.Fatalf("board of size %s could not be created\n", size.String())
	}
	// Row 1
	expected[4].SetValue(Location{X: 1, Y: 2}, 0)
	expected[4].SetValue(Location{X: 3, Y: 2}, 0)
	// Row 2
	expected[4].SetValue(Location{X: 2, Y: 3}, 0)
	expected[4].SetValue(Location{X: 3, Y: 3}, 0)
	// Row 3
	expected[4].SetValue(Location{X: 2, Y: 4}, 0)

	return size, Gliders, expected
}

func generateBlock(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = newBoard(size)
		if err != nil {
			t.Fatalf("board of size %s could not be created\n", size.String())
		}
		for row := 0; row < 2; row++ {
			for col := 0; col < 2; col++ {
				expected[period].SetValue(Location{X: col, Y: row}, 0)
			}
		}
	}

	return size, Blocks, expected
}

func generateBeehive(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = newBoard(size)
		if err != nil {
			t.Fatalf("board of size %s could not be created\n", size.String())
		}
		for row := 0; row < 3; row++ {
			switch row {
			case 0, 2:
				expected[period].SetValue(Location{X: 1, Y: row}, 0)
				expected[period].SetValue(Location{X: 2, Y: row}, 0)
			case 1:
				expected[period].SetValue(Location{X: 0, Y: row}, 0)
				expected[period].SetValue(Location{X: 3, Y: row}, 0)
			}
		}
	}

	return size, Beehive, expected
}

func generateLoaf(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = newBoard(size)
		if err != nil {
			t.Fatalf("board of size %s could not be created\n", size.String())
		}
		// ROW 1
		expected[period].SetValue(Location{X: 1, Y: 0}, 0)
		expected[period].SetValue(Location{X: 2, Y: 0}, 0)
		// ROW 2
		expected[period].SetValue(Location{X: 0, Y: 1}, 0)
		expected[period].SetValue(Location{X: 3, Y: 1}, 0)
		// ROW 3
		expected[period].SetValue(Location{X: 1, Y: 2}, 0)
		expected[period].SetValue(Location{X: 3, Y: 2}, 0)
		// ROW 4
		expected[period].SetValue(Location{X: 2, Y: 3}, 0)
	}

	return size, Loaf, expected
}

func generateBoat(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*board) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*board, 4)

	var err error
	for period := 0; period < len(expected); period++ {
		expected[period], err = newBoard(size)
		if err != nil {
			t.Fatalf("board of size %s could not be created\n", size.String())
		}
		// ROW 1
		expected[period].SetValue(Location{X: 0, Y: 0}, 0)
		expected[period].SetValue(Location{X: 1, Y: 0}, 0)
		// ROW 2
		expected[period].SetValue(Location{X: 0, Y: 1}, 0)
		expected[period].SetValue(Location{X: 2, Y: 1}, 0)
		// ROW 3
		expected[period].SetValue(Location{X: 1, Y: 2}, 0)
	}

	return size, Boat, expected
}

//////////////////////// Simultaneous processor ////////////////////////

func testProcessorSimultaneousRulesConway(t *testing.T,
	size Dimensions,
	init func(Dimensions, Location) []Location,
	expected []*board) {

	testProcessor(t,
		SimultaneousProcessor,
		ConwayTester(),
		size,
		init,
		expected)
}

func TestProcessorSimultaneousRulesConwayRandom(t *testing.T) {
	// Build the initial pond
	size := Dimensions{Height: 16, Width: 16}
	initialLocations := Random(size, Location{}, 80)

	pondInitialSnapshot, err := newpond(size, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondInitialSnapshot.SetOrganisms(initialLocations)

	pondWorker, err := newpond(size, NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondWorker.SetOrganisms(initialLocations)

	// Go through one generation
	SimultaneousProcessor(pondWorker, ConwayTester())

	// Compare the pond with the expected version
	if pondWorker.board.Equals(pondInitialSnapshot.board) {
		t.Error("board did not change after one generation of random intialization")
	}
}

func TestProcessorSimultaneousRulesConwayBlinker(t *testing.T) {
	size, init, expected := generateBlinkers(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayToad(t *testing.T) {
	size, init, expected := generateToads(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayBeacon(t *testing.T) {
	size, init, expected := generateBeacons(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayBlock(t *testing.T) {
	size, init, expected := generateBlock(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayBeehive(t *testing.T) {
	size, init, expected := generateBeehive(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayLoaf(t *testing.T) {
	size, init, expected := generateLoaf(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayBoat(t *testing.T) {
	size, init, expected := generateBoat(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayPulsar(t *testing.T) {
	size, init, expected := generatePulsar(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func TestProcessorSimultaneousRulesConwayGliders(t *testing.T) {
	size, init, expected := generateGlider(t)
	testProcessorSimultaneousRulesConway(t, size, init, expected)
}

func BenchmarkProcessorSimultaneousRulesConwayPulsar(b *testing.B) {
	// Build the initial pond
	size := Dimensions{Height: 33, Width: 33}
	pond, err := newpond(size, NEIGHBORS_ALL)
	if err != nil {
		b.Fatalf("Unable to create pond: %s\n", err)
	}
	pond.SetOrganisms(Pulsar(size, Location{}))

	// Ok, do the benchmark now
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SimultaneousProcessor(pond, ConwayTester())
	}
}
