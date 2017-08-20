package life

func SimultaneousProcessor(pond *pond, rules func(int, bool) bool) {
	// Blocks the completion of this function
	done := make(chan bool, 1)

	////// Modifications handler /////
	type ModifiedOrganism struct {
		loc   Location
		alive bool
	}

	modifications := make(chan ModifiedOrganism, pond.Dims.Capacity())
	blockModifications := make(chan bool, 1)

	// This routine will make the actual modifications to the pond
	go func() {
		<-blockModifications

		for {
			if mod, more := <-modifications; more {
				pond.setOrganismState(mod.loc, mod.alive)
			} else {
				break
			}
		}

		// Send a value to notify that we're done.
		done <- true
	}()

	///// Start processing stuffs /////
	// Process the queue
	processingQueue := make(chan Location, pond.Dims.Capacity())
	go func() {
		processed := make(map[int]map[int]int)
		for {
			// Retrieve organism from channel, get its neighbors and see if it is alive
			if organism, more := <-processingQueue; more {
				// Should not process an organism which has already been processed
				unprocessed := true
				if _, keyExists := processed[organism.Y]; keyExists {
					_, keyExists = processed[organism.Y][organism.X]
					unprocessed = !keyExists
				} else {
					processed[organism.Y] = make(map[int]int)
				}

				// Since this is a new one, go ahead and process it
				if unprocessed {
					// Add organism to list of processed
					processed[organism.Y][organism.X] = 1

					// Retrieve all the infos
					neighbors, err := pond.GetNeighbors(organism)
					if err == nil {
						numLivingNeighbors := 0
						for _, neighbor := range neighbors {
							if pond.isOrganismAlive(neighbor) {
								numLivingNeighbors++
							}
						}
						currentlyAlive := pond.isOrganismAlive(organism)

						// Check with the ruleset what this organism's current status is
						organismStatus := rules(numLivingNeighbors, currentlyAlive)

						if currentlyAlive != organismStatus { // If its status has changed, then we do stuff
							modifications <- ModifiedOrganism{loc: organism, alive: organismStatus}
						}
					}
				}
			} else {
				close(modifications)
				blockModifications <- false
				break
			}
		}
	}()

	// Add living organisms to processing queue
	for _, organism := range pond.living.GetAll() {
		processingQueue <- organism

		// Now process the neighbors!
		neighbors, err := pond.GetNeighbors(organism)
		if err == nil {
			for _, neighbor := range neighbors {
				processingQueue <- neighbor
			}
		}
	}
	close(processingQueue)

	// Block until all modifications are done
	<-done
}

// vim: set foldmethod=marker:
