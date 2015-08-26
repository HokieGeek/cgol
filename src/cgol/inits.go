package cgol

import (
	"math/rand"
	"time"
)

// TODO: maybe this should take a gameboard and not a pond?
func InitRandom(pond *Pond, percent int) []GameboardLocation {
	// TODO: keep trying until at least on living organism has been created?
	initialLiving := make([]GameboardLocation, 0)

	for i := 0; i < pond.gameboard.Dims.Height; i++ {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < pond.gameboard.Dims.Width; j++ {
			if rand.Intn(100) > percent {
				initialLiving = append(initialLiving, GameboardLocation{X: i, Y: j})
			}
		}
	}

	return initialLiving
}

func Blocks(pond *Pond, size int) []GameboardLocation {
	initialLiving := make([]GameboardLocation, 0)

	// TODO: put in as many sizexsize blocks as you can fit

	return initialLiving
}
