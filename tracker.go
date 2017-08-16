package life

type trackerAddOp struct {
	loc  Location
	resp chan bool
}

type trackerRemoveOp struct {
	loc  Location
	resp chan bool
}

type trackerTestOp struct {
	loc  Location
	resp chan bool
}

type trackerGetAllOp struct {
	resp chan []Location
}

type trackerCountOp struct {
	resp chan int
}

type tracker struct {
	trackerAdd    chan *trackerAddOp
	trackerRemove chan *trackerRemoveOp
	trackerTest   chan *trackerTestOp
	trackerGetAll chan *trackerGetAllOp
	trackerCount  chan *trackerCountOp
}

func (t *tracker) living() {
	var livingMap = make(map[int]map[int]Location)
	var count int
	// logger := log.New(os.Stderr, "tracker: ", log.Ltime)
	// logger := log.New(ioutil.Discard, "tracker: ", log.Ltime)

	for {
		select {
		case add := <-t.trackerAdd:
			added := true
			if _, keyExists := livingMap[add.loc.Y]; !keyExists {
				livingMap[add.loc.Y] = make(map[int]Location)
			}
			if _, keyExists := livingMap[add.loc.Y][add.loc.X]; !keyExists {
				livingMap[add.loc.Y][add.loc.X] = add.loc
				count++
			}
			add.resp <- added
		case remove := <-t.trackerRemove:
			removed := false
			_, keyExists := livingMap[remove.loc.Y]
			if keyExists {
				_, keyExists = livingMap[remove.loc.Y][remove.loc.X]
				if keyExists {
					delete(livingMap[remove.loc.Y], remove.loc.X)
					removed = true
					count--

					// TODO Delete the row if it has no children?
					// if len(livingMap[remove.loc.Y]) <= 0 {
					// 	delete(livingMap, remove.loc.Y)
					// }
				}
			}
			remove.resp <- removed
		case test := <-t.trackerTest:
			_, keyExists := livingMap[test.loc.Y]
			if keyExists {
				_, keyExists = livingMap[test.loc.Y][test.loc.X]
			}
			test.resp <- keyExists
		case getall := <-t.trackerGetAll:
			all := make([]Location, 0)
			for rowNum := range livingMap {
				for _, col := range livingMap[rowNum] {
					all = append(all, col)
				}
			}
			getall.resp <- all
		case countOp := <-t.trackerCount:
			countOp.resp <- count
		}
	}
}

func (t *tracker) Set(location Location) bool {
	add := &trackerAddOp{loc: location, resp: make(chan bool)}
	t.trackerAdd <- add
	val := <-add.resp

	return val
}

func (t *tracker) Remove(location Location) bool {
	remove := &trackerRemoveOp{loc: location, resp: make(chan bool)}
	t.trackerRemove <- remove
	val := <-remove.resp

	return val
}

func (t *tracker) Test(location Location) bool {
	read := &trackerTestOp{loc: location, resp: make(chan bool)}
	t.trackerTest <- read
	val := <-read.resp

	return val
}

func (t *tracker) GetAll() []Location {
	get := &trackerGetAllOp{resp: make(chan []Location)}
	t.trackerGetAll <- get
	val := <-get.resp

	return val
}

func (t *tracker) GetCount() int {
	count := &trackerCountOp{resp: make(chan int)}
	t.trackerCount <- count
	val := <-count.resp

	return val
}

func newTracker() *tracker {
	t := new(tracker)

	t.trackerAdd = make(chan *trackerAddOp)
	t.trackerRemove = make(chan *trackerRemoveOp)
	t.trackerTest = make(chan *trackerTestOp)
	t.trackerGetAll = make(chan *trackerGetAllOp)
	t.trackerCount = make(chan *trackerCountOp)

	go t.living()

	return t
}

// vim: set foldmethod=marker:
