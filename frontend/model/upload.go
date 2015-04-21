package model

import "github.com/neelance/dom/bind"

type Upload struct {
	Scope *bind.Scope

	Status string

	UploadErr, FormErr error

	UploadInflight             bool
	UploadSuccess, FormSuccess bool
}
