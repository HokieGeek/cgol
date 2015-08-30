package main

import "cgol"
import "fmt"
import "time"

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

	fmt.Println("===== Starting blinkers sim =====")
	blinkers := cgol.NewStrategy("Standard,Blinkers",
		cgol.NewPond(9, 9, cgol.NEIGHBORS_ALL),
		cgol.Blinkers,
		cgol.Standard,
		cgol.SimultaneousProcessor)

	fmt.Println(blinkers)
	blinkers.Start()
	time.Sleep(blinkers.UpdateRate * 2)
	blinkers.Stop()
	fmt.Println(blinkers)

	fmt.Println("===== Starting toads sim =====")
	toads := cgol.NewStrategy("Standard,Toads",
		// cgol.NewPond(10, 15, cgol.NEIGHBORS_ALL),
		cgol.NewPond(4, 4, cgol.NEIGHBORS_ALL),
		cgol.Toads,
		cgol.Standard,
		cgol.SimultaneousProcessor)

	fmt.Println(toads)
	toads.Start()
	time.Sleep(toads.UpdateRate * 1)
	toads.Stop()
	fmt.Println(toads)
}
