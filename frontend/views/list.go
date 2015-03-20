package views

import (
	"github.com/soroushjp/humble"
	"github.com/soroushjp/humble/view"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/rpcClient"
)

type TrackList struct {
	humble.Identifier

	Navbar *Navbar
	Client *rpcClient.Client

	Tracks []*Track
}

const (
	trackListSelector = "ul#tracks"
)

func (v *TrackList) RenderHTML() string {
	return `<div id="navbar"></div>
<div class="container">
<div class="page-header">
  <h1>List <small>uploaded by TODO</small></h1>
</div>
<ul id="tracks"></ul>
</div>
`
}

func (v *TrackList) OuterTag() string {
	return "div"
}

func (v *TrackList) OnLoad() error {
	tracks, err := v.Client.Tracks.All()
	if err != nil {
		return err
	}
	for _, track := range tracks {
		console.Dir(track)
		tv := &Track{
			Track:  track,
			Parent: v,
		}
		v.Tracks = append(v.Tracks, tv)
	}

	// Add each child view to the DOM
	for _, tv := range v.Tracks {
		view.AppendToParentHTML(tv, trackListSelector)
	}

	v.Navbar, err = NewNavbar("profile", v.Client.Users)
	if err != nil {
		return err
	}
	return view.ReplaceParentHTML(v.Navbar, "#navbar")
}
