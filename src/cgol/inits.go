package cgol

import (
	"math/rand"
	"time"
)

func InitRandom(pond *Pond, percent int) []OrganismReference {
	initialLiving := make([]OrganismReference, 10)

	for i := 0; i < pond.Rows; i++ {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < pond.Cols; j++ {
			if rand.Intn(100) > percent {
				initialLiving = append(initialLiving, OrganismReference{X: i, Y: j})
			}
		}
	}

	return initialLiving
}
