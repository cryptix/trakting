package views

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/prop"
)

type Upload struct{}

func (v *Upload) Render() dom.Aspect {
	return elem.Div(prop.Class("container"),
		pageHeader("Upload", "moar trackz plzzz...!"),
		elem.Div(prop.Class("row"),

			elem.Div(prop.Class("col-md-6"),
				elem.Paragraph(
					dom.Text("Select File:"),
					elem.Input(prop.Class("form-control"), prop.Type(prop.TypeFile)),
					elem.Input(prop.Class("btn", "btn-primary"),
						prop.Type(prop.TypeButton),
						prop.Value("Upload!")),
				),
			),

			elem.Div(prop.Class("col-md-6"),
				elem.Div(prop.Class("panel", "panel-default"), prop.Id("uploadPanel"),
					elem.Div(prop.Class("panel-heading"), dom.Text("Upload Status")),
					elem.Div(prop.Class("panel-body"),
						elem.Paragraph(
							prop.Id("uploadStatus"),
							dom.Text("Not started")),
						elem.Div(prop.Class("progress", "progress-striped", "active"),
							elem.Div(prop.Class("progress-bar", "progress-bar-striped")), //role="progressbar"
						),
						// <script type="text/javascript" src="/public/js/uploadui.js"></script>
					),
				),
			),
		),
	)
}
