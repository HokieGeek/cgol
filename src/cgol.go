package main

import "cgol"
import "fmt"

func main() {
	/*
		random := cgol.NewStrategy("Standard,Random",
			cgol.NewPond(5, 20, cgol.NEIGHBORS_ALL),
			func(pond *cgol.Pond) []cgol.GameboardLocation { return cgol.InitRandom(pond, 80) },
			cgol.Standard,
			cgol.NewQueueProcessor())

		fmt.Println("===== Starting random sim =====")
		fmt.Print(random)

		random.Start()
	*/

	blinkers := cgol.NewStrategy("Standard,Blinkers",
		cgol.NewPond(3, 3, cgol.NEIGHBORS_ALL),
		cgol.Blinkers,
		cgol.Standard,
		cgol.SimultaneousProcessor)

	fmt.Println("===== Starting blinkers sim =====")
	fmt.Println(blinkers)

	blinkers.Start()
	fmt.Println(blinkers)
}
