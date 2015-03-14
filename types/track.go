package types

import (
	"fmt"
	"time"

	"gopkg.in/errgo.v1"
)

var (
	ErrEmptyTrackName = errgo.New("empty name")
)

type Track struct {
	ID     int64
	By     string
	Name   string
	BoomID string
	Added  time.Time
}

func (t Track) String() string {
	return fmt.Sprintf("%q (by %s)", t.Name, t.By)
}

type Tracker interface {
	Add(Track) error
	Get(id string) (Track, error)
	All() ([]Track, error)
	ByUserName(name string) ([]Track, error)
}
