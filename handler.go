package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/trakting/store"
	"github.com/gorilla/mux"
)

func init() {
	render.Init(Asset, []string{"tmpl/base.tmpl", "tmpl/navbar.tmpl"})
	render.AddTemplates([]string{
		"tmpl/error.tmpl",
		"tmpl/index.tmpl",
		"tmpl/list.tmpl",
		"tmpl/upload.tmpl",
		"tmpl/listen.tmpl",
	})
}

// ugly hack to access mux.Vars in httputil ReverseProxy Director func
func pushMuxVarsToReqUrl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		qry := req.URL.Query()
		for key, value := range mux.Vars(req) {
			qry.Set(key, value)
		}
		req.URL.RawQuery = qry.Encode()
		next.ServeHTTP(rw, req)
	})
}

func Handler(m *mux.Router) http.Handler {
	if m == nil {
		m = mux.NewRouter()
	}

	m.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		to := "/list"
		if _, err := ah.AuthenticateRequest(r); err != nil {
			to = "/start"
		}

		http.Redirect(w, r, to, http.StatusTemporaryRedirect)
	})
	m.Get(Start).Handler(render.StaticHTML("index.tmpl"))

	m.Get(AuthLogin).HandlerFunc(ah.Authorize)
	m.Get(AuthLogout).HandlerFunc(ah.Logout)

	m.Get(List).Handler(ah.Authenticate(render.HTML(list)))
	m.Get(UploadForm).Handler(ah.Authenticate(render.HTML(uploadForm)))
	m.Get(Upload).Handler(ah.Authenticate(render.HTML(upload)))
	m.Get(Listen).Handler(ah.Authenticate(render.HTML(listen)))
	m.Get(Fetch).Handler(ah.Authenticate(pushMuxVarsToReqUrl(boomProxy)))

	m.Get(UserProfile).Handler(ah.Authenticate(render.HTML(todo)))

	return m
}

func todo(w http.ResponseWriter, r *http.Request) error {
	return errors.New("TODO")
}

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

func uploadForm(w http.ResponseWriter, r *http.Request) error {
	// allready authenticated
	i, _ := ah.AuthenticateRequest(r)
	user, ok := i.(store.User)
	if !ok {
		return errors.New("type conversion error")
	}

	var data = struct {
		User store.User
	}{
		User: user,
	}

	return render.Render(w, r, "upload.tmpl", http.StatusOK, data)
}

func upload(w http.ResponseWriter, r *http.Request) error {
	i, _ := ah.AuthenticateRequest(r)
	user, ok := i.(store.User)
	if !ok {
		return errors.New("type conversion error")
	}

	file, header, err := r.FormFile("fupload")
	if err != nil {
		return err
	}

	fname := filepath.Join(os.TempDir(), header.Filename)
	input, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer input.Close()

	if _, err = io.Copy(input, file); err != nil {
		return err
	}
	l.Notice("upload taken:", fname)

	req := uplinkRequest{
		User: user,
		Name: fname,
		Resp: make(chan error),
	}
	pushUp <- req

	if err = <-req.Resp; err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Upload done.", fname)
	// http.Redirect(w, r, "/list", http.StatusOK)
	return nil
}

func listen(w http.ResponseWriter, r *http.Request) error {
	i, _ := ah.AuthenticateRequest(r)
	user, ok := i.(store.User)
	if !ok {
		return errors.New("type conversion error")
	}

	id := r.URL.Query().Get("t")
	if id == "" {
		return errors.New("missing id parameter")
	}

	t, err := trackStore.Get(id)
	if err != nil {
		return err
	}

	var data = struct {
		User  store.User
		Track store.Track
	}{
		User:  user,
		Track: t,
	}

	return render.Render(w, r, "listen.tmpl", http.StatusOK, data)
}
