package rpcClient

import (
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/types"
)

var ErrTODO = errgo.New("TODO")

type Client struct{}

var _ types.TrackService = (*Client)(nil)

func (c *Client) Add(t types.Track) error {
	return ErrTODO
}

func (c *Client) All() ([]types.Track, error) {
	return nil, ErrTODO
}

func (c *Client) ByUserName(u string) ([]types.Track, error) {
	return nil, ErrTODO
}

func (c *Client) Get(boomid string) (types.Track, error) {
	return types.Track{}, ErrTODO
}
