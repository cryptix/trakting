package rpcClient

import (
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/types"
)

var ErrTODO = errgo.New("TODO")

type Client struct {
	Tracks types.Tracker
	Users  types.Userer
}
