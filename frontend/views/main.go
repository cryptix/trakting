package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/prop"
)

type Upload struct{}

func (v *Upload) Render() dom.Aspect {
	return dom.Group(
		navbar(),
		elem.Div(prop.Class("container"),
			pageHeader("Upload", "moar trackz plzzz...!"),
			elem.Div(prop.Class("row")),
		),
	)
}
