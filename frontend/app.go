package main

import (
	"time"

	"github.com/neelance/dom"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/controllers"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/frontend/wsclient"
	"github.com/cryptix/trakting/router"
)

func main() {
	wc, err := wsclient.New("ws://localhost:3000/wsrpc") // inject correct url somehow
	check(err)
	console.Log("rpc connected")

	list, err := controllers.NewList(wc)
	check(err)

	r, err := router.New(list)
	// router.Mode("history"),
	// router.Delay(5*time.Second),

	check(err)

	r.Add("list", list)
	r.Add("upload", &views.Upload{})
	r.Add("profile", &views.Profile{})

	go r.Listen(func(match string, ren router.Renderer) {
		dom.SetTitle("Trakting • " + match)
		dom.SetBody(ren.Render())
	})

	go func() {
		time.Sleep(1 * time.Second)
		r.Navigate("start")
	}()
}

func check(err error) {
	if err != nil {
		console.Error(err)
		panic(err)
	}
}
