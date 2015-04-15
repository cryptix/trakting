package main

import (
	"github.com/cryptix/trakting/rpcClient"
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/model"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/frontend/wsclient"
)

func main() {
	wc, err := wsclient.New("ws://localhost:3000/wsrpc") // inject correct url somehow
	if err != nil {
		panic(err)
	}
	console.Log("rpc connected")

	m := &model.TrackList{
		Scope: bind.NewScope(),
	}

	l := createListeners(m)

	getTracks(wc, m)

	dom.SetTitle("Trakting • Landing")
	// dom.AddStylesheet("css/tt.css")
	dom.SetBody(views.Page(m, l))

}

func createListeners(m *model.TrackList) *views.PageListeners {
	l := &views.PageListeners{}

	l.Search = func(c *dom.EventContext) {
		console.Log("search...")
		m.Scope.Digest()
	}
	l.TogglePlay = func(t *model.Track) dom.Listener {
		return func(c *dom.EventContext) {
			console.Log("toggle track", t.Name)
			console.Dir(t)
			m.Scope.Digest()
		}
	}
	l.QueueAll = func(c *dom.EventContext) {
		console.Log("queue all tracks")
		m.Scope.Digest()
	}

	return l
}

func getTracks(wc *rpcClient.Client, m *model.TrackList) {
	tracks, err := wc.Tracks.All()
	if err != nil {
		panic(err)
	}

	for _, track := range tracks {
		m.Tracks = append(m.Tracks, &model.Track{
			Scope: m.Scope,
			Track: &track,
		})
	}
}
