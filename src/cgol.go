package main

import "cgol"

func main() {
	/*
		TODO: Will need a daemon that intercepts new pond requests and provides data on all current ponds. Use json
	*/
	s := cgol.CreateStrategy("Standard,Orthogonal,Random",
		cgol.CreatePond(5, 20, cgol.NEIGHBORS_ORTHOGONAL),
		func(pond *cgol.Pond) []cgol.OrganismReference { return cgol.InitRandom(pond, 80) },
		cgol.Standard,
		new(cgol.QueueProcessor))

	s.Display()

	// s.Start()
	// s.Display()
}
