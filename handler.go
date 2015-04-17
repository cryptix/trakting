package main

import (
	"errors"
	"fmt"
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

//go:generate gopherjs build -m -o public/js/app.js github.com/cryptix/trakting/frontend
//go:generate go-bindata -pkg=$GOPACKAGE tmpl/... public/...

func init() {
	render.Init(Asset, []string{"tmpl/base.tmpl", "tmpl/navbar.tmpl"})
	render.AddTemplates([]string{
		"tmpl/error.tmpl",
		"tmpl/index.tmpl",
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
		if _, err := ah.AuthenticateRequest(r); err != nil {
			l.WithField("addr", r.RemoteAddr).Error(errgo.Notef(err, "AuthenticateRequest failed"))
			http.Redirect(w, r, "/start", http.StatusTemporaryRedirect)
			return
		}

		fmt.Fprint(w, `<!doctype html>
<html lang="en" data-framework="jquery">
<head>
	<title>Trakting * Loading</title>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta charset="utf-8" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <link rel="stylesheet" href="/public/css/bootstrap.min.css">
    <link rel="stylesheet" href="/public/css/tt.css">
    <link rel="shortcut icon" type="image/png" href="/public/images/favicon.png">
    <script type="text/javascript" src="/public/js/jquery-2.1.0.min.js"></script>
    <script type="text/javascript" src="/public/js/bootstrap.min.js"></script>
    <script type="text/javascript" src="/public/js/app.js"></script>
</head>
<body>
	<div id="app"></div>
</body>
</html>`)
	})

	m.Path("/wsrpc").Handler(websocket.Handler(wsRpcHandler))

	m.Get(Start).Handler(render.StaticHTML("index.tmpl"))
	m.Get(AuthLogin).HandlerFunc(ah.Authorize)
	m.Get(AuthLogout).HandlerFunc(ah.Logout)

	// protected

	m.Get(Upload).Handler(render.Binary(upload))
	m.Get(Fetch).Handler(ah.Authenticate(pushMuxVarsToReqUrl(boomProxy)))

	return m
}

func wsRpcHandler(conn *websocket.Conn) {
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
		return errors.New("user type conversion error")
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data;") {
		return errors.New("illegal content-type")
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
		return errors.New("no stat returned.. really weird error")
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
