package main

import (
	"log"

	"github.com/cryptix/goBoom"
	"github.com/kr/pretty"
)

func main() {
	client := goBoom.NewClient(nil)

	code, resp, err := client.User.Login("el.rey.de.wonns@gmail.com", "70e878c4")
	check(err)

	log.Printf("Login Response[%d] %# v\n", code, pretty.Formatter(resp))

	code, duMap, err := client.Info.Du()
	check(err)
	log.Printf("Du() Response[%d] %# v\n", code, pretty.Formatter(duMap))

	code, info, err := client.Info.Info("1", "1C")
	check(err)
	log.Printf("Info() Response[%d] %# v\n", code, pretty.Formatter(info))

	code, ls, err := client.Info.Ls("1")
	check(err)
	log.Printf("Ls(#2) Response[%d] %# v\n", code, pretty.Formatter(ls))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
