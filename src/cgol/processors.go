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
	t.queue = append(t.queue, organisms...)
}

func (t *QueueProcessor) Process(pond *Pond, rules func(*Pond, OrganismReference) bool) bool {
	if len(t.queue) > 0 {
		fmt.Println("Processing...")

		// 1. Consider an organism
		front := t.queue[0]
		// TODO: error handling when queue len == 1
		t.queue = append(t.queue[:0], t.queue[1:]...)

		// 2. Apply rules to organism
		modified := rules(pond, front)

		// 3. Propogate any changes to neighbors
		if modified {
			// TODO: does this logic properly handle when an organism dies?

			neighbors := pond.GetNeighbors(front)
			for _, neighbor := range neighbors {
				if pond.isOrganismAlive(neighbor) {
					pond.incrementNeighborCount(neighbor)
				}
			}
			t.schedule(neighbors)
		}

		return true
	} else {
		return false
	}
}
