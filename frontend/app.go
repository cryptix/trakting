package main

import (
	"github.com/soroushjp/humble/router"
	"github.com/soroushjp/humble/view"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/frontend/wsclient"
)

func main() {
	wc, err := wsclient.New("ws://localhost:3000/wsrpc") // inject correct url somehow
	if err != nil {
		panic(err)
	}
	console.Log("rpc connected")

	r := router.New()

	r.HandleFunc("/", func(_ map[string]string) {
		console.Log("overview")
		overView := &views.Main{Client: wc}
		if err := view.ReplaceParentHTML(overView, "#app"); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/list", func(_ map[string]string) {
		console.Log("list for all")
		listView := &views.TrackList{Client: wc}
		if err := view.ReplaceParentHTML(listView, "#app"); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/list/{user}", func(params map[string]string) {
		user, ok := params["user"]
		if !ok {
			console.Warn("no user => all")
		}
		console.Log("list for", user)
	})

	r.HandleFunc("/profile", func(_ map[string]string) {
		console.Log("profile..")
		profView := &views.Profile{Client: wc}
		if err := view.ReplaceParentHTML(profView, "#app"); err != nil {
			panic(err)
		}
	})

	r.Start()
}
