package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/prop"
)

func info() dom.Aspect {
	return elem.Footer(
		prop.Id("info"),

		elem.Paragraph(
			dom.Text("listen and stuff"),
		),
		elem.Paragraph(
			dom.Text("build with"),
			elem.Anchor(
				prop.Href("https://github.com/gopherjs/gopherjs"),
				dom.Text("gopherjs"),
			),
		),
		elem.Paragraph(
			dom.Text("Created by "),
			elem.Anchor(
				prop.Href("http://github.com/cryptix"),
				dom.Text("Richard Musiol"),
			),
		),
	)
}
