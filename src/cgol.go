package main

import "cgol"
import "fmt"

func main() {
	/*
		TODO: Will need a daemon that intercepts new pond requests and provides data on all current ponds. Use json
	*/
	s := cgol.CreateStrategy("Standard,Orthogonal,Random",
		cgol.CreatePond(5, 20, cgol.NEIGHBORS_ORTHOGONAL),
		func(pond *cgol.Pond) []cgol.OrganismReference { return cgol.InitRandom(pond, 80) },
		cgol.Standard,
		new(cgol.QueueProcessor))

	fmt.Println("===== Starting the thing =====")
	fmt.Print(s)

	s.Start()
}
