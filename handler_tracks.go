package trakting

import (
	"errors"
	"net/http"

	"github.com/cryptix/trakting/store"
	"github.com/gin-gonic/gin/render"
)

func list(w http.ResponseWriter, r *http.Request) error {
	// allready authenticated
	i, _ := ah.AuthenticateRequest(r)
	user, ok := i.(store.User)
	if !ok {
		return errors.New("type conversion error")
	}

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
