package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/cryptix/go/http/render"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/trakting/rpcServer"
	"github.com/cryptix/trakting/types"
)

const parentUploadFolder = "24RWR71O"

func copyHTMLAsset(w http.ResponseWriter, fname string) {
	f, err := assets.Open(fname)
	if err != nil {
		err = errgo.Notef(err, "asset.Open(%s) failed", fname)
		l.WithField("err", err).Error("copyHTMLAsset failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, f); err != nil {
		err = errgo.Notef(err, "io.Copy(w,f) failed")
		l.WithField("err", err).Error("copyHTMLAsset failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler hooks up the passed mux.Rotuer to this apps http handlers
func Handler(m *mux.Router) http.Handler {
	if m == nil {
		m = mux.NewRouter()
	}
	m.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := ah.AuthenticateRequest(r); err != nil {
			err = errgo.Notef(err, "AuthenticateRequest failed")
			l.WithField("addr", r.RemoteAddr).Error(err)
			copyHTMLAsset(w, "/login.html")
			return
		}

		copyHTMLAsset(w, "/load.html")
	})

	m.Path("/wsrpc").Handler(websocket.Handler(wsRPCHandler))
	m.Path("/upload").Methods("POST").Handler(render.Binary(upload))
	m.Path("/fetch").Methods("GET").Handler(ah.Authenticate(boomProxy))

	m.Path("/auth/login").Methods("POST").HandlerFunc(ah.Authorize)
	m.Path("/auth/logout").Methods("GET").HandlerFunc(ah.Logout)

	return m
}

func wsRPCHandler(conn *websocket.Conn) {
	l = l.WithField("addr", conn.Request().RemoteAddr)
	i, err := ah.AuthenticateRequest(conn.Request())
	defer func(err error) {
		if err != nil {
			l.WithField("err", err).Error("wsRPC AuthenticateRequest failed")
			fmt.Fprintln(conn, err)
			conn.Close()
		}
	}(err)
	if err != nil {
		return
	}
	user, ok := i.(types.User)
	if !ok {
		err = errgo.New("user type conversion error")
		return
	}

	s := rpc.NewServer()

	ts, err := rpcServer.NewTrackService(user, trackStore)
	if !ok {
		err = errgo.Notef(err, "NewTrackService failed")
		return
	}
	s.RegisterName("TrackService", ts)

	us, e := rpcServer.NewUserService(user, userStore)
	if e != nil {
		err = errgo.Notef(err, "NewUserService failed")
		return
	}
	s.RegisterName("UserService", us)

	fmt.Fprintln(conn, "OK")
	conn.PayloadType = websocket.BinaryFrame
	s.ServeConn(conn)
}

func upload(w http.ResponseWriter, r *http.Request) error {
	i, err := ah.AuthenticateRequest(r)
	if err != nil {
		return err
	}
	user, ok := i.(types.User)
	if !ok {
		return errgo.New("user type conversion error")
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data;") {
		return errgo.New("illegal content-type")
	}

	clen, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	stat, err := boomClient.FS.RawUpload(parentUploadFolder, ct, clen, r.Body)
	if err != nil {
		return err
	}

	l.WithField("stat", stat).Infof("uplink done")
	if len(stat) != 1 {
		return errgo.New("no stat returned.. really weird error")
	}

	track := types.Track{
		By:     user.Name,
		Name:   stat[0].Name(),
		BoomID: stat[0].ID,
	}

	if err := trackStore.Add(track); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Upload done.", stat[0].Name())
	return nil
}
