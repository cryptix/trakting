package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cryptix/goBoom"
	"github.com/cryptix/trakting2/store"
)

type uplinkRequest struct {
	User store.User
	Name string
	Resp chan error
}

var (
	pushUp     chan<- uplinkRequest
	boomClient *goBoom.Client
)

func handlePushes(reqs <-chan uplinkRequest) {
	for r := range reqs {
		go handleUpload(r)
	}
}

func handleUpload(req uplinkRequest) {
	defer func() {
		os.Remove(req.Name)
		close(req.Resp)
	}()
	f, err := os.Open(req.Name)
	if err != nil {
		req.Resp <- err
		return
	}
	defer f.Close()

	stat, err := boomClient.FS.Upload(req.Name, f)
	if err != nil {
		req.Resp <- err
		return
	}

	l.Noticef("uplink done: %v", stat)
	if len(stat) != 1 {
		req.Resp <- errors.New("no stat returned.. really weird error")
		return
	}

	track := store.Track{
		By:     req.User.Name,
		Name:   filepath.Base(req.Name),
		BoomID: stat[0].ID,
	}
	req.Resp <- trackStore.Add(track)
	return
}
