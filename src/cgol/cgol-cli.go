package main

import (
	"cgol/core"
	"flag"
	"fmt"
	"time"
)

func displayPond(strategy *cgol.Strategy, iterations time.Duration, static bool) {
	// Clear the screen and put the cursor on the top left
	if static {
		fmt.Print("\033[2J")
		fmt.Print("\033[H")
	}

	// Print the seed
	fmt.Print(strategy)

	updates := make(chan bool)
	strategy.Start(updates)

	if iterations > -1 {
		go func() {
			time.Sleep(strategy.UpdateRate * iterations)
			strategy.Stop()
		}()
	}

	for {
		select {
		case <-updates:
			if static {
				fmt.Print("\033[H")
			}
			fmt.Print(strategy)
		}
	}
}

func displayTestPond(width int, height int, rate time.Duration, initializer func(cgol.LifeboardDims) []cgol.LifeboardLocation) {
	pond, err := cgol.NewPond(height, width, cgol.NEIGHBORS_ALL)
	if err == nil {
		strategy := cgol.NewStrategy("Test",
			pond,
			initializer,
			cgol.GetConwayTester(),
			cgol.SimultaneousProcessor)
		if rate > 0 {
			strategy.UpdateRate = rate
		}
		displayPond(strategy, -1, true)
	} else {
		fmt.Printf("Could not create: %s\n", err)
	}
}

func main() {

	patternPtr := flag.String("pattern", "blinkers", "Specify the pattern to run")
	widthPtr := flag.Int("width", 1, "Width of the Life board")
	heightPtr := flag.Int("height", 1, "Height of the Life board")
	ratePtr := flag.Duration("rate", -1, "Rate at which the board should be updated in milliseconds")

	flag.Parse()

	switch *patternPtr {
	case "blinkers":
		width := 9
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 9
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestPond(width, height, *ratePtr, cgol.Blinkers)
	case "toads":
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 10
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestPond(width, height, *ratePtr, cgol.Toads)
	case "glider":
		width := 30
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 30
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestPond(width, height, *ratePtr,
			func(dimensions cgol.LifeboardDims) []cgol.LifeboardLocation {
				return cgol.Gliders(cgol.LifeboardDims{Height: 4, Width: 4})
			})

	case "pulsar":
		width := 15
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 15
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestPond(width, height, *ratePtr, cgol.Pulsar)
	default:
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 20
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestPond(width, height, *ratePtr,
			func(dimensions cgol.LifeboardDims) []cgol.LifeboardLocation {
				return cgol.Random(dimensions, 80)
			})
	}
}
