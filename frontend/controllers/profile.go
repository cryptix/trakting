package controllers

import (
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"honnef.co/go/js/console"

	"github.com/cryptix/trakting/frontend/model"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/rpcClient"
)

func NewProfile(c *rpcClient.Client) (*views.Profile, error) {
	m := &model.Profile{
		Scope: bind.NewScope(),
	}

	u, err := c.Users.Current()
	if err != nil {
		return nil, err
	}

	m.Name = u.Name

	lis := &views.ProfileListeners{}

	lis.Save = func(ctx *dom.EventContext) {
		defer m.Scope.Digest()
		if !m.Valid() {
			console.Error("form invalid...")
			return
		}

		u, err := c.Users.Current()
		if err != nil {
			m.CurrentErr = err
			console.Error(err)
			return
		}

		_, err = c.Users.Check(u.Name, m.Current)
		if err != nil {
			m.CurrentErr = err
			console.Error(err)
			return
		}

		err = c.Users.ChangePassword(u.ID, m.New)
		if err != nil {
			m.NewErr = err
			console.Error(err)
			return
		}

		m.Success = true
		return
	}

	return views.NewProfile(m, lis), nil
}
