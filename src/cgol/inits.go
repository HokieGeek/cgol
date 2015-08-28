package cgol

import (
	"math/rand"
	"time"
)

// TODO: maybe this should take a gameboard and not a pond?
func InitRandom(gameboard *Gameboard, percent int) []GameboardLocation {
	// TODO: keep trying until at least on living organism has been created?
	initialLiving := make([]GameboardLocation, 0)

	for i := 0; i < gameboard.Dims.Height; i++ {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < gameboard.Dims.Width; j++ {
			if rand.Intn(100) > percent {
				initialLiving = append(initialLiving, GameboardLocation{X: i, Y: j})
			}
		}
	}

	return initialLiving
}

func Blinkers(gameboard *Gameboard) []GameboardLocation {
	const LENGTH = 4 // 3 for the line itself and 1 for the spacer

	initialLiving := make([]GameboardLocation, 0)

	// put in as many lengthx1 vertical lines as you can fit
	// Period 1
	// -1-
	// -1-
	// -1-
	// Period 2
	// ---
	// 111
	// ---

	numPerRow := gameboard.Dims.Width / LENGTH
	numPerCol := gameboard.Dims.Height / LENGTH

	// Special case for when the spacer is not needed
	if numPerRow == 0 && gameboard.Dims.Height == 3 {
		numPerRow = 1
	}
	if numPerCol == 0 && gameboard.Dims.Width == 3 {
		numPerCol = 1
	}

	// fmt.Println(pond.gameboard)
	// fmt.Printf("numPerRow = %d\n", numPerRow)
	// fmt.Printf("numPerCol = %d\n", numPerCol)

	currentY := 1
	for row := 0; row < numPerCol; row++ {
		currentY = (row * LENGTH) + 1
		currentX := 1
		for col := 0; col < numPerRow; col++ {
			currentX = (col * LENGTH) + 1
			// fmt.Printf("X: %d, Y: %d\n", currentX, currentY)
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY - 1})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY + 1})
		}
	}

	return initialLiving
}

func Toads(gameboard *Gameboard) []GameboardLocation {
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
