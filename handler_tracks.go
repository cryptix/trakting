package main

import (
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/trakting/store"
	"github.com/gorilla/mux"
)

func list(user store.User, w http.ResponseWriter, r *http.Request) error {
	tracks, err := trackStore.All()
	if err != nil {
		return err
	}

	return render.Render(w, r, "list.tmpl", http.StatusOK, map[string]interface{}{
		"User":   user,
		"By":     "All",
		"Tracks": tracks,
	})
}

func listByUser(user store.User, w http.ResponseWriter, r *http.Request) error {
	qry := mux.Vars(r)["user"]
	if qry == "" {
		qry = user.Name
	}

	tracks, err := trackStore.ByUserName(qry)
	if err != nil {
		return err
	}

	return render.Render(w, r, "list.tmpl", http.StatusOK, map[string]interface{}{
		"User":   user,
		"By":     qry,
		"Tracks": tracks,
	})
}
