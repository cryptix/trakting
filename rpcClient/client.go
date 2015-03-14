package rpcClient

import "github.com/cryptix/trakting/types"

type Client struct {
	Tracks types.Tracker
	Users  types.Userer
}

func New(t types.Tracker, u types.Userer) *Client {
	return &Client{
		Tracks: t,
		Users:  u,
	}
}
