package main

import "cgol"
import "fmt"

func main() {
	s := cgol.NewStrategy("Standard,Orthogonal,Random",
		cgol.NewPond(5, 20, cgol.NEIGHBORS_ORTHOGONAL),
		func(pond *cgol.Pond) []cgol.GameboardLocation { return cgol.InitRandom(pond, 80) },
		cgol.Standard,
		cgol.NewQueueProcessor())

	fmt.Println("===== Starting the thing =====")
	fmt.Print(s)

	s.Start()
}
