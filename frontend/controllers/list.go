package controllers

import (
	"github.com/neelance/dom"
	"gopkg.in/errgo.v1"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/model"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/rpcClient"
)

func NewList(c *rpcClient.Client) (*views.List, error) {
	m := model.NewTrackList(c)

	lis := &views.ListListeners{}

	lis.Search = func(c *dom.EventContext) {
		console.Log("search..." + m.SearchText)
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

	lis.Reload = func(c *dom.EventContext) {
		console.Log("reload clicked")
		if err := m.Load(); err != nil {
			console.Error(err)
		}
		m.Scope.Digest()
	}

	if err := m.Load(); err != nil {
		return nil, errgo.Notef(err, "getTracks failed")
	}

	return views.NewList(m, lis), nil
}
