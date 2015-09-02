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

		blinkers := cgol.NewStrategy("RulesConwayLife,Blinkers",
			cgol.NewPond(height, width, cgol.NEIGHBORS_ALL),
			cgol.Blinkers,
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		if *ratePtr > 0 {
			blinkers.UpdateRate = *ratePtr
		}
		displayPond(blinkers, -1, true)
	case "toads":
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 10
		if *heightPtr > height {
			height = *heightPtr
		}

		toads := cgol.NewStrategy("RulesConwayLife,Toads",
			cgol.NewPond(height, width, cgol.NEIGHBORS_ALL),
			cgol.Toads,
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		if *ratePtr > 0 {
			toads.UpdateRate = *ratePtr
		}
		displayPond(toads, -1, true)
	case "glider":
		width := 30
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 30
		if *heightPtr > height {
			height = *heightPtr
		}

		glider := cgol.NewStrategy("RulesConwayLife,Glider",
			cgol.NewPond(height, width, cgol.NEIGHBORS_ALL),
			func(dimensions cgol.GameboardDims) []cgol.GameboardLocation {
				return cgol.Gliders(cgol.GameboardDims{Height: 4, Width: 4})
			},
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		if *ratePtr > 0 {
			glider.UpdateRate = *ratePtr
		}
		displayPond(glider, -1, true)
	case "pulsar":
		width := 15
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 15
		if *heightPtr > height {
			height = *heightPtr
		}

		pulsar := cgol.NewStrategy("RulesConwayLife,Pulsar",
			cgol.NewPond(height, width, cgol.NEIGHBORS_ALL),
			cgol.Pulsar,
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		if *ratePtr > 0 {
			pulsar.UpdateRate = *ratePtr
		}
		displayPond(pulsar, -1, true)
	default:
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 20
		if *heightPtr > height {
			height = *heightPtr
		}

		random := cgol.NewStrategy("RulesConwayLife,Random",
			cgol.NewPond(height, width, cgol.NEIGHBORS_ALL),
			func(dimensions cgol.GameboardDims) []cgol.GameboardLocation { return cgol.Random(dimensions, 80) },
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		if *ratePtr > 0 {
			random.UpdateRate = *ratePtr
		}
		displayPond(random, -1, true)
	}
}
