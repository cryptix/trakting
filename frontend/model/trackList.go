package model

import (
	"github.com/neelance/dom/bind"

	"github.com/cryptix/trakting/types"
)

type TrackList struct {
	Scope *bind.Scope

	Tracks       []*Track
	PlayingTrack *Track
	SearchText   string
}

func (t *TrackList) QueueCount() int {
	return len(t.Tracks)
}

type Track struct {
	Scope *bind.Scope `json:"-"`
	*types.Track
}
