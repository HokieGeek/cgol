package cgol

import "fmt"

type Processor interface {
	schedule(organisms []OrganismReference)
	Process(pond *Pond, rules func(*Pond, OrganismReference) bool) bool
}

//////////////////// QUEUE ////////////////////

type QueueProcessor struct {
	queue []OrganismReference
}

func (t *QueueProcessor) schedule(organisms []OrganismReference) {
	fmt.Printf("schedule(%v)\n", organisms)
	t.queue = append(t.queue, organisms...)
}

func (t *QueueProcessor) Process(pond *Pond, rules func(*Pond, OrganismReference) bool) bool {
	if len(t.queue) > 0 {
		fmt.Printf("\n>>>> Processing (queue size: %d)...\n", len(t.queue))
		fmt.Printf("QUEUE: %v\n", t.queue)

		// 1. Consider an organism
		front := t.queue[0]
		if len(t.queue) > 1 {
			t.queue = append(t.queue[:0], t.queue[1:]...)
		} else {
			t.queue = nil
		}
		fmt.Printf(" Organism: %s\n", front.String())

		// 2. Apply rules to organism
		modified := rules(pond, front)

		// 3. Propogate any changes to neighbors
		if modified {
			neighbors := pond.GetNeighbors(front)
			// fmt.Printf("Found %d neighbors: %v\n", len(neighbors), neighbors)
			/*
				for _, neighbor := range neighbors {
					if pond.isOrganismAlive(neighbor) {
						// fmt.Printf("  Neighbor %s is alive\n", neighbor.String())
						// pond.incrementNeighborCount(neighbor)
						t.schedule([]OrganismReference{neighbor})
					}
				}
			*/
			t.schedule(neighbors)
		}

		return true
	} else {
		return false
	}
}
