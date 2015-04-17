package controllers

import (
	"github.com/neelance/dom"
	"gopkg.in/errgo.v1"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/model"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/rpcClient"
)

type List struct {
	*views.List
}

func NewList(c *rpcClient.Client) (*List, error) {
	m := model.NewTrackList(c)

	lis := &views.ListListeners{}

	lis.Search = func(c *dom.EventContext) {
		console.Log("search...")
		m.Scope.Digest()
	}

	lis.Play = func(t *model.Track) dom.Listener {
		return func(c *dom.EventContext) {
			console.Log("playing track", t.Name)
			// console.Dir(t)
			m.PlayingTrack = t
			m.Scope.Digest()
		}
	}

	if err := m.Load(); err != nil {
		return nil, errgo.Notef(err, "getTracks failed")
	}

	v := views.NewList(m, lis)

	return &List{v}, nil
}
