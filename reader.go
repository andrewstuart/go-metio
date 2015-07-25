package metio

import (
	"bufio"
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
	emitters map[chan int]int //Map a channel to its index
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		br:       bufio.NewReader(r),
		snaps:    make([]Snapshot, 0, 1000),
		initTime: time.Now(),
		emitters: make(map[chan int]int),
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

//Since returns a byte count since the time provided (until now).
func (r *Reader) Since(t time.Time) (int, error) {
	ct := 0
	var err error

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

func (r *Reader) Emitter() chan int {
	ch := make(chan int)
	r.emitters[ch] = 0

	go func() {
		for {

		}
	}()

	return nil
}

func sum(vals []int) int {
	sum := 0
	for j := range vals {
		sum += vals[j]
	}
	return sum
}
