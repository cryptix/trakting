package main

import (
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/cryptix/go/http/auth"
	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
	"github.com/cryptix/trakting/store"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/graceful"
	"gopkg.in/unrolled/secure.v1"
)

var (
	ah         *auth.Handler
	userStore  *store.UserStore
	trackStore *store.TrackStore
)

func serveCmd(ctx *cli.Context) {
	boomClient = goBoom.NewClient(nil)

	code, err := boomClient.User.Login(
		os.Getenv("OBOOM_USER"),
		os.Getenv("OBOOM_PW"))
	logging.CheckFatal(err)

	l.Noticef("Login Response[%d]\n", code)

	c := make(chan uplinkRequest)
	go handlePushes(c)
	pushUp = c

	var s store.Settings

	err = store.DBH.SelectOne(&s, `SELECT * FROM appsettings`)
	logging.CheckFatal(err)

	if len(s.HashKey) < 32 || len(s.BlockKey) < 32 {
		l.Critical("Warning! cookie keys too short, generating new..")
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
		},
	}

	sessStore.Options.Secure = ctx.Bool("ssl")

	render.Reload = ctx.Bool("reload")
	r := App()
	render.SetAppRouter(r)
	render.Load()

	ah, err = auth.NewHandler(userStore,
		auth.SetStore(sessStore),
		auth.SetLanding("/"),
		auth.SetLifetime(24*time.Hour),
	)
	logging.CheckFatal(err)

	secuirtyHeaders := secure.New(secure.Options{
		AllowedHosts:          []string{"trakting.herokuapp.com"},
		SSLRedirect:           true,
		STSSeconds:            315360000,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IsDevelopment:         true,
	})
	app := negroni.New(
		negroni.NewRecovery(),
		logging.NewNegroni("trakting"),
	)
	app.Use(negroni.HandlerFunc(secuirtyHeaders.HandlerFuncWithNext))

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(
		&assetfs.AssetFS{
			Asset:    Asset,
			AssetDir: AssetDir,
			Prefix:   "public",
		},
	)))
	app.UseHandler(Handler(r))

	listenAddr := ":" + os.Getenv("PORT")
	l.Notice("Listening on", listenAddr)

	srv := &graceful.Server{
		Timeout: 1 * time.Minute,
		Server:  &http.Server{Addr: listenAddr, Handler: app},
	}
	// heroku has it's own ssl reverseproxy
	logging.CheckFatal(srv.ListenAndServe())
}
