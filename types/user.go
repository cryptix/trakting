package types

import (
	"fmt"

	"github.com/cryptix/go/http/auth"
)

type User struct {
	ID     int64
	Name   string
	Level  int
	PwHash []byte
}

func (u *User) String() string {
	return fmt.Sprintf("User(%d) %q (lvl %d)", u.ID, u.Name, u.Level)
}

type Userer interface {
	Add(name, passw string, level int) error
	ChangePassword(id int64, newpw string) error
	Current() (*User, error)
	List() ([]User, error)
	auth.Auther
}

type ArgAddUser struct {
	Name, Passw string
	Level       int
}

type ArgChangePassword struct {
	ID    int64
	Passw string
}
