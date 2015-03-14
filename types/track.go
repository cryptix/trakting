package types

import (
	"fmt"
	"time"
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

type TrackService interface {
	Add(Track) error
	Get(id string) (Track, error)
	All() ([]Track, error)
	ByUserName(name string) ([]Track, error)
}
