package views

import (
	"fmt"

	"github.com/soroushjp/humble"

	"github.com/cryptix/trakting/types"
)

type Navbar struct {
	humble.Identifier
	User   *types.User
	Active string
}

func NewNavbar(a string, u types.Userer) (*Navbar, error) {
	n := &Navbar{}
	n.Active = "profile"
	var err error
	n.User, err = u.Current()
	return n, err
}

func (n *Navbar) RenderHTML() string {
	return fmt.Sprintf(`<div class="navbar navbar-inverse navbar-fixed-top" role="navigation">
  <div class="container">
    <div class="navbar-header">
      <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
        <span class="sr-only">Toggle navigation</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a class="navbar-brand" href="/">Trakting</a>
    </div>
    <div class="collapse navbar-collapse">
      <ul class="nav navbar-nav">
        <li><a href="#/list">List</a></li>
        <li><a href="#/upload">Upload</a></li>
      </ul>
      <ul class="nav navbar-nav navbar-right">
        <li><a href="#/profile">%s</a></li>
        <li><a href="/auth/logout">Logout</a></li>
      </ul>
    </div>
  </div>
</div>`, n.User.Name)
}

func (n *Navbar) OuterTag() string {
	return "div"
}
