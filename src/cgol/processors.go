package cgol

import "fmt"

func SimultaneousProcessor(pond *Pond, rules func(int, bool) bool) {
	fmt.Printf("SimultaneousProcessor()\n")
	// For each living organism, push to processing channel
	//	calculate num neighbors
	//	if living and over or under pop, push to kill channel and send neighbors to processing channel
	//	if dead and can be revived, then send to revive channel and send neighbors to processing channel

	done := make(chan bool, 1)

	////// Modifications handler /////
	type ModifiedOrganism struct {
		loc GameboardLocation
		val int
	}

	modifications := make(chan ModifiedOrganism, pond.gameboard.Dims.GetCapacity())
	blockModifications := make(chan bool, 1)
	numModifications := 0

	queueModification := make(chan ModifiedOrganism)
	go func() {
		for {
			mod, more := <-queueModification
			if more {
				modifications <- mod
				fmt.Printf(" queued up organism mod: %s, %d\n", mod.loc.String(), mod.val)
				numModifications++
			} else {
				fmt.Printf(" stopped accepting modifications\n")
				close(modifications)
				break
			}
		}
	}()

	go func() {
		fmt.Println("blocking modifications...")
		<-blockModifications

		fmt.Printf("making %d modifications...\n", numModifications)
		for {
			mod, more := <-modifications
			if more {
				fmt.Printf("%d\n", mod)
				pond.setOrganismValue(mod.loc, mod.val)
			} else {
				fmt.Println("no more modifications")
				break
			}
		}

		// Send a value to notify that we're done.
		fmt.Println("done")
		done <- true
	}()

	///// Start processing stuffs /////

	// numToProcess := pond.GetNumLiving() // FIXME Is this correct?
	fmt.Printf(" living = %v\n", pond.living)
	// fmt.Printf(" numToProcess = %d\n", numToProcess)
	processingQueue := make(chan GameboardLocation, pond.gameboard.Dims.GetCapacity()+pond.GetNumLiving()+10)
	// processingQueue := make(chan GameboardLocation)
	doneProcessing := make(chan bool, 1)

	// Process the queue
	go func() {
		for { //i := 0; i < numToProcess; i++ {
			// Retrieve organism from channel, get its neighbors and see if it is alive
			organism, more := <-processingQueue
			if more {
				// numNeighbors, neighbors := pond.calculateNeighborCount(organism)
				numNeighbors, _ := pond.calculateNeighborCount(organism)
				currentlyAlive := pond.isOrganismAlive(organism)
				fmt.Printf("======= processing organism at %s with %d neighbors and alive status of '%t'\n", organism.String(), numNeighbors, currentlyAlive)

				// Check with the ruleset what this organism's current status is
				organismStatus := rules(numNeighbors, currentlyAlive)
				fmt.Printf("   ruleset isalive verdict: %t\n", organismStatus)

				if currentlyAlive != organismStatus { // If its status has changed, then we do stuff
					pond.Status = Active
					fmt.Printf("   organism will be modified\n")

					if organismStatus { // If is alive
						queueModification <- ModifiedOrganism{loc: organism, val: 0}
					} else {
						queueModification <- ModifiedOrganism{loc: organism, val: -1}
					}
				} else {
					fmt.Printf("   nothing to do for organism\n")
				}

				// Now process the neighbors!
				/*
					for _, neighbor := range neighbors {
						// processingQueue <- neighbor
						// numToProcess++
						fmt.Printf("TODO   processingQueue <- neighbor: %s\n", neighbor.String())
					}
				*/
			} else {
				fmt.Printf("   No longer processing organisms\n")
				close(queueModification)
				doneProcessing <- true
				break
			}
		}
	}()

	// Add living organisms to processing queue
	fmt.Printf("processing >%d living< organisms\n", len(pond.living))
	for _, row := range pond.living {
		for _, organism := range row {
			processingQueue <- organism
			fmt.Printf(">> processingQueue <- organism: %s\n", organism.String())

			// Now process the neighbors!
			_, neighbors := pond.calculateNeighborCount(organism)
			for _, neighbor := range neighbors {
				// processingQueue <- neighbor
				// numToProcess++
				fmt.Printf("    > processingQueue <- neighbor: %s\n", neighbor.String())
			}
		}
	}
	close(processingQueue)

	<-doneProcessing

	// Start processing modifications
	blockModifications <- false

	if pond.NumLiving > 0 {
		pond.Status = Stable
	} else {
		pond.Status = Dead
	}

	// Block until all modifications are done
	<-done
}
