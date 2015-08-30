package main

import (
	"cgol"
	"fmt"
	// "time"
)

// TODO: Use variadic arguments to add a stop?
func displayPond(strategy *cgol.Strategy) {
	// Clear the screen and put the cursor on the top left
	fmt.Print("\033[2J")
	fmt.Print("\033[H")

	// Seed
	fmt.Print(strategy)

	updates := make(chan bool)
	strategy.Start(updates)

	/*
		go func() {
			time.Sleep(strategy.UpdateRate * 2)
			strategy.Stop()
		}()
	*/

	for {
		select {
		case <-updates:
			fmt.Print("\033[H")
			fmt.Print(strategy)
		}
	}
}

func main() {
	/*
		fmt.Println("===== Starting random sim =====")
		random := cgol.NewStrategy("Standard,Random",
			cgol.NewPond(5, 20, cgol.NEIGHBORS_ALL),
			func(pond *cgol.Pond) []cgol.GameboardLocation { return cgol.InitRandom(pond, 80) },
			cgol.Standard,
			cgol.NewQueueProcessor())

		fmt.Print(random)

		random.Start()
	*/

	/*
		fmt.Println("===== Starting blinkers sim =====")
		blinkers := cgol.NewStrategy("Standard,Blinkers",
			cgol.NewPond(9, 9, cgol.NEIGHBORS_ALL),
			cgol.Blinkers,
			cgol.Standard,
			cgol.SimultaneousProcessor)
		displayPond(blinkers)
	*/

	fmt.Println("===== Starting toads sim =====")
	toads := cgol.NewStrategy("Standard,Toads",
		// cgol.NewPond(10, 15, cgol.NEIGHBORS_ALL),
		cgol.NewPond(4, 4, cgol.NEIGHBORS_ALL),
		cgol.Toads,
		cgol.Standard,
		cgol.SimultaneousProcessor)
	displayPond(toads)
}
