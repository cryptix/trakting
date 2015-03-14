// + build js

package views

import (
	"github.com/soroushjp/humble"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/rpcClient"
)

type Main struct {
	humble.Identifier

	Client *rpcClient.Client
}

func (m *Main) RenderHTML() string {
	return `
	<section id="traktingapp">
		<header id="header">
			<h1>Trakting</h1>
		</header>

		<section id="main">
			<p>List</p>
			<button class="new btn btn-primary">New</button>
			<input type="text" name="number"></input>
			<ul id="mul-list"></ul>
		</section>
	</section>
	<footer id="footer"></footer>`
}

func (m *Main) OuterTag() string {
	return "div"
}

func (m *Main) Init() {
	tracks, err := m.Client.Tracks.All()
	if err != nil {
		console.Error(err)
		return
	}

	for _, track := range tracks {
		console.Dir(track)
		// multiView := &Multi{
		// 	Args:   a,
		// 	Result: i,
		// 	Parent: m,
		// }
		// m.addChild(multiView)
	}
}
