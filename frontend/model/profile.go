package model

import (
	"github.com/neelance/dom/bind"
	"gopkg.in/errgo.v1"
)

type Profile struct {
	Scope *bind.Scope

	Name string

	Current, New, Repeat          string
	CurrentErr, NewErr, RepeatErr error

	Success bool
}

func (p *Profile) Valid() bool {
	var valid = true

	p.CurrentErr = nil
	if p.Current == "" {
		p.CurrentErr = errgo.New("Current to short")
		valid = false
	}

	p.NewErr = nil
	if len(p.New) < 5 {
		p.NewErr = errgo.New("New too short")
		valid = false
	}

	p.RepeatErr = nil
	if p.New != p.Repeat {
		p.RepeatErr = errgo.New("not equal")
		valid = false
	}

	return valid
}
