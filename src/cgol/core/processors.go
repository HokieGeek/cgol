package cgol

import (
	"io/ioutil"
	"log"
	// "os"
)

func SimultaneousProcessor(pond *Pond, rules func(int, bool) bool) {
	// logger := log.New(os.Stderr, "DEBUG: ", log.Ltime)
	logger := log.New(ioutil.Discard, "DEBUG: ", log.Ltime)
	logger.Printf("SimultaneousProcessor()\n")
	// For each living organism, push to processing channel
	//	calculate num neighbors
	//	if living and over or under pop, push to kill channel and send neighbors to processing channel
	//	if dead and can be revived, then send to revive channel and send neighbors to processing channel

	// Blocks the completion of this function
	done := make(chan bool, 1)

	////// Modifications handler /////
	type ModifiedOrganism struct {
		loc GameboardLocation
		val int
	}

	modifications := make(chan ModifiedOrganism, pond.gameboard.Dims.GetCapacity()*2)
	blockModifications := make(chan bool, 1)
	numModifications := 0

	// queueModification := make(chan ModifiedOrganism)
	queueModification := make(chan ModifiedOrganism, pond.gameboard.Dims.GetCapacity()*2)
	go func() {
		for {
			mod, more := <-queueModification
			if more {
				modifications <- mod
				logger.Printf(" queued up organism mod: %s, %d\n", mod.loc.String(), mod.val)
				numModifications++
			} else {
				logger.Printf(" stopped accepting modifications\n")
				close(modifications)
				break
			}
		}
	}()

	go func() {
		logger.Println("blocking modifications...")
		<-blockModifications

		logger.Printf("making %d modifications...\n", numModifications)
		for {
			mod, more := <-modifications
			if more {
				logger.Printf("%d\n", mod)
				pond.setOrganismValue(mod.loc, mod.val)
			} else {
				logger.Println("no more modifications")
				break
			}
		}

		// Send a value to notify that we're done.
		logger.Println("done")
		done <- true
	}()

	///// Start processing stuffs /////

	logger.Printf(" living = %v\n", pond.living)
	processingQueue := make(chan GameboardLocation, pond.gameboard.Dims.GetCapacity()+pond.GetNumLiving()+10)
	// processingQueue := make(chan GameboardLocation)
	doneProcessing := make(chan bool, 1)

	// Process the queue
	go func() {
		for {
			// Retrieve organism from channel, get its neighbors and see if it is alive
			organism, more := <-processingQueue
			// TODO: should not process an organism which has already been processed
			if more {
				// numNeighbors, neighbors := pond.calculateNeighborCount(organism)
				numNeighbors, _ := pond.calculateNeighborCount(organism)
				currentlyAlive := pond.isOrganismAlive(organism)
				logger.Printf("======= processing organism at %s with %d neighbors and alive status of '%t'\n", organism.String(), numNeighbors, currentlyAlive)

				// Check with the ruleset what this organism's current status is
				organismStatus := rules(numNeighbors, currentlyAlive)
				logger.Printf("   ruleset isalive verdict: %t\n", organismStatus)

				if currentlyAlive != organismStatus { // If its status has changed, then we do stuff
					pond.Status = Active
					logger.Printf("   organism will be modified\n")

					if organismStatus { // If is alive
						queueModification <- ModifiedOrganism{loc: organism, val: 0}
					} else {
						queueModification <- ModifiedOrganism{loc: organism, val: -1}
					}
				} else {
					logger.Printf("   nothing to do for organism\n")
				}

				// Now process the neighbors!
				// TODO: make this work. Need to somehow figure out when to close the channel
				// for _, neighbor := range neighbors {
				// 	processingQueue <- neighbor
				// 	logger.Printf("    > processingQueue <- neighbor: %s\n", neighbor.String())
				// }
			} else {
				logger.Printf("   No longer processing organisms\n")
				close(queueModification)
				doneProcessing <- true
				break
			}
		}
	}()

	////// Ensures no duplication occurs /////
	/*
		type ProcessingCandidate struct {
			loc  GameboardLocation
			resp chan bool
		}
		queueForProcessing := make(chan ProcessingCandidate, pond.gameboard.Dims.GetCapacity()*2)
		go func() {
			for {
			}
		}()
	*/

	// Add living organisms to processing queue
	logger.Printf("processing >%d living< organisms\n", len(pond.living))
	for _, row := range pond.living {
		for _, organism := range row {
			processingQueue <- organism
			logger.Printf(">> processingQueue <- organism: %s\n", organism.String())

			// Now process the neighbors!
			_, neighbors := pond.calculateNeighborCount(organism)
			for _, neighbor := range neighbors {
				processingQueue <- neighbor
				logger.Printf("    > processingQueue <- neighbor: %s\n", neighbor.String())
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
