package goBoom

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/tools/godoc/vfs"
	"gopkg.in/errgo.v1"
	"sourcegraph.com/sourcegraph/rwvfs"
)

type VirtualFSService struct {
	fs   *FilesystemService
	info *InformationService

	parent string
}

var _ (rwvfs.FileSystem) = (*VirtualFSService)(nil)

// new rwvfs API

func (v *VirtualFSService) String() string {
	return fmt.Sprintf("boomVFS(%s)", v.parent)
}

func (v *VirtualFSService) Open(path string) (vfs.ReadSeekCloser, error) {
	// FIXME: insert map lookup from full/file/path to (id,name)
	url, err := v.fs.Download(path)
	if err != nil {
		return nil, errgo.Notef(err, "vfs.Open(%s) fs.Download failed", path)
	}

	resp, err := v.fs.c.c.Get(url.String())
	if err != nil {
		return nil, errgo.Notef(err, "vfs.Open(%s) http.Get failed", path)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errgo.Newf("vfs.Open(%s) http.Get: Status not OK: %s", path, resp.Status)
	}

	if resp.ContentLength > 10*1024*1024 {
		return nil, errgo.Newf("vfs.Open(%s) larger than 10meg - not supported", path)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errgo.Notef(err, "vfs.Open(%s) ReadAll(body) failed", path)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, errgo.Notef(err, "vfs.Open(%s) body.Close() failed", path)
	}

	// FIXME: alternatively make a ReadSeekCloser that layzly does range requests

	return &readSeekCloser{bytes.NewReader(b)}, errgo.New("TODO")
}

type readSeekCloser struct{ *bytes.Reader }

func (r *readSeekCloser) Close() error { return nil }

func (v *VirtualFSService) Create(path string) (io.WriteCloser, error) {
	servers, err := v.fs.GetULServer()
	if err != nil {
		return nil, errgo.Notef(err, "GetULServer failed")
	}

	if len(servers) < 1 {
		return nil, errgo.New("no servers available for upload")
	}

	return &fileUpload{fs: v.fs, path: path}, nil
}

type fileUpload struct {
	bytes.Buffer
	fs     *FilesystemService
	path   string
	closed bool
}

func (ul *fileUpload) Close() error {
	if ul.closed {
		return nil
	}
	ul.closed = true

	// FIXME: insert map lookup from full/file/path to (id,name)
	_, err := ul.fs.Upload(filepath.Dir(ul.path), filepath.Base(ul.path), ul)
	return err
}

func (v *VirtualFSService) Lstat(path string) (os.FileInfo, error) {
	// FIXME: insert map lookup from full/file/path to (id,name)
	infos, err := v.info.Info(path)
	if err != nil {
		return nil, errgo.Notef(err, "vfs.Lstat(%s) failed", path)
	}

	if len(infos) != 1 {
		return nil, errgo.Newf("len(infos) != 1. got %d", len(infos))
	}

	return infos[0], nil
}

func (v *VirtualFSService) Stat(path string) (os.FileInfo, error) {
	infos, err := v.info.Info(path)
	if err != nil {
		return nil, errgo.Notef(err, "vfs.Stat(%s) Info() failed", path)
	}

	if len(infos) != 1 {
		return nil, errgo.Newf("len(infos) != 1. got %d", len(infos))
	}

	return infos[0], nil
}

func (v *VirtualFSService) ReadDir(path string) ([]os.FileInfo, error) {

	// FIXME: insert map lookup from full/file/path to (id,name)
	infos, err := v.info.Info(path)
	if err != nil {
		return nil, errgo.Notef(err, "vfs.ReadDir(%s) Info() failed", path)
	}

	finfo := make([]os.FileInfo, 0, len(infos))
	cnt := 0
	for _, i := range infos {
		if i.Parent == "corretParentID" {
			finfo[cnt] = i
			cnt++
		}
	}

	return finfo, nil
}

func (v *VirtualFSService) Mkdir(path string) error {
	// FIXME: insert map lookup from full/file/path to (id,name)
	return v.fs.Mkdir("1", filepath.Base(path))
}

func (v *VirtualFSService) Remove(path string) error {
	// FIXME: insert map lookup from full/file/path to (id,name)
	return v.fs.Rm(true, path)
}
