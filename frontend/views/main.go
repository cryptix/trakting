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

type Profile struct{}

func (v *Profile) Render() dom.Aspect {
	return dom.Group(
		navbar(),
		elem.Div(prop.Class("container"),

			pageHeader("Profile", "change pw and stuff..."),
			elem.Div(prop.Class("row")),

			elem.Form(

				// dom.PreventDefault(event.Submit()),

				elem.Label(dom.Text("Current Password")),
				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					dom.SetProperty("placeholder", "SuperSecret"),
				),

				elem.Label(dom.Text("New Password")),
				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					dom.SetProperty("placeholder", "NewPW"),
				),
				elem.Input(prop.Class("form-control"),
					prop.Type(prop.TypePassword),
					dom.SetProperty("placeholder", "Repeat"),
				),
			),

			//	<label for="inputCurrent">Current Password</label>
			//    <input type="password" id="inputCurrent" required="" autofocus="" name="current">
			//
			//    <label for="inputNewPW">NewPassword</label>
			//    <input type="password" required="" id="inputNewPW" name="new">
			//    <input type="password"required="" name="repeat">
			//
			//    <button class="btn btn-lg btn-primary btn-block" type="submit">Save</button>
			//</form>
		),
	)
}
