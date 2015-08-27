package cgol

import (
	// "fmt"
	"time"
)

type Processor interface {
	Schedule(organisms []GameboardLocation)
	Process(pond *Pond, rules func(*Pond, GameboardLocation) bool)
	Stop()
}

//////////////////// DEFAULT ////////////////////

type defaultProcessor struct{}

func (t *defaultProcessor) processOrganism(pond *Pond, rules func(*Pond, GameboardLocation) bool, organism GameboardLocation) bool {
	// Apply rules to organism
	if rules(pond, organism) {
		pond.Status = Active
		return true
	} else {
		if pond.NumLiving > 0 {
			pond.Status = Stable
		} else {
			pond.Status = Dead
		}
	}

	return false
}

//////////////////// QUEUE ////////////////////

type QueueProcessor struct {
	defaultProcessor
	// queueStatus *Gameboard // TODO: need to initialize it
	queueObj   []GameboardLocation // TODO: mutex protect
	ticker     *time.Ticker
	updateRate time.Duration
}

func (t *QueueProcessor) isQueued(organism GameboardLocation) bool {
	return false
	// return t.queueStatus.GetGameboardValue(organism) == 1
}

func (t *QueueProcessor) queue(organisms []GameboardLocation) {
	/*
		    for i,organism := organisms {
		        if ! t.isQueued(organism) {
			        t.queue = append(t.queue, organism)
		            t.queueStatus.SetGameboardValue(organism, 1)
		        }
		    }
	*/
	if len(organisms) > 0 {
		t.queueObj = append(t.queueObj, organisms...)
	}
}

func (t *QueueProcessor) dequeue() GameboardLocation {
	if len(t.queueObj) > 0 {
		front := t.queueObj[0]
		if len(t.queueObj) > 1 {
			t.queueObj = append(t.queueObj[:0], t.queueObj[1:]...)
		} else {
			t.queueObj = nil
		}

		// t.queueStatus.SetGameboardValue(front, 0)
		return front
	} else {
		return GameboardLocation{X: 0, Y: 0}
	}
}

func (t *QueueProcessor) processQueue(pond *Pond, rules func(*Pond, GameboardLocation) bool) {
	if len(t.queueObj) > 0 {
		// 1. Consider an organism. Pop it off the front of the queue
		front := t.dequeue()
		// front := t.queue[0]
		// if len(t.queue) > 1 {
		// 	t.queue = append(t.queue[:0], t.queue[1:]...)
		// } else {
		// 	t.queue = nil
		// }

		// 2. Apply rules to organism
		if t.processOrganism(pond, rules, front) {
			// 3. Propogate any changes to neighbors
			t.Schedule(pond.GetNeighbors(front))
		}
	}
}

func (t *QueueProcessor) Schedule(organisms []GameboardLocation) {
	// fmt.Printf("Schedule(%v)\n", organisms)
	// if len(organisms) > 0 {
	//    t.queue = append(t.queue, organisms...)
	// }
	t.queue(organisms)
}

func (t *QueueProcessor) Process(pond *Pond, rules func(*Pond, GameboardLocation) bool) {
	t.processQueue(pond, rules)
	/* TODO
	t.ticker = time.NewTicker(t.updateRate)
	for {
		select {
		case <-t.ticker.C:
			t.processQueue()
			fmt.Println(t) // TODO: remove
		}
	}
	*/
}

/*
func (t *QueueProcessor) Stop() {
    t.ticker.Stop()
}
*/

func NewQueueProcessor() *QueueProcessor {
	p := new(QueueProcessor)
	p.updateRate = time.Millisecond * 250
	return p
}

//////////////////// SIMULTANEOUS ////////////////////
/*
type SimultaneousProcessor struct {
	defaultProcessor
	list []GameboardLocation // TODO: mutex protect
}

func (t *SimultaneousProcessor) Schedule(organisms []GameboardLocation) {
	t.list = append(t.list, organisms...)
}

func (t *SimultaneousProcessor) Process(pond *Pond, rules func(*Pond, GameboardLocation) bool) {
	if len(t.list) > 0 {
		   // go func() {
		   // // TODO: seriously need some mutex protection here, I think...
		   //     for _,organism := t.list {
		   //         go func() {
		   //         }
		// 1. Consider an organism. Pop it off the front of the queue
		// next := t.dequeue()

		// 2. Apply rules to organism
		// 3. Propogate any changes to neighbors
		//if t.processOrganism(pond, rules, next) {
		// 3. Propogate any changes to neighbors
		// TODO: t.Process(pond.GetNeighbors(front))
		//}
	}
}

func NewSimultaneousProcessor() *SimultaneousProcessor {
	p := new(SimultaneousProcessor)
	return p
}
*/
