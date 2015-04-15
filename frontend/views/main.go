package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"

	"github.com/cryptix/trakting/frontend/model"
)

type PageListeners struct {
	Search     dom.Listener
	TogglePlay func(*model.Track) dom.Listener
	QueueAll   dom.Listener
}

func Page(m *model.TrackList, l *PageListeners) dom.Aspect {
	return dom.Group(
		elem.Section(

			prop.Id("traktingapp"),

			listHeader(m, l),
			bind.IfFunc(func() bool { return len(m.Tracks) != 0 }, m.Scope,
				trackList(m, l),
				listFooter(m, l),
			),
		),
		info(),
	)
}

func trackList(m *model.TrackList, l *PageListeners) dom.Aspect {
	return elem.Section(
		prop.Id("main"),

		elem.Button(
			prop.Id("queue-all"),
		),

		elem.UnorderedList(
			prop.Id("track-list"),

			bind.Dynamic(m.Scope, func(aspects *bind.Aspects) {
				for _, item := range m.Tracks {
					if !aspects.Reuse(item) {
						theTrack := item
						playing := func() bool { return theTrack == m.PlayingTrack }
						aspects.Add(item, trackElem(item, playing, l))
					}
				}
			}),
		),
	)
}

func trackElem(track *model.Track, playing func() bool, l *PageListeners) dom.Aspect {
	return elem.ListItem(
		bind.IfFunc(playing, track.Scope,
			prop.Class("playing"),
		),

		elem.Div(
			prop.Class("view"),

			// elem.Input(
			// 	prop.Class("toggle"),
			// 	prop.Type(prop.TypeCheckbox),
			// 	bind.Checked(&track.Completed, track.Scope),
			// ),
			elem.Button(
				prop.Class("playpause"),
				event.Click(l.TogglePlay(track)),
			),
			elem.Label(
				bind.TextPtr(&track.Name, track.Scope),
				// event.DblClick(l.StartEdit(track)),
			),
		),
		// elem.Form(
		// 	style.Margin(style.Px(0)),
		// 	dom.PreventDefault(event.Submit(l.StopEdit)),
		// 	elem.Input(
		// 		prop.Class("edit"),
		// 		bind.Value(&track.Title, track.Scope),
		// 	),
		// ),
	)
}
