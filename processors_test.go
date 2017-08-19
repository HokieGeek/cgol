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
	expected []*tracker) {

	// Build the initial pond
	pond, err := newPond(size, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pond.SetOrganisms(init(size, Location{}))

	// Go through one generation
	for i := 0; i < len(expected); i++ {
		processor(pond, rules)

		// Compare the pond with the expected version
		if !pond.living.Equals(expected[i]) {
			t.Fatalf("At iteration %d, actual board\n%s\ndoes not match expected\n", len(expected), pond.String())
		}
	}
}

func generateBlinkers(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 3, Width: 3}

	// Build the expected board
	expected := make([]*tracker, 2)

	// Period 2
	expected[0] = newTracker()
	for i := 0; i < 3; i++ {
		expected[0].Set(Location{X: 1, Y: i})
	}

	// Period 1
	expected[1] = newTracker()
	for _, val := range Blinkers(size, Location{}) {
		expected[1].Set(val)
	}

	return size, Blinkers, expected
}

func generateToads(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 2)

	// Period 2
	expected[0] = newTracker()
	expected[0].Set(Location{X: 2, Y: 0})

	for row := 1; row <= 2; row++ {
		expected[0].Set(Location{X: 0, Y: row})
		expected[0].Set(Location{X: 3, Y: row})
	}

	expected[0].Set(Location{X: 1, Y: 3})

	// Period 1
	expected[1] = newTracker()
	for _, val := range Toads(size, Location{}) {
		expected[1].Set(val)
	}

	return size, Toads, expected
}

func generateBeacons(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 2)

	// Period 2
	expected[0] = newTracker()
	for row := 0; row < 4; row++ {
		adjust := 0
		if row == 2 || row == 3 {
			adjust = 2
		}

		for i := 0; i < 2; i++ {
			expected[0].Set(Location{X: i + adjust, Y: row})
		}
	}

	// Period 1
	expected[1] = newTracker()
	for _, val := range Beacons(size, Location{}) {
		expected[1].Set(val)
	}

	return size, Beacons, expected
}

func generatePulsar(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 15, Width: 15}

	// Build the expected board
	expected := make([]*tracker, 3)

	// Period 2
	expected[0] = newTracker()
	for i := 0; i < size.Height; i++ {
		switch i {
		case 0, 1, 2, 6, 8, 12, 13, 14:
			expected[0].Set(Location{X: 4, Y: i})
			expected[0].Set(Location{X: 10, Y: i})
			switch i {
			case 2, 6, 8, 12:
				expected[0].Set(Location{X: 5, Y: i})
				expected[0].Set(Location{X: 9, Y: i})
			}
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 0; k < 3; k++ {
					expected[0].Set(Location{X: k + (12 * j), Y: i})
				}

				for k := 5; k <= 6; k++ {
					expected[0].Set(Location{X: k + (3 * j), Y: i})
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[0].Set(Location{X: j, Y: i})
			}
		}
	}

	// Period 3
	expected[1] = newTracker()
	for i := 0; i < size.Height; i++ {
		switch i {
		case 1, 13:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].Set(Location{X: (3 + k) + (7 * j), Y: i})
				}
			}
		case 2, 12:
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					expected[1].Set(Location{X: (4 + k) + (5 * j), Y: i})
				}
			}
		case 3, 11:
			expected[1].Set(Location{X: 1, Y: i})
			for j := 4; j <= 10; j += 2 {
				expected[1].Set(Location{X: j, Y: i})
			}
			expected[1].Set(Location{X: 13, Y: i})
		case 4, 10:
			for j := 0; j < 2; j++ {
				for k := 1; k < 4; k++ {
					expected[1].Set(Location{X: k + (10 * j), Y: i})
				}

				for k := 5; k <= 6; k++ {
					expected[1].Set(Location{X: k + (3 * j), Y: i})
				}
			}
		case 5, 9:
			for j := 2; j <= 12; j += 2 {
				expected[1].Set(Location{X: j, Y: i})
			}
		case 6, 8:
			for j := 3; j <= 5; j++ {
				for k := 0; k < 2; k++ {
					expected[1].Set(Location{X: j + (k * 6), Y: i})
				}
			}
		}
	}

	// Period 1
	expected[2] = newTracker()
	for _, val := range Pulsar(size, Location{}) {
		expected[2].Set(val)
	}

	return size, Pulsar, expected
}

func generateGlider(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 5, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 5)

	////// Period 2
	expected[0] = newTracker()
	// Row 1
	expected[0].Set(Location{X: 0, Y: 1})
	expected[0].Set(Location{X: 2, Y: 1})
	// Row 2
	expected[0].Set(Location{X: 1, Y: 2})
	expected[0].Set(Location{X: 2, Y: 2})
	// Row 3
	expected[0].Set(Location{X: 1, Y: 3})

	////// Period 3
	expected[1] = newTracker()
	// Row 1
	expected[1].Set(Location{X: 2, Y: 1})
	// Row 2
	expected[1].Set(Location{X: 0, Y: 2})
	expected[1].Set(Location{X: 2, Y: 2})
	// Row 3
	expected[1].Set(Location{X: 1, Y: 3})
	expected[1].Set(Location{X: 2, Y: 3})

	////// Period 4
	expected[2] = newTracker()
	// Row 1
	expected[2].Set(Location{X: 1, Y: 1})
	// Row 2
	expected[2].Set(Location{X: 2, Y: 2})
	expected[2].Set(Location{X: 3, Y: 2})
	// Row 3
	expected[2].Set(Location{X: 1, Y: 3})
	expected[2].Set(Location{X: 2, Y: 3})

	////// Period 5
	// This is the same seed except shifted over
	expected[3] = newTracker()
	// Row 1
	expected[3].Set(Location{X: 2, Y: 1})
	// Row 2
	expected[3].Set(Location{X: 3, Y: 2})
	// Row 3
	for i := 1; i < 4; i++ {
		expected[3].Set(Location{X: i, Y: 3})
	}

	////// Period 6
	// -0--       ----       ----       ----      ----      ----
	// --0-       0-0-       --0-       -0--      --0-      ----
	// 000-       -00-       0-0-       --00      ---0      -0-0
	// ----       -0--       -00-       -00-      -000      --00
	// ----       ----       ----       ----      ----      --0-
	expected[4] = newTracker()
	// Row 1
	expected[4].Set(Location{X: 1, Y: 2})
	expected[4].Set(Location{X: 3, Y: 2})
	// Row 2
	expected[4].Set(Location{X: 2, Y: 3})
	expected[4].Set(Location{X: 3, Y: 3})
	// Row 3
	expected[4].Set(Location{X: 2, Y: 4})

	return size, Gliders, expected
}

func generateBlock(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 4)

	for period := 0; period < len(expected); period++ {
		expected[period] = newTracker()
		for row := 0; row < 2; row++ {
			for col := 0; col < 2; col++ {
				expected[period].Set(Location{X: col, Y: row})
			}
		}
	}

	return size, Blocks, expected
}

func generateBeehive(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 4)

	for period := 0; period < len(expected); period++ {
		expected[period] = newTracker()
		for row := 0; row < 3; row++ {
			switch row {
			case 0, 2:
				expected[period].Set(Location{X: 1, Y: row})
				expected[period].Set(Location{X: 2, Y: row})
			case 1:
				expected[period].Set(Location{X: 0, Y: row})
				expected[period].Set(Location{X: 3, Y: row})
			}
		}
	}

	return size, Beehive, expected
}

func generateLoaf(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 4)

	for period := 0; period < len(expected); period++ {
		expected[period] = newTracker()
		// ROW 1
		expected[period].Set(Location{X: 1, Y: 0})
		expected[period].Set(Location{X: 2, Y: 0})
		// ROW 2
		expected[period].Set(Location{X: 0, Y: 1})
		expected[period].Set(Location{X: 3, Y: 1})
		// ROW 3
		expected[period].Set(Location{X: 1, Y: 2})
		expected[period].Set(Location{X: 3, Y: 2})
		// ROW 4
		expected[period].Set(Location{X: 2, Y: 3})
	}

	return size, Loaf, expected
}

func generateBoat(t *testing.T) (Dimensions, func(Dimensions, Location) []Location, []*tracker) {
	size := Dimensions{Height: 4, Width: 4}

	// Build the expected board
	expected := make([]*tracker, 4)

	for period := 0; period < len(expected); period++ {
		expected[period] = newTracker()
		// ROW 1
		expected[period].Set(Location{X: 0, Y: 0})
		expected[period].Set(Location{X: 1, Y: 0})
		// ROW 2
		expected[period].Set(Location{X: 0, Y: 1})
		expected[period].Set(Location{X: 2, Y: 1})
		// ROW 3
		expected[period].Set(Location{X: 1, Y: 2})
	}

	return size, Boat, expected
}

//////////////////////// Simultaneous processor ////////////////////////

func testProcessorSimultaneousRulesConway(t *testing.T,
	size Dimensions,
	init func(Dimensions, Location) []Location,
	expected []*tracker) {

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

	pondInitialSnapshot, err := newPond(size, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondInitialSnapshot.SetOrganisms(initialLocations)

	pondWorker, err := newPond(size, newTracker(), NEIGHBORS_ALL)
	if err != nil {
		t.Fatalf("Unable to create pond: %s\n", err)
	}
	pondWorker.SetOrganisms(initialLocations)

	// Go through one generation
	SimultaneousProcessor(pondWorker, ConwayTester())

	// Compare the pond with the expected version
	if pondWorker.Equals(pondInitialSnapshot) {
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
	pond, err := newPond(size, newTracker(), NEIGHBORS_ALL)
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

// vim: set foldmethod=marker:
