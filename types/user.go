package types

import "fmt"

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
}

type ArgAddUser struct {
	Name, Passw string
	Level       int
}

type ArgChangePassword struct {
	ID    int64
	Passw string
}
