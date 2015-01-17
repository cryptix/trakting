package store

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/cryptix/go/http/auth"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gob.Register(User{})
}

const userBucket = "users"

type User struct {
	Name   string
	Level  int
	PwHash []byte
}

type UserStore struct {
	db *bolt.DB
}

func NewUserStore(db *bolt.DB) (*UserStore, error) {
	if db == nil {
		return nil, errors.New("Init() first")
	}

	return &UserStore{db}, nil
}

func (u *UserStore) Add(name, passw string, level int) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		var (
			err error
			buf bytes.Buffer
			u   User
		)
		u.Level = level
		u.Name = name

		u.PwHash, err = bcrypt.GenerateFromPassword([]byte(passw), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		err = gob.NewEncoder(&buf).Encode(u)
		if err != nil {
			return err
		}

		tx.CreateBucket([]byte(userBucket))
		b := tx.Bucket([]byte(userBucket))
		b.Put([]byte(name), buf.Bytes())

		return nil
	})
}

func (u *UserStore) Check(name, pass string) (interface{}, error) {
	var user User
	err := u.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(userBucket)).Get([]byte(name))
		if v == nil {
			return auth.ErrBadLogin
		}

		err := gob.NewDecoder(bytes.NewReader(v)).Decode(&user)
		if err != nil {
			return err
		}

		err = bcrypt.CompareHashAndPassword(user.PwHash, []byte(pass))
		if err != nil {
			return auth.ErrBadLogin
		}

		return nil
	})

	return user, err
}
