package metio

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type Snapshot struct {
	Count int
	Time  time.Time
}

type Reader struct {
	br       *bufio.Reader
	snaps    []Snapshot
	initTime time.Time
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		br:       bufio.NewReader(r),
		snaps:    make([]Snapshot, 0, 1000),
		initTime: time.Now(),
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.br.Read(p)
	r.snaps = append(r.snaps, Snapshot{
		Count: n,
		Time:  time.Now(),
	})
	return n, err
}

var PriorToInit = fmt.Errorf("`since` time prior to reader init time")

//Since returns a byte count since the time provided (until now). It may also
//return an error if the time requested was before
func (r *Reader) Since(t time.Time) (int, error) {
	ct := 0
	var err error

	if t.Before(r.initTime) {
		err = PriorToInit
	}

	for i := len(r.snaps) - 1; i >= 0; i-- {
		s := r.snaps[i]

		if s.Time.Before(t) {
			return ct, err
		} else {
			ct += s.Count
		}
	}

	return ct, err
}
