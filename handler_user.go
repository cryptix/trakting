package main

import (
	"errors"
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/trakting/types"
)

func userProfile(user types.User, w http.ResponseWriter, r *http.Request) error {
	return render.Render(w, r, "profile.tmpl", http.StatusOK, struct {
		User types.User
	}{
		User: user,
	})
}

func userUpdate(user types.User, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	currentPW := r.PostFormValue("current")

	newPW := r.PostFormValue("new")
	repeatPW := r.PostFormValue("repeat")
	if newPW != repeatPW {
		return errors.New("passwords didn't match")
	}

	if _, err := userStore.Check(user.Name, currentPW); err != nil {
		return err
	}

	if err := userStore.ChangePassword(user.ID, newPW); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}
