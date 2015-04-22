package model

import (
	"github.com/neelance/dom/bind"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/rpcClient"
	"github.com/cryptix/trakting/types"
)

type TrackList struct {
	Scope *bind.Scope

	Users        []*User
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
		return errgo.Notef(err, "tracks.All() failed")
	}

	m.Tracks = make([]*Track, len(tracks))
	for idx, track := range tracks {
		m.Tracks[idx] = &Track{
			Scope: m.Scope,
			Track: &track,
		}
	}

	users, err := m.client.Users.List()
	if err != nil {
		return errgo.Notef(err, "Users.List() failed")
	}

	m.Users = make([]*User, len(users))
	for idx, user := range users {
		m.Users[idx] = &User{
			Scope: m.Scope,
			User:  &user,
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

type User struct {
	Scope *bind.Scope `json:"-"`
	*types.User
}
