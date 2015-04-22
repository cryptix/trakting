package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/cryptix/go/http/auth"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
	"github.com/cryptix/trakting/store"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/unrolled/secure.v1"
)

var (
	ah         *auth.Handler
	userStore  *store.UserStore
	trackStore *store.TrackStore
)

func serveCmd(ctx *cli.Context) {
	boomClient = goBoom.NewClient(nil)

	_, err := boomClient.User.Login(
		os.Getenv("OBOOM_USER"),
		os.Getenv("OBOOM_PW"))
	logging.CheckFatal(err)

	l.Info("boomClient.Login done")

	var s store.Settings

	err = store.DBH.SelectOne(&s, `SELECT * FROM appsettings`)
	logging.CheckFatal(err)

	if len(s.HashKey) < 32 || len(s.BlockKey) < 32 {
		l.Error("Warning! cookie keys too short, generating new..")
		s.HashKey = securecookie.GenerateRandomKey(32)
		s.BlockKey = securecookie.GenerateRandomKey(32)
	}

	sessStore := &sessions.CookieStore{

		Codecs: []securecookie.Codec{
			securecookie.New(s.HashKey, s.BlockKey),
		},
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			HttpOnly: true,
			Secure:   true,
		},
	}

	sessStore.Options.Secure = ctx.Bool("ssl")

	ah, err = auth.NewHandler(userStore,
		auth.SetStore(sessStore),
		auth.SetLanding("/"),
		auth.SetLifetime(24*time.Hour),
	)
	logging.CheckFatal(err)

	app := negroni.New(
		negroni.NewRecovery(),
		logging.NewNegroni(l.WithField("module", "http")),
	)
	secuirtyHeaders := secure.New(secure.Options{
		// IsDevelopment:         true,
		AllowedHosts:          []string{"trakting.herokuapp.com"},
		STSSeconds:            315360000,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: `default-src 'self'; connect-src 'self' ws://localhost:3000`,
	})
	app.Use(negroni.HandlerFunc(secuirtyHeaders.HandlerFuncWithNext))

	r := App()
	r.PathPrefix("/public/").Handler(handlers.CompressHandler(
		http.StripPrefix("/public/", http.FileServer(
			&assetfs.AssetFS{
				Asset:    Asset,
				AssetDir: AssetDir,
				Prefix:   "public",
			},
		))))

	app.UseHandler(Handler(r))

	listenAddr := ":" + os.Getenv("PORT")
	lis, err := net.Listen("tcp", listenAddr)
	l.Info("Listening on", lis.Addr())

	logging.CheckFatal(http.Serve(lis, app))
}

// 	certPem, err := ioutil.ReadFile("server.crt")
// 	logging.CheckFatal(err)

// 	keyPem, err := ioutil.ReadFile("server.key")
// 	logging.CheckFatal(err)

// 	cert, err := tls.X509KeyPair(certPem, keyPem)
// 	logging.CheckFatal(err)

// 	srv := &http.Server{
// 		TLSConfig: &tls.Config{
// 			Certificates: []tls.Certificate{cert},
// 		},
// 		Handler: app,
// 	}
// 	http2.ConfigureServer(srv, &http2.Server{})
// 	ln, err := net.Listen("tcp", ":443")
// 	logging.CheckFatal(err)

// 	err = srv.Serve(tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig))
// 	logging.CheckFatal(err)

// type tcpKeepAliveListener struct {
// 	*net.TCPListener
// }

// func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
// 	tc, err := ln.AcceptTCP()
// 	if err != nil {
// 		return
// 	}
// 	tc.SetKeepAlive(true)
// 	tc.SetKeepAlivePeriod(3 * time.Minute)
// 	return tc, nil
// }
