package life

/*
import (
// "io/ioutil"
"log"
"os"
)
*/

func SimultaneousProcessor(pond *pond, rules func(int, bool) bool) {
	// logger := log.New(os.Stderr, "SimultaneousProcessor: ", log.Ltime)
	// logger := log.New(ioutil.Discard, "SimultaneousProcessor: ", log.Ltime)

	// Blocks the completion of this function
	done := make(chan bool, 1)

	////// Modifications handler /////
	type ModifiedOrganism struct {
		loc Location
		val int
	}

	modifications := make(chan ModifiedOrganism, pond.board.Dims.GetCapacity())
	blockModifications := make(chan bool, 1)

	// This routine will make the actual modifications to the pond
	go func() {
		<-blockModifications

		for {
			mod, more := <-modifications
			if more {
				pond.setOrganismValue(mod.loc, mod.val)
			} else {
				break
			}
		}

		// Send a value to notify that we're done.
		done <- true
	}()

	///// Start processing stuffs /////
	processingQueue := make(chan Location, pond.board.Dims.GetCapacity()+pond.GetNumLiving()+10)
	doneProcessing := make(chan bool, 1)

	// Process the queue
	go func() {
		processed := make(map[int]map[int]int)
		for {
			// Retrieve organism from channel, get its neighbors and see if it is alive
			organism, more := <-processingQueue
			if more {

				// Should not process an organism which has already been processed
				unprocessed := true
				_, keyExists := processed[organism.Y]
				if keyExists {
					_, keyExists = processed[organism.Y][organism.X]
					if keyExists {
						unprocessed = false
					}
				} else {
					processed[organism.Y] = make(map[int]int)
				}

				// Since this is a new one, go ahead and process it
				if unprocessed {
					// Add organism to list of processed
					processed[organism.Y][organism.X] = 1

					// Retrieve all the infos
					// numNeighbors, neighbors := pond.calculateNeighborCount(organism)
					numNeighbors, _ := pond.calculateNeighborCount(organism)
					currentlyAlive := pond.isOrganismAlive(organism)

					// Check with the ruleset what this organism's current status is
					organismStatus := rules(numNeighbors, currentlyAlive)

					if currentlyAlive != organismStatus { // If its status has changed, then we do stuff
						pond.Status = Active

						if organismStatus { // If is alive
							modifications <- ModifiedOrganism{loc: organism, val: 0}
						} else {
							modifications <- ModifiedOrganism{loc: organism, val: -1}
						}
					}

				}
			} else {
				close(modifications)
				doneProcessing <- true
				break
			}
		}
	}()

	// Add living organisms to processing queue
	for _, organism := range pond.living.GetAll() {
		processingQueue <- organism

		// Now process the neighbors!
		_, neighbors := pond.calculateNeighborCount(organism)
		for _, neighbor := range neighbors {
			processingQueue <- neighbor
		}
	}
	close(processingQueue)

	<-doneProcessing

	// Start processing modifications
	blockModifications <- false

	if pond.GetNumLiving() > 0 {
		pond.Status = Stable
	} else {
		pond.Status = Dead
	}

	// Block until all modifications are done
	<-done
}
