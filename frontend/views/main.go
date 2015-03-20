package views

import (
	"github.com/soroushjp/humble"
	"github.com/soroushjp/humble/view"

	"github.com/cryptix/trakting/rpcClient"
)

type Main struct {
	humble.Identifier

	Navbar *Navbar
	Client *rpcClient.Client
}

func (v *Main) RenderHTML() string {
	return `<div id="navbar"></div>
<div class="container">
<div class="page-header">
  <h1>Hello <small>welcome to tracking</small></h1>
</div>
</div>
`
}

func (v *Main) OuterTag() string {
	return "div"
}

func (v *Main) OnLoad() error {
	var err error
	v.Navbar, err = NewNavbar("profile", v.Client.Users)
	if err != nil {
		return err
	}
	return view.ReplaceParentHTML(v.Navbar, "#navbar")
}
