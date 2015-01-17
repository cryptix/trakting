package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
)

var client *goBoom.Client

func init() {
	start := time.Now()
	client = goBoom.NewClient(nil)

	code, _, err := client.User.Login("email", "clearPassword")
	logging.CheckFatal(err)

	log.Printf("Login Response: %d (took %v)\n", code, time.Since(start))
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
				log.Println("Listing ", wd)

				_, ls, err := client.Info.Ls(wd)
				logging.CheckFatal(err)
				for _, item := range ls.Items {
					log.Printf("%8s - %s\n", item.ID, item.Name())
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

				log.Println("uploading ", file)
				stats, err := client.FS.Upload(filepath.Base(fname), file)
				logging.CheckFatal(err)
				for _, item := range stats {
					log.Printf("%8s - %s\n", item.ID, item.Name())
				}
			},
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a file",
			Action: func(c *cli.Context) {

				item := c.Args().First()
				if item == "" {
					println("no item id")
					os.Exit(1)
				}

				log.Println("Requesting link for", item)
				_, url, err := client.FS.Download(item)
				logging.CheckFatal(err)
				println(url.String())
			},
		},
	}

	app.Run(os.Args)
}
