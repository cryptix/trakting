// +build js

package main

import (
	"github.com/soroushjp/humble/router"
	"github.com/soroushjp/humble/view"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/frontend/wsclient"
)

func main() {
	console.Log("Starting...")
	wc, err := wsclient.New("ws://localhost:3000/wsrpc")
	if err != nil {
		panic(err)
	}
	console.Warn("new client")

	//Start main app view, appView
	appView := &views.Main{Client: wc}
	if err := view.ReplaceParentHTML(appView, "#app"); err != nil {
		panic(err)
	}

	r := router.New()
	r.HandleFunc("/", func(params map[string]string) {
		appView.Init()
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		// if err := view.Update(appView.Footer); err != nil {
		// 	panic(err)
		// }
	})
	r.Start()
}
