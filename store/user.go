package store

import (
	"encoding/gob"
	"errors"

	"github.com/cryptix/go/http/auth"
	"github.com/jmoiron/modl"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/types"
)

func init() {
	gob.Register(types.User{}) // for auth
	DB.AddTable(types.User{}).SetKeys(true, "id")
	createSql = append(createSql, `alter table "user" ADD UNIQUE (name)`)
}

type UserStore struct {
	dbh modl.SqlExecutor
}

var (
	_ types.Userer = (*UserStore)(nil)
	_ auth.Auther  = (*UserStore)(nil)
)

func NewUserStore() (*UserStore, error) {
	if DBH == nil {
		return nil, errors.New("connect db first")
	}

	return &UserStore{DBH}, nil
}

func (u *UserStore) Add(name, passw string, level int) error {
	var err error
	user := types.User{
		Name:  name,
		Level: level,
	}
	user.PwHash, err = bcrypt.GenerateFromPassword([]byte(passw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.dbh.Insert(&user)
}

func (u *UserStore) Check(name, pass string) (interface{}, error) {
	var user types.User
	err := u.dbh.SelectOne(&user, `SELECT * from "user" WHERE Name = $1`, name)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PwHash, []byte(pass))
	if err != nil {
		return nil, auth.ErrBadLogin
	}

	return user, nil
}

func (u *UserStore) ChangePassword(id int64, newpw string) error {
	var (
		err  error
		user types.User
	)

	err = u.dbh.Get(&user, id)
	if err != nil {
		return err
	}

	user.PwHash, err = bcrypt.GenerateFromPassword([]byte(newpw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.dbh.Update(&user)
	return err
}

func (u *UserStore) Current() (*types.User, error) {
	return nil, errgo.New("not applicable")
}
