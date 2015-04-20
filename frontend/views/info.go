package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/prop"
)

func Navbar() dom.Aspect {
	return elem.Navigation(prop.Class("navbar", "navbar-inverse", "navbar-fixed-top"),

		elem.Div(prop.Class("container"),

			elem.Div(prop.Class("navbar-header"),

				elem.Button(prop.Type(prop.TypeButton), prop.Class("navbar-toggle"),
					// TODO: data- doesnt work..
					dom.SetProperty("data-toggle", "collpase"),
					dom.SetProperty("data-target", ".nvarbar-collapse"),

					elem.Span(prop.Class("sr-only"), dom.Text("Toggle navigation")),
					elem.Span(prop.Class("icon-bar")),
					elem.Span(prop.Class("icon-bar")),
					elem.Span(prop.Class("icon-bar")),
				),
				elem.Anchor(prop.Class("navbar-brand"), prop.Href("/"), dom.Text("Trakting")),
			),

			elem.Div(prop.Class("collapse", "navbar-collapse"),
				elem.UnorderedList(prop.Class("nav", "navbar-nav"),
					elem.ListItem(elem.Anchor(prop.Href("#list"), dom.Text("List"))),
					elem.ListItem(elem.Anchor(prop.Href("#upload"), dom.Text("Upload"))),
				),
				elem.UnorderedList(prop.Class("nav", "navbar-nav", "navbar-right"),
					elem.ListItem(elem.Anchor(prop.Href("#profile"), dom.Text("Profile"))),
					elem.ListItem(elem.Anchor(prop.Href("/auth/logout"), dom.Text("Logout"))),
				),
			),
		),
	)
}

func pageHeader(head, lead string) dom.Aspect {
	return elem.Div(prop.Class("tt-header"),
		elem.Header1(prop.Class("tt-title"), dom.Text(head)),
		elem.Paragraph(prop.Class("lead", "tt-description"), dom.Text(lead)),
	)

}

func Footer() dom.Aspect {
	return elem.Footer(prop.Class("tt-footer"),

		elem.Paragraph(
			dom.Text("listen and stuff"),
		),
		elem.Paragraph(
			dom.Text("build with"),
			elem.Anchor(
				prop.Href("https://github.com/gopherjs/gopherjs"),
				dom.Text("gopherjs"),
			),
			dom.Text(" by "),
			elem.Anchor(
				prop.Href("https://twitter.com/neelance"),
				dom.Text("neelance"),
			),
		),
		elem.Paragraph(
			dom.Text("Created by "),
			elem.Anchor(
				prop.Href("http://github.com/cryptix"),
				dom.Text("Henry"),
			),
		),
	)
}
