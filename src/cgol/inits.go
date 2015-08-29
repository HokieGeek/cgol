package cgol

import (
	"fmt"
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
	// put in as many lengthx1 vertical lines as you can fit
	// Period 1
	// -0-
	// -0-
	// -0-
	// Period 2
	// ---
	// 000
	// ---

	const HEIGHT = 4 // 3 for the line itself and 1 for the spacer

	numPerRow := gameboard.Dims.Width / HEIGHT
	numPerCol := gameboard.Dims.Height / HEIGHT

	// Special case for when the spacer is not needed
	if numPerRow == 0 && gameboard.Dims.Height == HEIGHT-1 {
		numPerRow = 1
	}
	if numPerCol == 0 && gameboard.Dims.Width == HEIGHT-1 {
		numPerCol = 1
	}

	// fmt.Println(pond.gameboard)
	// fmt.Printf("numPerRow = %d\n", numPerRow)
	// fmt.Printf("numPerCol = %d\n", numPerCol)

	initialLiving := make([]GameboardLocation, 0)
	currentY := 1
	for row := 0; row < numPerCol; row++ {
		currentY = (row * HEIGHT) + 1
		currentX := 1
		for col := 0; col < numPerRow; col++ {
			currentX = (col * HEIGHT) + 1
			// fmt.Printf("X: %d, Y: %d\n", currentX, currentY)
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY - 1})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY + 1})
		}
	}

	return initialLiving
}

func Toads(gameboard *Gameboard) []GameboardLocation {

	// TODO
	// Period 1
	// -000
	// 000-
	// Period 2
	// --0-
	// 0--0
	// 0--0
	// -0--

	const (
		HEIGHT = 5
	)

	numPerRow := gameboard.Dims.Width / HEIGHT
	numPerCol := gameboard.Dims.Height / HEIGHT

	// Special case for when the spacer is not needed
	if numPerRow == 0 && gameboard.Dims.Height == HEIGHT-1 {
		numPerRow = 1
	}
	if numPerCol == 0 && gameboard.Dims.Width == HEIGHT-1 {
		numPerCol = 1
	}
	fmt.Printf(">>>> numPerRow = %d\n", numPerRow)
	fmt.Printf(">>>> numPerCol = %d\n", numPerCol)

	initialLiving := make([]GameboardLocation, 0)

	for row := 0; row < numPerCol; row++ {
		currentY := (row * HEIGHT)

		for col := 0; col < numPerRow; col++ {
			currentX := (col * HEIGHT)
			// fmt.Printf("X: %d, Y: %d\n", currentX, currentY)
			// ROW 1
			initialLiving = append(initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 1})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX + 2, Y: currentY + 1})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX + 3, Y: currentY + 1})
			// ROW 2
			initialLiving = append(initialLiving, GameboardLocation{X: currentX, Y: currentY + 2})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 2})
			initialLiving = append(initialLiving, GameboardLocation{X: currentX + 2, Y: currentY + 2})
		}

	}

	return initialLiving
}
