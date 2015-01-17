package main

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/trakting2/store"
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

	db, err := bolt.Open("trakting2.db", 0600, nil)
	logging.CheckFatal(err)
	defer db.Close() // TODO: add notify of kill signal...

	userStore, err = store.NewUserStore(db)
	logging.CheckFatal(err)

	trackStore, err = store.NewTrackStore(db)
	logging.CheckFatal(err)

	app.Run(os.Args)
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
	l.Notice("User added.")
}
