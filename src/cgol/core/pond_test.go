package cgol

import "testing"

func TestPondSettingInitialPatterns(t *testing.T) {
	rows := 3
	cols := 3
	pond := NewPond(rows, cols, NEIGHBORS_ALL)

	// Create a pattern and call the pond's init function
	initialLiving := make([]GameboardLocation, 3)
	initialLiving[0] = GameboardLocation{X: 0, Y: 0}
	initialLiving[1] = GameboardLocation{X: 1, Y: 1}
	initialLiving[2] = GameboardLocation{X: 2, Y: 2}

	pond.init(initialLiving)

	// Check each expected value
	for _, loc := range initialLiving {
		if !pond.isOrganismAlive(loc) {
			t.Fatalf("Seed organism is not alive!: %s\n", loc.String())
		}
	}
}

func TestPondNeighborSelection(t *testing.T) {
	// TODO
}

func TestPondOrganismValue(t *testing.T) {
	expectedVal := 2
	pos := GameboardLocation{X: 0, Y: 0}
	pond := NewPond(1, 1, NEIGHBORS_ALL)
	pond.setOrganismValue(pos, expectedVal)

	actualVal := pond.GetOrganismValue(pos)

	if actualVal != expectedVal {
		t.Fatalf("Retrieved actual value %d instead of expected value %d\n", actualVal, expectedVal)
	}
}

func TestPondGetNumLiving(t *testing.T) {
	// TODO
}

func TestPondNeighborCountCalutation(t *testing.T) {
	// TODO
}
