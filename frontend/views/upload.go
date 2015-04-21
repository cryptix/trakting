package views

import (
	"github.com/cryptix/trakting/frontend/model"
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"
)

type Upload struct {
	m *model.Upload
	l *UploadListeners
}

type UploadListeners struct {
	Send dom.Listener
}

func NewUpload(m *model.Upload, l *UploadListeners) *Upload {
	return &Upload{
		m: m,
		l: l,
	}
}

func (v *Upload) Render() dom.Aspect {
	return elem.Div(prop.Class("container"),
		pageHeader("Upload", "moar trackz plzzz...!"),
		elem.Div(prop.Class("row"),

			elem.Div(prop.Class("col-md-6"),
				elem.Form(
					dom.PreventDefault(event.Submit(v.l.Send)),

					elem.Div(prop.Class("form-group"),
						elem.Label(prop.Class("control-label"),
							dom.Text("Select File:"),
						),

						elem.Input(prop.Class("form-control"),
							prop.Type(prop.TypeFile),
							prop.Id("tt-file"),
						),

						hasError(&v.m.FormErr, v.m.Scope),
						bind.IfPtr(&v.m.FormSuccess, v.m.Scope, prop.Class("has-success")),

						elem.Input(prop.Class("btn", "btn-primary"),
							bind.IfPtr(&v.m.FormSuccess, v.m.Scope, prop.Class("disabled")),
							prop.Type(prop.TypeSubmit),
							prop.Value("Upload!")),
					),
				),
			),

			elem.Div(prop.Class("col-md-6"),
				elem.Div(prop.Class("panel", "panel-default"),

					elem.Div(prop.Class("panel-heading"), dom.Text("Upload Status")),

					bind.IfFunc(func() bool { return v.m.UploadErr != nil }, v.m.Scope,
						prop.Class("panel-danger"),
					),
					bind.IfPtr(&v.m.UploadSuccess, v.m.Scope, prop.Class("panel-success")),
					bind.IfPtr(&v.m.UploadInflight, v.m.Scope, prop.Class("panel-primary")),

					elem.Div(prop.Class("panel-body"),
						bind.IfFunc(func() bool { return v.m.UploadErr != nil }, v.m.Scope,
							elem.Paragraph(bind.TextFunc(func() string {
								if v.m.UploadErr == nil {
									return ""
								}
								return v.m.UploadErr.Error()
							}, v.m.Scope)),
						),

						elem.Paragraph(
							prop.Id("uploadStatus"),
							bind.TextPtr(&v.m.Status, v.m.Scope)),
						elem.Div(prop.Class("progress", "progress-striped"),

							bind.IfPtr(&v.m.UploadInflight, v.m.Scope, prop.Class("active", "progress-primary")),
							bind.IfFunc(func() bool { return v.m.UploadErr != nil }, v.m.Scope,
								prop.Class("progress-danger"),
							),
							elem.Div(prop.Class("progress-bar", "progress-bar-striped"),
								bind.IfPtr(&v.m.UploadSuccess, v.m.Scope, prop.Class("progress-bar-success")),
								bind.IfFunc(func() bool { return v.m.UploadErr != nil }, v.m.Scope,
									prop.Class("progress-bar-danger"),
								),
							), //role="progressbar"
						),
					),
				),
			),
		),
	)
}
