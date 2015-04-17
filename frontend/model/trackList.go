package model

import (
	"github.com/neelance/dom/bind"

	"github.com/cryptix/trakting/rpcClient"
	"github.com/cryptix/trakting/types"
)

type TrackList struct {
	Scope *bind.Scope

	Tracks       []*Track
	PlayingTrack *Track
	SearchText   string

	client *rpcClient.Client
}

func NewTrackList(c *rpcClient.Client) *TrackList {
	return &TrackList{
		Scope:  bind.NewScope(),
		client: c,
	}
}

func (m *TrackList) Load() error {
	tracks, err := m.client.Tracks.All()
	if err != nil {
		return err
	}

	m.Tracks = make([]*Track, len(tracks))
	for idx, track := range tracks {
		m.Tracks[idx] = &Track{
			Scope: m.Scope,
			Track: &track,
		}
	}

	return nil
}

func (t *TrackList) QueueCount() int {
	return len(t.Tracks)
}

type Track struct {
	Scope *bind.Scope `json:"-"`
	*types.Track
}
