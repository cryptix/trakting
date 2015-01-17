package store

import (
	"encoding/gob"
	"errors"

	"github.com/cryptix/go/http/auth"
	"github.com/jmoiron/modl"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gob.Register(User{}) // for auth
	DB.AddTable(User{}).SetKeys(true, "id")
}

type User struct {
	ID     int64
	Name   string
	Level  int
	PwHash []byte
}

type UserStore struct {
	dbh modl.SqlExecutor
}

func NewUserStore() (*UserStore, error) {
	if DBH == nil {
		return nil, errors.New("connect db first")
	}

	return &UserStore{DBH}, nil
}

func (u *UserStore) Add(name, passw string, level int) error {
	var err error
	user := User{
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
	var user User
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
