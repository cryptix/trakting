package views

import (
	"fmt"

	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"

	"github.com/cryptix/trakting/frontend/model"
)

func trackList(m *model.TrackList, l *ListListeners) dom.Aspect {
	return elem.Section(
		prop.Id("tt-main"),

		bind.Dynamic(m.Scope, func(aspects *bind.Aspects) {
			for _, item := range m.Tracks {
				if !aspects.Reuse(item) {
					theTrack := item
					playing := func() bool { return theTrack == m.PlayingTrack }
					aspects.Add(item, trackElem(item, playing, l))
				}
			}
		}),
	)
}
func trackElem(track *model.Track, playing func() bool, l *ListListeners) dom.Aspect {
	return elem.Div(prop.Class("tt-track"),

		bind.IfFunc(playing, track.Scope,
			prop.Class("playing"),
		),

		elem.Header2(prop.Class("tt-track-title"),
			dom.Text(track.Name),
		),

		elem.Paragraph(prop.Class("tt-track-meta"),
			dom.Text(track.Added.Format("2006-02-01")),
			dom.Text(" by "),
			elem.Anchor(prop.Href("#/list/usr1"), dom.Text(track.By)),
		),

		elem.Paragraph(prop.Class("tt-track-player"),
			elem.Button(prop.Class("btn", "btn-default"),
				elem.Span(prop.Class("glyphicon", "glyphicon-cloud-download")),
				dom.Text("Load"),
				event.Click(l.Play(track)),
			),

			elem.Audio(
				dom.SetProperty("controls", "test"),
				bind.IfFunc(playing, track.Scope,
					prop.Src(fmt.Sprintf("/fetch/%s", track.BoomID)),
				),
			),
		),
	)
}
