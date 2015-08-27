package cgol

import (
	"math/rand"
	"time"
)

// TODO: maybe this should take a gameboard and not a pond?
func InitRandom(pond *Pond, percent int) []GameboardLocation {
	// TODO: keep trying until at least on living organism has been created?
	initialLiving := make([]GameboardLocation, 0)

	for i := 0; i < pond.gameboard.Dims.Height; i++ {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < pond.gameboard.Dims.Width; j++ {
			if rand.Intn(100) > percent {
				initialLiving = append(initialLiving, GameboardLocation{X: i, Y: j})
			}
		}
	}

	return initialLiving
}

func Blinkers(pond *Pond) []GameboardLocation {
	// length := 3

	initialLiving := make([]GameboardLocation, 0)

	// TODO: put in as many lengthx1 vertical lines as you can fit
	// Period 1
	// -1-
	// -1-
	// -1-
	// Period 2
	// ---
	// 111
	// ---
	//
	// 1. Determine how many can fit vertically (rows / length+1) [+1 == spacer]
	// 2. Determine how many can fit horizontally (cols / length+1)
	// 3. Determine each viable centerPoint
	//    Really just need the very first one, after that I need to do
	//    centerPoint.X + width for each col and
	//    centerPoint.Y + height for each row
	// 4. For each centerPoint
	//    1. initialiLiving = append(initialLiving, centerPoint)
	//    2. initialiLiving = append(initialLiving, centerPoint.Y-1)
	//    3. initialiLiving = append(initialLiving, centerPoint.Y+1)

	return initialLiving
}

func Toads(pond *Pond) []GameboardLocation {
	initialLiving := make([]GameboardLocation, 0)

	// TODO
	// Period 1
	// -111
	// 111-
	// Period 2
	// --1-
	// 1--1
	// 1--1
	// -1--

	return initialLiving
}
