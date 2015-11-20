package metio

import (
	"testing"
	"time"
)

func tp(t string) time.Time {
	tt, _ := time.Parse(time.Kitchen, t)
	return tt
}

func getTestReader() *Reader {
	return &Reader{
		snaps:    []Increment{{100, tp("2:30PM")}, {100, tp("3:00PM")}},
		initTime: tp("1:00PM"),
	}
}

func TestSince(t *testing.T) {
	r := getTestReader()

	s, err := r.Since(tp("2:00PM"))
	if err != nil {
		t.Errorf("Error since: %v", err)
	}

	if s != uint64(200) {
		t.Errorf("wrong number returned: %d should be %d", s, 200)
	}

	s2, err := r.Since(tp("12:30PM"))
	if err != ErrPriorToInit {
		t.Errorf("Error was not equal to ErrPriorToInit: %v", err)
	}
	if s2 != uint64(200) {
		t.Errorf("Wrong number returned: %d should be %d", s2, 200)
	}
}

func TestBetween(t *testing.T) {
	r := getTestReader()

	ct, err := r.Between(tp("2:00PM"), tp("3:00PM"))

	if ct != uint64(100) {
		t.Errorf("Wrong number returned: %d, should be %d", ct, 100)
	}
	if err != nil {
		t.Errorf("Error returned for Between: %v", err)
	}
}
