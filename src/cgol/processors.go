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
			mod := <-queueModification
			modifications <- mod
			fmt.Printf(" queued up organism mod: %s, %d\n", mod.loc.String(), mod.val)
			numModifications++
		}
	}()

	go func() {
		fmt.Println("blocking modifications...")
		<-blockModifications

		fmt.Printf("making %d modifications...\n", numModifications)
		for i := 0; i < numModifications; i++ {
			select {
			case mod := <-modifications:
				fmt.Printf("%d ", mod)
				pond.setOrganismValue(mod.loc, mod.val)
			}
		}

		// Send a value to notify that we're done.
		fmt.Println("done")
		done <- true
	}()
	/*

		// Fake out some modifications
		for n := range make([]int, 10) {
			queueModification <- n
		}

		// Wait a few seconds before unblocking modifications
		time.Sleep(time.Second * 2)
		blockModifications <- false
	*/

	///// Start processing stuffs /////

	// numToProcess := len(pond.living) // FIXME Is this correct?
	fmt.Printf(" living = %v\n", pond.living)
	// fmt.Printf(" numToProcess = %d\n", numToProcess)
	processingQueue := make(chan GameboardLocation)
	// doneProcessing := make(chan bool, 1)

	// Process the queue
	go func() {
		for { // i := 0; i < numToProcess; i++ {
			// Retrieve organism from channel, get its neighbors and see if it is alive
			organism := <-processingQueue
			// organism := GameboardLocation{X: 1, Y: 1}
			numNeighbors, neighbors := pond.calculateNeighborCount(organism)
			currentlyAlive := pond.isOrganismAlive(organism)
			fmt.Printf("   processing organism at %s with %d neighbors and alive status of '%t'\n", organism.String(), numNeighbors, currentlyAlive)

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
				for _, neighbor := range neighbors {
					// processingQueue <- neighbor
					fmt.Printf("TODO   processingQueue <- neighbor: %s\n", neighbor.String())
					// numToProcess++
				}
			} else {
				fmt.Printf("   nothing to do for organism\n")
			}
		}
		// doneProcessing <- true
	}()

	// Add living organisms to processing queue
	fmt.Printf("processing >%d living< organisms\n", len(pond.living))
	for _, row := range pond.living {
		for _, organism := range row {
			processingQueue <- organism
			fmt.Printf("processingQueue <- organism: %s\n", organism.String())
		}
	}

	// <-doneProcessing

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
