package cgol

import "fmt"

func SimultaneousProcessor(pond *Pond, rules func(int, bool) bool) {
	fmt.Printf("SimultaneousProcessor()\n")
	// For each living organism, push to processing channel
	//	calculate num neighbors
	//	if living and over or under pop, push to kill channel and send neighbors to processing channel
	//	if dead and can be revived, then send to revive channel and send neighbors to processing channel

	type ModifiedOrganism struct {
		loc GameboardLocation
		val int
	}
	// modifiedOrganisms := make(chan ModifiedOrganism)
	processingQueue := make(chan GameboardLocation)

	// Add living organisms to processing queue
	for _, row := range pond.living {
		for _, organism := range row {
			// processingQueue <- organism
			fmt.Printf("DEBUG: processingQueue <- organism: %s\n", organism.String())
		}
	}
	numToProcess := len(pond.living)
	fmt.Printf(" numToProcess = %d\n", numToProcess)

	// Process the queue
	for i := 0; i < numToProcess; i++ {
		// Retrieve organism from channel, get its neighbors and see if it is alive
		// organism := <-processingQueue
		organism := GameboardLocation{X: 1, Y: 1}
		fmt.Printf("DEBUG: organism := <-processingQueue\n")
		numNeighbors, neighbors := pond.calculateNeighborCount(organism)
		currentlyAlive := pond.isOrganismAlive(organism)

		// Check with the ruleset what this organism's current status is
		organismStatus := rules(numNeighbors, currentlyAlive)

		if currentlyAlive != organismStatus { // If its status has changed, then we do stuff
			pond.Status = Active

			if organismStatus { // If is alive
				// modifiedOrganisms <- ModifiedOrganism{loc: organism, val: 0}
				fmt.Printf("DEBUG: modifiedOrganisms <- ModifiedOrganism{loc: organism, val: 0}\n")
				for _, neighbor := range neighbors {
					// processingQueue <- neighbor
					fmt.Printf("DEBUG: processingQueue <- neighbor: %s\n", neighbor.String())
					numToProcess++
				}
			}
		}
	}
	close(processingQueue)
	// close(modifiedOrganisms)

	// TODO: loop that handles modifiedOrganisms

	if pond.NumLiving > 0 {
		pond.Status = Stable
	} else {
		pond.Status = Dead
	}
}
