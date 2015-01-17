package goBoom

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fsHandler struct {
	client *Client
	files  map[string]boomFile
}

// NewHTTPFS returns a (net/http).FileSystem
func (c *Client) NewHTTPFS() (*fsHandler, error) {
	tree, _, err := c.Info.Tree("")
	if err != nil {
		return nil, err
	}

	fh := &fsHandler{
		client: c,
	}

	// BUG(Henry): reload map, add lock and subscribe to changes

	fh.files = make(map[string]boomFile, len(tree))
	for _, f := range tree {
		fh.files[f.Iname] = boomFile{
			id:      f.ID,
			info:    f,
			handler: fh,
		}
	}

	return fh, nil
}

// Open implements net/http FileSystem
func (h *fsHandler) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}

	_, name = filepath.Split(name)
	if name == "" {
		name = "public"
	}

	file, ok := h.files[name]
	if !ok {

		return nil, os.ErrNotExist
	}

	if file.info.IsDir() {
		return &file, nil
	}

	if file.info.Size() > 1024*1024*5 {
		return nil, errors.New("can't inline file transfer over 5mb")
	}

	return &file, file.load()
}

type boomFile struct {
	*bytes.Reader
	id            string
	readDirCalled bool
	info          ItemStat
	handler       *fsHandler
}

func (b *boomFile) load() error {
	url, err := b.handler.client.FS.Download(b.id)
	if err != nil {
		return err
	}

	resp, err := b.handler.client.c.Get(url.String())
	if err != nil {
		return err
	}

	err = checkResponse(resp)
	if err != nil {
		return err
	}

	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	b.Reader = bytes.NewReader(c)

	return resp.Body.Close()
}

func (b *boomFile) Close() error {
	b.readDirCalled = false
	return nil
}

func (b *boomFile) Readdir(n int) ([]os.FileInfo, error) {

	if b.readDirCalled {
		return nil, io.EOF
	}

	var finfo []os.FileInfo
	for _, v := range b.handler.files {
		if v.info.Parent == b.id {
			finfo = append(finfo, v.info)
		}
	}
	b.readDirCalled = true
	return finfo, nil
}

func (b *boomFile) Stat() (os.FileInfo, error) {
	return b.info, nil
}
