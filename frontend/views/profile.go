package views

import (
	"github.com/soroushjp/humble"
	"github.com/soroushjp/humble/view"

	"github.com/cryptix/trakting/rpcClient"
)

type Profile struct {
	humble.Identifier

	Navbar *Navbar
	Client *rpcClient.Client
}

func (p *Profile) RenderHTML() string {
	return `<div id="navbar"></div>
<div class="container">
<div class="page-header">
  <h1>Profile <small>update password and stuff</small></h1>
</div>

<form method="POST" action="#">
	<label for="inputCurrent">Current Password</label>
    <input type="password" id="inputCurrent" class="form-control" placeholder="SuperSecret" required="" autofocus="" name="current">

    <label for="inputNewPW">NewPassword</label>
    <input type="password" class="form-control" placeholder="New PW" required="" id="inputNewPW" name="new">
    <input type="password" class="form-control" placeholder="Repeat" required="" name="repeat">

    <button class="btn btn-lg btn-primary btn-block" type="submit">Save</button>
</form>
</div>
`
}

func (p *Profile) OuterTag() string {
	return "div"
}

func (p *Profile) OnLoad() error {
	var err error
	p.Navbar, err = NewNavbar("profile", p.Client.Users)
	if err != nil {
		return err
	}
	return view.ReplaceParentHTML(p.Navbar, "#navbar")
}
