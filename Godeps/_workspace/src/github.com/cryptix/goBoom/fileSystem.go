package goBoom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"
)

type FilesystemService struct {
	c *Client
}

func newFilesystemService(c *Client) *FilesystemService {
	i := &FilesystemService{}
	if c == nil {
		i.c = NewClient(nil)
	} else {
		i.c = c
	}

	return i
}

func (s *FilesystemService) GetULServer() ([]string, error) {
	resp, err := s.c.api.Res("ul/server").Get(nil)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	var servers []string
	if err = decodeInto(&servers, arr[1]); err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *FilesystemService) Mkdir(parent, name string) error {
	params := map[string]string{
		"token":  s.c.User.session,
		"name":   name,
		"parent": parent,
	}
	resp, err := s.c.api.Res("mkdir").Get(params)
	_, err = processResponse(resp, err)
	return err
}

func (s *FilesystemService) Rm(toTrash bool, items ...string) error {
	params := map[string]string{
		"token": s.c.User.session,
		"items": strings.Join(items, ","),
	}
	if toTrash == false {
		params["move_to_trash"] = "false"
	} else {
		params["move_to_trash"] = "true"
	}

	resp, err := s.c.api.Res("rm").Get(params)
	_, err = processResponse(resp, err)
	return err
}

// Upload pushes the input io.Reader to the service
func (s *FilesystemService) Upload(parent, fname string, input io.Reader) ([]ItemStat, error) {

	servers, err := s.GetULServer()
	if err != nil {
		return nil, err
	}

	if len(servers) < 1 {
		return nil, errors.New("no servers available for upload")
	}

	var bodyBuf bytes.Buffer
	writer := multipart.NewWriter(&bodyBuf)

	part, err := writer.CreateFormFile("file", filepath.Base(fname))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, input)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// prepare request
	oldHost := s.c.api.Api.BaseUrl.Host
	s.c.api.Api.BaseUrl.Host = strings.Replace(oldHost, "api.oboom.com", servers[0], 1)

	res := s.c.api.Res("ul")
	res.Payload = &bodyBuf
	res.Headers.Set("Content-Type", writer.FormDataContentType())

	// set  token
	params := map[string]string{
		"token":       s.c.User.session,
		"name_policy": "rename",
		"parent":      parent,
	}

	// do the request
	resp, err := res.FormPost(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}
	s.c.api.Api.BaseUrl.Host = oldHost

	var items []ItemStat
	if err = decodeInto(&items, arr[1]); err != nil {
		return nil, err
	}

	return items, nil
}

// RawUpload expects a multipart io.Reader and pushes it to the service
func (s *FilesystemService) RawUpload(parent, ct string, clen int64, multiBody io.Reader) ([]ItemStat, error) {

	servers, err := s.GetULServer()
	if err != nil {
		return nil, err
	}

	if len(servers) < 1 {
		return nil, errors.New("no servers available for upload")
	}

	// prepare request
	req, err := http.NewRequest("POST", "https://"+servers[0]+"/1.0/ul", multiBody)
	if err != nil {
		return nil, err
	}

	req.ContentLength = clen
	req.Header.Set("Content-Type", ct)

	qry := req.URL.Query()
	qry.Set("token", s.c.User.session)
	qry.Set("name_policy", "rename")
	qry.Set("parent", parent)
	req.URL.RawQuery = qry.Encode()

	resp, err := s.c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonResp, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	arr, err := jsonResp.Array()
	if err != nil {
		return nil, err
	}

	if len(arr) < 1 {
		return nil, ErrorResponse{resp, "Illegal oBoom response"}
	}

	statusCode, err := jsonResp.GetIndex(0).Int()
	if err != nil {
		return nil, err
	}

	if statusCode != resp.StatusCode {
		return nil, ErrorResponse{resp, fmt.Sprintf("StatusCode missmatch. %d vs %d\n%+v",
			statusCode,
			resp.StatusCode,
			jsonResp)}
	}

	var items []ItemStat
	if err = decodeInto(&items, arr[1]); err != nil {
		return nil, err
	}

	return items, nil
}

// Download requests a download url for item
func (s *FilesystemService) Download(item string) (*url.URL, error) {
	if s.c.User == nil {
		return nil, errors.New("non pro download not supported")
	}

	params := map[string]string{
		"token": s.c.User.session,
		"item":  item,
	}

	resp, err := s.c.api.Res("dl").Get(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	var (
		u  url.URL
		ok bool
	)

	u.Host, ok = arr[1].(string)
	if !ok {
		return nil, errors.New("arr[1] is not a string")
	}

	ticket, ok := arr[2].(string)
	if !ok {
		return nil, errors.New("arr[2] is not a string")
	}

	u.Scheme = "https"
	u.Path = libraryVersion + "/dlh"

	qry := u.Query()
	qry.Set("ticket", ticket)
	u.RawQuery = qry.Encode()

	return &u, nil
}
