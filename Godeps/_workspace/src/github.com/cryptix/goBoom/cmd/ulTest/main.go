package main

import (
	"log"
	"os"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
	"github.com/kr/pretty"
)

func main() {
	client := goBoom.NewClient(nil)

	code, _, err := client.User.Login("el.rey.de.wonns@gmail.com", "70e878c4")
	logging.CheckFatal(err)

	// log.Printf("Login Response[%d] %# v\n", code, pretty.Formatter(resp))

	file, err := os.Open(os.Args[1])
	logging.CheckFatal(err)
	defer file.Close()

	code, info, err := client.FS.Upload(os.Args[1], file)
	logging.CheckFatal(err)

	log.Printf("Upload[%d] %# v\n", code, pretty.Formatter(info))
}
