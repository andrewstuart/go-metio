package metio

import (
	"sync"
	"time"
)

//Meterer is an abstraction for the ability to report on quantities over time
type Meterer interface {
	Since(time.Time) (uint64, error)
	Between(time.Time, time.Time) (uint64, error)
}

//ReaderGroup is a helper for reporting across multiple concurrent readers
type ReaderGroup struct {
	m       sync.Mutex
	readers map[*Reader]*Reader
}

//NewReaderGroup returns an initialized ReaderGroup.
func NewReaderGroup() *ReaderGroup {
	return &ReaderGroup{
		m:       sync.Mutex{},
		readers: make(map[*Reader]*Reader),
	}
}

//Add adds a reader to the reporting set
func (rg *ReaderGroup) Add(r *Reader) {
	rg.m.Lock()
	defer rg.m.Unlock()

	rg.readers[r] = r
}

//Remove takes a reader from the reporting set
func (rg *ReaderGroup) Remove(r *Reader) {
	delete(rg.readers, r)
}

//Since implements Meterer
func (rg *ReaderGroup) Since(t time.Time) (ct uint64, err error) {
	rg.m.Lock()
	defer rg.m.Unlock()

	for r := range rg.readers {
		c, _ := r.Since(t)
		ct += c
	}

	return
}

//Between implements Meterer
func (rg *ReaderGroup) Between(from, to time.Time) (ct uint64, err error) {
	rg.m.Lock()
	defer rg.m.Unlock()

	var c uint64
	for r := range rg.readers {
		c, err = r.Between(from, to)
		ct += c
	}

	return
}
