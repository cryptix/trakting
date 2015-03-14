// +build js

package wsclient

import (
	"bufio"
	"net/rpc"

	"github.com/gopherjs/websocket"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/rpcClient"
)

func New(host string) (*rpcClient.Client, error) {
	conn, err := websocket.Dial(host)
	if err != nil {
		return nil, errgo.Notef(err, "Dial failed (url %s)", host)
	}

	l, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, errgo.Notef(err, "bufio.ReadString")
	}

	if l != "OK\n" {
		return nil, errgo.New("websock not ok")
	}

	rpcc := rpc.NewClient(conn)

	tcClient, e := rpcClient.NewTracksClient(rpcc)
	if e != nil {
		return nil, errgo.Notef(err, "NewTracksClient failed")
	}

	tuClient, e := rpcClient.NewUsersClient(rpcc)
	if e != nil {
		return nil, errgo.Notef(err, "NewUsersClient failed")
	}

	return rpcClient.New(tcClient, tuClient), nil
}
