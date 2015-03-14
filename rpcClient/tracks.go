package rpcClient

import (
	"net/rpc"

	"github.com/cryptix/trakting/types"
)

type tracks struct {
	client *rpc.Client
}

var _ types.Tracker = (*tracks)(nil)

func NewTracksClient(c *rpc.Client) (types.Tracker, error) {
	return &tracks{
		client: c,
	}, nil
}

func (c *tracks) Add(t types.Track) error {
	// TODO: add proper validation
	if t.Name == "" {
		return types.ErrEmptyTrackName
	}
	return c.client.Call("TrackService.Add", t, nil)
}

func (c *tracks) All() ([]types.Track, error) {
	var tracks []types.Track
	return tracks, c.client.Call("TrackService.All", "", &tracks)
}

func (c *tracks) ByUserName(u string) ([]types.Track, error) {
	var tracks []types.Track
	return tracks, c.client.Call("TrackService.ByUserName", u, &tracks)
}

func (c *tracks) Get(boomid string) (types.Track, error) {
	var t types.Track
	return t, c.client.Call("TrackService.Get", boomid, &t)
}
