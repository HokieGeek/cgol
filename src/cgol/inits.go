package cgol

import (
	// "fmt"
	"math/rand"
	"time"
)

func RandomInit(pond *Pond) {
	pond.gameboard = make([][]int, pond.Rows)
	// pond.NumLiving = pond.Rows * pond.Cols

	// completion := make(chan int, pond.Rows)
	for i := 0; i < pond.Rows; i++ {
		// go func(row int, c chan int) {
		// 	fmt.Printf("Doing: %d\n", row)
		// 	c <- row
		// }(i, completion)
		// go func(c chan int) {
		// fmt.Printf("Row: %d\n", i)
		rand.Seed(time.Now().UnixNano())
		pond.gameboard[i] = make([]int, pond.Cols)
		for j := 0; j < pond.Cols; j++ {
			cell := OrganismReference{X: i, Y: j}
			if rand.Intn(1000) > 970 {
				// pond.updateOrganismLivingState(cell, true)
				pond.gameboard[i][j] = 0
				pond.NumLiving++
				pond.living = append(pond.living, cell)
			} else {
				pond.gameboard[i][j] = -1
				// pond.updateOrganismLivingState(cell, false)
			}
		}
		// c <- i
		// }(completion)
	}
	// for c := range completion {
	// 	fmt.Printf("%d is done\n", c)
	// }
}
