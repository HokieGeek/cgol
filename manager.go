package main

// import (
// 	"bytes"
// 	"github.com/hokiegeek/life/core"
// )

/*
type Manager struct {
	analyzers map[string]*life.Analyzer
}

func (t *Manager) stringId(id []byte) string {
	n := bytes.IndexByte(id, 0)
	return string(id[:n])
}

func (t *Manager) Analyzer(id []byte) *life.Analyzer {
	// TODO: validate the input
	return t.analyzers[t.stringId(id)]
}

func (t *Manager) Add(analyzer *life.Analyzer) {
	// TODO: validate the input
	t.analyzers[t.stringId(analyzer.Id)] = analyzer
}

func (t *Manager) Remove(id []byte) {
	// TODO: validate the input
	delete(t.analyzers, t.stringId(id))
}

func NewManager() *Manager {
	m := new(Manager)

	mgr := new(Manager)
	mgr.analyzers = make(map[string]*life.Analyzer, 0)

	return m
}
*/
