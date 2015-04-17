package main

import (
	"log"

	"github.com/cryptix/goBoom"
	"github.com/shurcooL/go-goon"
)

func main() {
	client := goBoom.NewClient(nil)

	resp, err := client.User.Login("el.rey.de.wonns@gmail.com", "Tinchen!123")
	check(err)

	log.Println("Login Response ok")
	goon.Dump(resp)

	duMap, err := client.Info.Du()
	check(err)
	log.Printf("Du() Response:")
	goon.Dump(duMap)

	info, err := client.Info.Info("1", "1C")
	check(err)
	log.Printf("Info() Response:")
	goon.Dump(info)

	ls, err := client.Info.Ls("1")
	check(err)
	log.Printf("Ls(#2) Response:")
	goon.Dump(ls)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
