package main

import (
	"errors"
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/trakting/store"
)

func userProfile(user store.User, w http.ResponseWriter, r *http.Request) error {
	return render.Render(w, r, "profile.tmpl", http.StatusOK, struct {
		User store.User
	}{
		User: user,
	})
}

func userUpdate(user store.User, w http.ResponseWriter, r *http.Request) error {
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
