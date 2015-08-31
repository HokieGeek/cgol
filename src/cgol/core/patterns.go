package cgol

import (
	// "fmt"
	"math/rand"
	"time"
)

/////////////////////////// COMMON ///////////////////////////

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

func Blinkers(dimensions GameboardDims) []GameboardLocation {
	// put in as many lengthx1 vertical lines as you can fit
	// Period 1   Period 2
	// -0-		  ---
	// -0-        000
	// -0-        ---

	const HEIGHT = 4 // 3 for the line itself and 1 for the spacer
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			for i := 0; i < 3; i++ {
				*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + i, Y: currentY + 1})
			}
		})

}

func Toads(dimensions GameboardDims) []GameboardLocation {
	// Period 1	  Period 2
	// ----       --0-
	// -000       0--0
	// 000-       0--0
	// ----       -0--

	const HEIGHT = 5
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			// ROW 1
			for i := 1; i < 4; i++ {
				*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + i, Y: currentY + 1})
			}
			// ROW 2
			for i := 0; i < 3; i++ {
				*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + i, Y: currentY + 2})
			}
		})
}

func Beacons(dimensions GameboardDims) []GameboardLocation {
	// Period 1   Period 2
	// 00--       00--
	// 0---       00--
	// ---0       --00
	// --00       --00

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

func Pulsar(dimensions GameboardDims) []GameboardLocation {
	// Period 1          Period 2           Period 3
	//                                        012345678901234
	// ---------------   ----0-----0----    0 ---------------
	// ---000---000---   ----0-----0----    1 ---00-----00---
	// ---------------   ----00---00----    2 ----00---00----
	// -0----0-0----0-   ---------------    3 -0--0-0-0-0--0-
	// -0----0-0----0-   000--00-00--000    4 -000-00-00-000-
	// -0----0-0----0-   --0-0-0-0-0-0--    5 --0-0-0-0-0-0--
	// ---000---000---   ----00---00----    6 ---000---000---
	// ---------------   ---------------    7 ---------------
	// ---000---000---   ----00---00----    8 ---000---000---
	// -0----0-0----0-   --0-0-0-0-0-0--    9 --0-0-0-0-0-0--
	// -0----0-0----0-   000--00-00--000    0 -000-00-00-000-
	// -0----0-0----0-   ---------------    1 -0--0-0-0-0--0-
	// ---------------   ----00---00----    2 ----00---00----
	// ---000---000---   ----0-----0----    3 ---00-----00---
	// ---------------   ----0-----0----    4 ---------------

	const HEIGHT = 16
	return getRepeatingPattern(dimensions, HEIGHT, HEIGHT,
		func(initialLiving *[]GameboardLocation, currentX int, currentY int) {
			for i := 0; i < HEIGHT-1; i++ {
				switch i {
				case 1, 6, 8, 13:
					// Repeating pattern 1: --000---000--, rows: 0,5,7,12
					for j := 0; j < 2; j++ {
						for k := 3; k <= 5; k++ {
							*initialLiving = append(*initialLiving,
								GameboardLocation{X: currentX + (k + (6 * j)), Y: currentY + i})
						}
					}
				case 3, 4, 5, 9, 10, 11:
					// Repeating pattern 2: 0----0-0----0, rows: 2,3,4,8,9,10
					for j := 0; j < 2; j++ {
						*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + (1 + (7 * j)), Y: currentY + i})
						*initialLiving = append(*initialLiving, GameboardLocation{X: currentX + (6 + (7 * j)), Y: currentY + i})
					}
				}
			}
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
