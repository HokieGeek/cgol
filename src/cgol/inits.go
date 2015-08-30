package cgol

import (
	// "fmt"
	"math/rand"
	"time"
)

/////////////////// RANDOM ///////////////////

func InitRandom(dimensions GameboardDims, percent int) []GameboardLocation {
	initialLiving := make([]GameboardLocation, 0)

	for i := 0; i < dimensions.Height; i++ {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < dimensions.Width; j++ {
			if rand.Intn(100) > percent {
				initialLiving = append(initialLiving, GameboardLocation{X: i, Y: j})
			}
		}
	}

	return initialLiving
}

/////////////////// OSCILLATORS ///////////////////

func getCountsForDimensions(dimensions GameboardDims, width int, height int) (int, int) {
	numPerRow := dimensions.Width / height
	numPerCol := dimensions.Height / width

	// Special case for when the spacer is not needed
	if numPerRow == 0 && dimensions.Height == height-1 {
		numPerRow = 1
	}
	if numPerCol == 0 && dimensions.Width == width-1 {
		numPerCol = 1
	}

	// fmt.Printf(">>>> numPerRow = %d\n", numPerRow)
	// fmt.Printf(">>>> numPerCol = %d\n", numPerCol)
	return numPerRow, numPerCol
}

func getRepeatingPattern(dimensions GameboardDims, height int, width int,
	pattern func(*[]GameboardLocation, int, int)) []GameboardLocation {
	numPerRow, numPerCol := getCountsForDimensions(dimensions, width, height)

	initialLiving := make([]GameboardLocation, 0)
	for row := 0; row < numPerCol; row++ {
		currentY := (row * height)

		for col := 0; col < numPerRow; col++ {
			currentX := (col * width)
			pattern(&initialLiving, currentX, currentY)
		}
	}

	return initialLiving
}

func Blinkers(dimensions GameboardDims) []GameboardLocation {
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
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 1})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 2})
		})

}

func Toads(dimensions GameboardDims) []GameboardLocation {
	// Period 1
	// -000
	// 000-
	// Period 2
	// --0-
	// 0--0
	// 0--0
	// -0--

	const HEIGHT = 5
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			// ROW 1
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 1})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 2, Y: currentY + 1})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 3, Y: currentY + 1})
			// ROW 2
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX, Y: currentY + 2})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 2})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 2, Y: currentY + 2})
		})
}

func Beacons(dimensions GameboardDims) []GameboardLocation {
	// Period 1
	// 00--
	// 0---
	// ---0
	// --00
	// Period 2
	// 00--
	// 00--
	// --00
	// --00

	const HEIGHT = 5
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			// ROW 1
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX, Y: currentY})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY})
			// ROW 2
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX, Y: currentY + 1})
			// ROW 3
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 3, Y: currentY + 2})
			// ROW 4
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 2, Y: currentY + 3})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 3, Y: currentY + 3})
		})
}

/////////////////// STILL ///////////////////

func Blocks(dimensions GameboardDims) []GameboardLocation {
	// ----
	// -00-
	// -00-
	// ----

	const HEIGHT = 5
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			// ROW 1
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX, Y: currentY})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY})
			// ROW 2
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX, Y: currentY + 1})
			*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + 1, Y: currentY + 1})
		})
}

/////////////////// GLIDERS ///////////////////
