package main

import (
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/trakting/store"
)

func list(user store.User, w http.ResponseWriter, r *http.Request) error {
	tracks, err := trackStore.All()
	if err != nil {
		return err
	}

	var data = struct {
		User   store.User
		Tracks []store.Track
	}{
		User:   user,
		Tracks: tracks,
	}

	return render.Render(w, r, "list.tmpl", http.StatusOK, data)
}
