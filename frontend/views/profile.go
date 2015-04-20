package views

import (
	"github.com/cryptix/trakting/frontend/model"
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"github.com/neelance/dom/elem"
	"github.com/neelance/dom/event"
	"github.com/neelance/dom/prop"
)

type Profile struct {
	m *model.Profile
	l *ProfileListeners
}

type ProfileListeners struct {
	Save dom.Listener
}

func NewProfile(m *model.Profile, l *ProfileListeners) *Profile {
	return &Profile{
		m: m,
		l: l,
	}
}

func (v *Profile) Render() dom.Aspect {
	return elem.Div(prop.Class("container"),

		pageHeader("Profile", "Hey "+v.m.Name+", you can change your password here."),
		elem.Div(prop.Class("row")),

		elem.Form(

			dom.PreventDefault(event.Submit(v.l.Save)),

			elem.Div(prop.Class("form-group"),
				elem.Label(prop.Class("control-label"),
					dom.Text("Current Password")),

				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					// dom.SetProperty("autofocus", ""), Doesnt work..
					dom.SetProperty("placeholder", "SuperSecret"),
					bind.Value(&v.m.Current, v.m.Scope)),

				hasError(&v.m.CurrentErr, v.m.Scope),
				bind.IfPtr(&v.m.Success, v.m.Scope, prop.Class("has-success")),
			),

			elem.Div(prop.Class("form-group"),
				elem.Label(prop.Class("control-label"),
					dom.Text("New Password")),

				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					dom.SetProperty("placeholder", "NewPW"),
					bind.Value(&v.m.New, v.m.Scope)),

				hasError(&v.m.NewErr, v.m.Scope),
				bind.IfPtr(&v.m.Success, v.m.Scope, prop.Class("has-success")),
			),

			elem.Div(prop.Class("form-group"),
				elem.Label(prop.Class("control-label"),
					dom.Text("Repeat")),

				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					dom.SetProperty("placeholder", "Repeat"),
					bind.Value(&v.m.Repeat, v.m.Scope)),

				hasError(&v.m.RepeatErr, v.m.Scope),
				bind.IfPtr(&v.m.Success, v.m.Scope, prop.Class("has-success")),
			),

			bind.IfPtr(&v.m.Success, v.m.Scope,
				dom.Text("Saved. please "),
				elem.Anchor(prop.Href("/auth/logout"), dom.Text("Logout")),
				dom.Text(" to use your new password now."),
			),
			elem.Button(prop.Class("btn", "btn-primary", "btn-block"),
				dom.Text("Save"),
			),
		),
	)
}

func hasError(err *error, scope *bind.Scope) dom.Aspect {
	return bind.IfFunc(func() bool { return *err != nil }, scope,
		prop.Class("has-error"),
		elem.Span(prop.Class("help-block"),
			bind.TextFunc(func() string {
				if *err == nil {
					return ""
				}
				return (*err).Error()
			}, scope),
		),
	)
}
