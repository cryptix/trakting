package rpcClient

import (
	"net/rpc"

	"github.com/cryptix/trakting/types"
)

type users struct {
	client *rpc.Client
}

var _ types.Userer = (*users)(nil)

func NewUsersClient(c *rpc.Client) (types.Userer, error) {
	return &users{
		client: c,
	}, nil
}

func (u *users) Add(name, pass string, lvl int) error {
	args := types.ArgAddUser{
		Name:  name,
		Passw: pass,
		Level: lvl,
	}
	return u.client.Call("UserService.Add", args, nil)
}

func (u *users) ChangePassword(id int64, passw string) error {
	args := types.ArgChangePassword{
		ID:    id,
		Passw: passw,
	}
	return u.client.Call("UserService.ChangePassword", args, nil)
}

func (u *users) Current() (*types.User, error) {
	var user types.User
	err := u.client.Call("UserService.Current", "", &user)
	return &user, err
}

func (u *users) Check(name, pass string) (interface{}, error) {
	args := types.ArgAddUser{
		Name:  name,
		Passw: pass,
	}
	return nil, u.client.Call("UserService.Check", args, nil)
}
