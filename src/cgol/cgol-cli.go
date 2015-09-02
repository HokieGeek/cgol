package main

import (
	"cgol/core"
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
	/*
		fmt.Println("===== Starting random sim =====")
		random := cgol.NewStrategy("RulesConwayLife,Random",
			cgol.NewPond(5, 20, cgol.NEIGHBORS_ALL),
			func(dimensions *GameboardDims) []cgol.GameboardLocation { return cgol.Random(dimensions, 80) },
			cgol.RulesConwayLife,
			cgol.NewQueueProcessor())

		fmt.Print(random)

		random.Start()
	*/

	fmt.Println("===== Starting toads sim =====")
	toads := cgol.NewStrategy("RulesConwayLife,Toads",
		cgol.NewPond(10, 10, cgol.NEIGHBORS_ALL),
		cgol.Toads,
		cgol.RulesConwayLife,
		cgol.SimultaneousProcessor)
	displayPond(toads, 2, false)

	/*
		fmt.Println("===== Starting glider sim =====")
		glider := cgol.NewStrategy("RulesConwayLife,Glider",
			cgol.NewPond(10, 10, cgol.NEIGHBORS_ALL),
			func(dimensions cgol.GameboardDims) []cgol.GameboardLocation {
				return cgol.Gliders(cgol.GameboardDims{Height: 4, Width: 4})
			},
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		glider.UpdateRate = time.Second * 2
		displayPond(glider, 2, false)

		fmt.Println("===== Starting pulsar sim =====")
		pulsar := cgol.NewStrategy("RulesConwayLife,Pulsar",
			// cgol.NewPond(10, 15, cgol.NEIGHBORS_ALL),
			cgol.NewPond(15, 15, cgol.NEIGHBORS_ALL),
			cgol.Pulsar,
			cgol.RulesConwayLife,
			cgol.SimultaneousProcessor)
		displayPond(pulsar, 3, false)
	*/
}
