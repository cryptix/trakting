package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"
	"github.com/neelance/dom/style"

	"github.com/cryptix/trakting/frontend/model"
)

func listHeader(m *model.TrackList, l *PageListeners) dom.Aspect {
	return elem.Header(
		prop.Id("header"),

		elem.Header1(
			dom.Text("Trackting"),
		),
		elem.Form(
			style.Margin(style.Px(0)),
			dom.PreventDefault(event.Submit(l.Search)),

			elem.Input(
				prop.Id("search-track"),
				prop.Placeholder("What do you want to hear?"),
				prop.Autofocus(),
				bind.Value(&m.SearchText, m.Scope),
			),
		),
	)
}

func listFooter(m *model.TrackList, l *PageListeners) dom.Aspect {
	return elem.Footer(
		prop.Id("footer"),

		elem.Span(
			prop.Id("player-stats"),

			elem.Strong(
				bind.TextFunc(bind.Itoa(m.QueueCount), m.Scope),
			),
			bind.IfFunc(func() bool { return m.QueueCount() == 1 }, m.Scope,
				dom.Text(" track left"),
			),
			bind.IfFunc(func() bool { return m.QueueCount() != 1 }, m.Scope,
				dom.Text(" tracks left"),
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
