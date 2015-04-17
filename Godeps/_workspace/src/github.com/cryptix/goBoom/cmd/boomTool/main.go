package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
)

var (
	client *goBoom.Client
	l      = logging.Logger("boomTool")
)

func init() {
	logging.SetupLogging(nil)

	start := time.Now()
	client = goBoom.NewClient(nil)
	_, err := client.User.Login(
		os.Getenv("OBOOM_USER"),
		os.Getenv("OBOOM_PW"))
	logging.CheckFatal(err)

	l.Infof("Login worked.(took %v)", time.Since(start))
}

func main() {
	app := cli.NewApp()
	app.Name = "boomTool"
	app.Commands = []cli.Command{
		{
			Name:  "ls",
			Usage: "list...",
			Action: func(c *cli.Context) {
				wd := c.Args().First()
				l.Info("Listing ", wd)

				ls, err := client.Info.Ls(wd)
				logging.CheckFatal(err)
				for _, item := range ls.Items {
					l.Infof("%8s - %s", item.ID, item.Name())
				}
			},
		},
		{
			Name:      "put",
			ShortName: "p",
			Usage:     "put a file",
			Action: func(c *cli.Context) {
				fname := c.Args().First()

				file, err := os.Open(fname)
				logging.CheckFatal(err)
				defer file.Close()

				l.Info("uploading ", file)
				stats, err := client.FS.Upload("1", filepath.Base(fname), file)
				logging.CheckFatal(err)
				for _, item := range stats {
					l.Infof("%8s - %s", item.ID, item.Name())
				}
			},
		},
		{
			Name:      "mkdir",
			ShortName: "m",
			Usage:     "create a folder",
			Action: func(c *cli.Context) {

				parent := c.Args().First()
				if parent == "" {
					l.Fatal("no parent id")
				}

				name := c.Args().Get(1)
				if parent == "" {
					l.Fatal("empty name")
				}

				err := client.FS.Mkdir(parent, name)
				logging.CheckFatal(err)
			},
		},
		{
			Name:  "rm",
			Flags: []cli.Flag{cli.BoolFlag{Name: "trash,t"}},
			Action: func(c *cli.Context) {
				item := c.Args().First()
				l.Info("deleting ", item)

				err := client.FS.Rm(c.Bool("trash"), item)
				logging.CheckFatal(err)
			},
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a file",
			Action: func(c *cli.Context) {

				item := c.Args().First()
				if item == "" {
					l.Fatal("no item id")
				}

				l.Info("Requesting link for", item)
				url, err := client.FS.Download(item)
				logging.CheckFatal(err)
				l.Info(url.String())
			},
		},
	}

	app.Run(os.Args)
}
