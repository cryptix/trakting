package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"

	"github.com/cryptix/trakting/frontend/model"
)

type List struct {
	m *model.TrackList
	l *ListListeners
}

type ListListeners struct {
	Reload   dom.Listener
	Search   dom.Listener
	Play     func(*model.Track) dom.Listener
	QueueAll dom.Listener
}

func NewList(m *model.TrackList, l *ListListeners) *List {
	return &List{
		m: m,
		l: l,
	}
}

func (l *List) Render() dom.Aspect {
	return elem.Div(prop.Class("container"),

		pageHeader("Trakting", "soundcloud is to rainy..."),
		elem.Div(prop.Class("row"),

			elem.Div(prop.Class("col-sm-8", "tt-main"),
				bind.IfFunc(func() bool { return len(l.m.Tracks) != 0 }, l.m.Scope,
					trackList(l.m, l.l),
					listNav(l.m, l.l),
				),
			),
			elem.Div(prop.Class("col-sm-3", "col-sm-offset-1", "tt-sidebar"),
				elem.Div(prop.Class("sidebar-module", "sidebar-module-inset"),
					elem.Header4(dom.Text("About")),
					elem.Paragraph(dom.Text("wawawawa"), elem.Emphasis(dom.Text("wat"))),
					elem.Button(prop.Class("btn", "btn-default"),
						elem.Span(prop.Class("glyphicon", "glyphicon-refresh")),
						dom.Text("Reload"),
						event.Click(l.l.Reload)),
				),
				elem.Div(prop.Class("sidebar-module"),
					elem.Header4(dom.Text("Search")),

					elem.Form(
						dom.PreventDefault(event.Submit(l.l.Search)),

						elem.Input(prop.Class("form-control"),
							// dom.SetProperty("disabled", "disabled"),
							prop.Placeholder("What do you want to hear?"),
							prop.Type(prop.TypeText),
							// prop.Autofocus(),
							bind.Value(&l.m.SearchText, l.m.Scope),
						),
					),
				),
				elem.Div(prop.Class("sidebar-module"),
					elem.Header4(dom.Text("By User")),
					elem.OrderedList(prop.Class("list-unstyled"),
						elem.ListItem(elem.Anchor(prop.Href("#by/usr1"), dom.Text("user1"))),
						elem.ListItem(elem.Anchor(prop.Href("#by/usr2"), dom.Text("user2"))),
						elem.ListItem(elem.Anchor(prop.Href("#by/usr3"), dom.Text("user3"))),
					),
				),
			),
		),
	)
}

func listNav(m *model.TrackList, l *ListListeners) dom.Aspect {
	return elem.Footer(
		prop.Id("footer"),

		elem.Span(
			prop.Id("player-stats"),

			elem.Strong(
				bind.TextFunc(bind.Itoa(m.QueueCount), m.Scope),
			),
			bind.IfFunc(func() bool { return m.QueueCount() == 1 }, m.Scope,
				dom.Text(" track total"),
			),
			bind.IfFunc(func() bool { return m.QueueCount() != 1 }, m.Scope,
				dom.Text(" tracks total"),
			),
		),

		// elem.UnorderedList(
		// 	prop.Id("filters"),
		// 	filterButton("All", model.All, m),
		// 	filterButton("Active", model.Active, m),
		// 	filterButton("Completed", model.Completed, m),
		// ),

		// bind.IfFunc(func() bool { return m.CompletedItemCount() != 0 }, m.Scope,
		// 	elem.Button(
		// 		prop.Id("clear-completed"),
		// 		dom.Text("Clear completed ("),
		// 		bind.TextFunc(bind.Itoa(m.CompletedItemCount), m.Scope),
		// 		dom.Text(")"),
		// 		event.Click(l.ClearCompleted),
		// 	),
		// ),
	)
}
