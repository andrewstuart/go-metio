package metio

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

//Increment is a type for tracking counts against time
type Increment struct {
	Count uint64
	Time  time.Time
}

//Reader implements the io.Reader interface and exposes methods for learning
//how much has been read.
type Reader struct {
	br       *bufio.Reader
	snaps    []Increment
	initTime time.Time
}

//NewReader returns a Reader given an io.Reader
func NewReader(r io.Reader) *Reader {
	return &Reader{
		br:       bufio.NewReader(r),
		snaps:    make([]Increment, 0, 1000),
		initTime: time.Now(),
	}
}

//Read implements io.Reader
func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.br.Read(p)
	r.snaps = append(r.snaps, Increment{
		Count: uint64(n),
		Time:  time.Now(),
	})
	return n, err
}

var (
	//ErrPriorToInit indicates that a start time provided was prior to the
	//initialization time of the reader. If returned, this error can be safely
	//ignored as the number will be accurate. It is more of a warning that
	//results may not be indicative of actual rates.
	ErrPriorToInit = fmt.Errorf("metio.Reader: time provided prior to reader init time; results may be unexpected")

	//ErrTimeInFuture indicates that a time provided is in the future. If
	//returned, this error can be safely ignored as the number will be accurate.
	//It is more of a warning that results may not be indicative of actual rates.
	ErrTimeInFuture = fmt.Errorf("metio.Reader: time provided is in the future; results may be unexpected")
)

//Since returns a byte count since the time provided (until now). It may also
//return an error if the time requested was before
func (r *Reader) Since(t time.Time) (uint64, error) {
	var ct uint64
	var err error

	if t.Before(r.initTime) {
		err = ErrPriorToInit
	}

	for i := len(r.snaps) - 1; i >= 0; i-- {
		s := r.snaps[i]

		if s.Time.Before(t) {
			return ct, err
		}

		ct += uint64(s.Count)
	}

	return ct, err
}

//Between returns a count for items between a given time
func (r *Reader) Between(from, to time.Time) (ct uint64, err error) {
	//TODO optimize with binary search for start position?
	for i := range r.snaps {
		t := r.snaps[i].Time
		if t.After(from) && t.Before(to) {
			ct += r.snaps[i].Count
		}
	}

	if to.After(time.Now()) {
		err = ErrTimeInFuture
	}
	if from.Before(r.initTime) {
		err = ErrPriorToInit
	}

	return
}
