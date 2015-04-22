package main

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/prop"
	"honnef.co/go/js/console"
	hdom "honnef.co/go/js/dom"

	"github.com/cryptix/trakting/frontend/controllers"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/frontend/wsclient"
	"github.com/cryptix/trakting/router"
)

var document = hdom.GetWindow().Document()

// const wshost = "ws://localhost:3000/wsrpc"
const wshost = "wss://trakting.herokuapp.com/wsrpc"

func main() {
	wc, err := wsclient.New(wshost) // inject correct url somehow
	check(err)
	console.Log("rpc connected")

	list, err := controllers.NewList(wc)
	check(err)

	upload, err := controllers.NewUpload(wc)
	check(err)

	profile, err := controllers.NewProfile(wc)
	check(err)

	r, err := router.New(list)
	// router.Mode("history"),
	// router.Delay(5*time.Second),
	check(err)

	r.Add("list", list)
	r.Add("upload", upload)
	r.Add("profile", profile)

	go r.Listen(func(match string, ren router.Renderer) {
		dom.SetTitle("Trakting • " + match)
		main := document.QuerySelector("#main")

		div := document.CreateElement("div")
		ren.Render().Apply(div.Underlying(), 0, 1)

		main.ReplaceChild(div, main.FirstChild())
	})

	dom.SetTitle("Trakting")
	ren, _ := r.Match("")
	dom.SetBody(views.Navbar(), elem.Div(prop.Id("main"), ren.Render()), views.Footer())

}

func check(err error) {
	if err != nil {
		console.Error(err)
		panic(err)
	}
}
