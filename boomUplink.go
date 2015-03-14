package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/cryptix/goBoom"
)

var boomClient *goBoom.Client

var boomProxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {

		id := req.URL.Query().Get("id")
		if id == "" {
			l.Warningln("boomProxy id missing")
			return
		}

		link, err := boomClient.FS.Download(id)
		if err != nil {
			l.Errorf("boomProxy Download failed: %v", err)
			return
		}

		req.URL.Scheme = link.Scheme
		req.URL.Path = "/" + link.Path
		req.URL.Host = link.Host
		req.URL.RawQuery = link.RawQuery
	},
}
