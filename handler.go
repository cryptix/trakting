package main

import (
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
//go:generate go-bindata -pkg=$GOPACKAGE public/...

// ugly hack to access mux.Vars in httputil ReverseProxy Director func
func pushMuxVarsToReqURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		qry := req.URL.Query()
		for key, value := range mux.Vars(req) {
			qry.Set(key, value)
		}
		req.URL.RawQuery = qry.Encode()
		next.ServeHTTP(rw, req)
	})
}

const (
	loadHTML = `<!doctype html>
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
		<script type="text/javascript" src="/public/js/bootstrapProgressbar.min.js"></script>
    <script type="text/javascript" src="/public/js/app.js"></script>
</head>
<body>
<a href="#wtf">Click me if you can.</a>
</body>
</html>`

	loginHTML = `<!doctype html>
<html lang="en" data-framework="jquery">
<head>
	<title>Trakting * Loading</title>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta charset="utf-8" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <link rel="stylesheet" href="/public/css/bootstrap.min.css">
    <link rel="stylesheet" href="/public/css/signin.css">
    <link rel="shortcut icon" type="image/png" href="/public/images/favicon.png">
    <script type="text/javascript" src="/public/js/jquery-2.1.0.min.js"></script>
    <script type="text/javascript" src="/public/js/bootstrap.min.js"></script>
</head>
<body>
<div class="container">
	<form class="form-signin" method="POST" action="/auth/login" >
		<h2 class="form-signin-heading">Please sign in</h2>
		<label for="inputUser" class="sr-only">Username</label>
		<input name="user" type="text" id="inputUser" class="form-control" placeholder="Username" required="" autofocus="">
		<label for="inputPassword" class="sr-only">Password</label>
		<input name="pass" type="password" id="inputPassword" class="form-control" placeholder="Password" required="">
		<button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
	</form>
</div>
</body>
</html>`
)

// Handler hooks up the passed mux.Rotuer to this apps http handlers
func Handler(m *mux.Router) http.Handler {
	if m == nil {
		m = mux.NewRouter()
	}

	m.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := ah.AuthenticateRequest(r); err != nil {
			l.WithField("addr", r.RemoteAddr).Error(errgo.Notef(err, "AuthenticateRequest failed"))
			fmt.Fprintf(w, loginHTML)
			return
		}

		fmt.Fprint(w, loadHTML)
	})

	m.Path("/wsrpc").Handler(websocket.Handler(wsRPCHandler))
	m.Path("/upload").Methods("POST").Handler(render.Binary(upload))
	m.Path("/fetch/{id}").Methods("GET").Handler(ah.Authenticate(pushMuxVarsToReqURL(boomProxy)))

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
