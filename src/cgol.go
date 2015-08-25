package main

import "cgol"

func main() {
	/*
		Will need a daemon that intercepts new pond requests and provides data on all current ponds. Use json
	*/
	// cgol.CreatePond("Standard,Random", 5, 20, cgol.StandardRules, cgol.RandomInit)
	tmp := cgol.CreatePond("Standard,Random", 5, 20,
		cgol.StandardOrthogonal,
		func(pond *cgol.Pond) { cgol.InitRandom(pond, 80) })
	tmp.Display()
}
