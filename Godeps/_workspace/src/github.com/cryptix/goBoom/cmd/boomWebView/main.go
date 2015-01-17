package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
)

var client *goBoom.Client

func main() {
	client = goBoom.NewClient(nil)

	code, _, err := client.User.Login("xxxx", "xxxxx")
	logging.CheckFatal(err)

	log.Println("Login Response: ", code)

	fs, err := client.NewHTTPFS()
	logging.CheckFatal(err)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(http.FileServer(fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "0"
	}

	l, err := net.Listen("tcp", ":"+port)
	logging.CheckFatal(err)
	log.Printf("Serving at http://%s/", l.Addr())

	logging.CheckFatal(http.Serve(l, n))
}
