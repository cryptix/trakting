package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/trakting/store"
)

var l = logging.Logger("main")

func main() {
	app := cli.NewApp()
	app.Name = "trakting"
	app.Version = "0.1"
	app.Usage = "best of music"
	app.Commands = []cli.Command{
		{
			Name:   "addUser",
			Usage:  "adds a user to the db",
			Action: addUserCmd,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "level,l"},
			},
		},
		{
			Name:   "createdb",
			Action: createDbCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "drop,d",
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "starts the webserver",
			Action: serveCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "reload,r",
					Usage: "reload templates on each request?",
				},
				cli.BoolFlag{
					Name:  "ssl,s",
					Usage: "listen using ssl - TODO: configure flags for cert file",
				},
			},
		},
	}
	// app.Action = serveCmd

	logging.SetupLogging(nil)

	var err error
	store.Connect()

	userStore, err = store.NewUserStore()
	logging.CheckFatal(err)

	trackStore, err = store.NewTrackStore()
	logging.CheckFatal(err)

	app.Run(os.Args)
}

func createDbCmd(ctx *cli.Context) {
	if ctx.Bool("drop") {
		store.Drop()
		l.Info("db dropped")
	}
	store.Create()
	l.Info("db created")
}

func addUserCmd(ctx *cli.Context) {
	a := ctx.Args()

	user := a.First()
	if user == "" {
		l.Fatal("we need a username")
	}

	if len(a.Tail()) != 1 {
		l.Fatal("we need a password...")
	}

	pass := a.Tail()[0]

	logging.CheckFatal(userStore.Add(user, pass, ctx.Int("level")))
	l.Info("User added.")
}
