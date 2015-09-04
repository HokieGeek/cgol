package main

import (
	"flag"
	"fmt"
	"github.com/hokiegeek/life"
	"time"
)

func displaypond(strategy *life.Life, iterations time.Duration, static bool) {
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

func displayTestpond(width int, height int, rate time.Duration, initializer func(life.LifeboardDims) []life.LifeboardLocation) {
	strategy, err := life.New("Test",
		life.LifeboardDims{Height: height, Width: width},
		life.NEIGHBORS_ALL,
		initializer,
		life.GetConwayTester(),
		life.SimultaneousProcessor)
	if err == nil {
		if rate > 0 {
			strategy.UpdateRate = rate
		}
		displaypond(strategy, -1, true)
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

		displayTestpond(width, height, *ratePtr, life.Blinkers)
	case "toads":
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 10
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr, life.Toads)
	case "glider":
		width := 30
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 30
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr,
			func(dimensions life.LifeboardDims) []life.LifeboardLocation {
				return life.Gliders(life.LifeboardDims{Height: 4, Width: 4})
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

		displayTestpond(width, height, *ratePtr, life.Pulsar)
	default:
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 20
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr,
			func(dimensions life.LifeboardDims) []life.LifeboardLocation {
				return life.Random(dimensions, 80)
			})
	}
}
