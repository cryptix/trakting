// Copyright 2014 Vadim Kravcenko
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package gocrayons is a Golang REST Client with which you can easily consume REST API's. Uses bily/simplejson
package gocrayons

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitly/go-simplejson"
)

var ErrCantUseAsQuery = errors.New("can't use options[0] as Query")

// Resource is basically an url relative to given API Baseurl.
type Resource struct {
	Api         *ApiStruct
	Url         string
	id          string
	QueryValues url.Values
	Payload     io.Reader
	Headers     http.Header
	Response    *simplejson.Json
	Raw         *http.Response
}

// Creates a new Resource.
func (r *Resource) Res(path string) *Resource {
	if path != "" {
		if len(r.Url) > 0 {
			path = r.Url + "/" + path
		}

		newR := &Resource{Url: path, Api: r.Api, Headers: http.Header{}, QueryValues: make(url.Values)}

		return newR
	}
	return r
}

// Same as Res() Method, but returns a Resource with url resource/:id
func (r *Resource) Id(id interface{}) *Resource {
	if id != nil {
		var idStr string
		switch v := id.(type) {
		case string:
			idStr = v
		case int:
			idStr = strconv.Itoa(v)
		case int64:
			idStr = strconv.FormatInt(v, 10)
		default:
			panic("unknown id type")
		}

		url := r.Url + "/" + idStr
		newR := &Resource{id: idStr, Url: url, Api: r.Api, Headers: http.Header{}, Response: r.Response}

		return newR
	}
	return r
}

// Sets QueryValues for current Resource

func (r *Resource) SetQuery(qry map[string]string) *Resource {
	r.QueryValues = make(url.Values)
	for k, v := range qry {
		r.QueryValues.Set(k, v)
	}
	return r
}

// Performs a GET request on given Resource
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Get(params map[string]string) (*Resource, error) {
	r.SetQuery(params)
	return r.do("GET")
}

// Performs a HEAD request on given Resource
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Head(params map[string]string) (*Resource, error) {
	r.SetQuery(params)
	return r.do("HEAD")
}

// Performs a PUT request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Put(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PUT")
}

// Performs a POST request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Post(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("POST")
}

// FormPost doesn't touch the payload
func (r *Resource) FormPost(params map[string]string) (*Resource, error) {
	r.SetQuery(params)
	return r.do("POST")
}

// Performs a Delete request on given Resource.
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Delete(params map[string]string) (*Resource, error) {
	r.SetQuery(params)
	return r.do("DELETE")
}

// Performs a Delete request on given Resource.
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Options(params map[string]string) (*Resource, error) {
	r.SetQuery(params)
	return r.do("OPTIONS")
}

// Performs a PATCH request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Patch(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PATCH")
}

// Main method, opens the connection, sets basic auth, applies headers,
// parses response json.
func (r *Resource) do(method string) (*Resource, error) {
	url := *r.Api.BaseUrl
	if len(url.Path) > 0 {
		url.Path += "/" + r.Url
	} else {
		url.Path = r.Url
	}
	if r.Api.PathSuffix != "" {
		url.Path += r.Api.PathSuffix
	}

	url.RawQuery = r.QueryValues.Encode()
	req, err := http.NewRequest(method, url.String(), r.Payload)
	if err != nil {
		return r, err
	}

	if r.Api.BasicAuth != nil {
		req.SetBasicAuth(r.Api.BasicAuth.Username, r.Api.BasicAuth.Password)
	}

	if r.Headers != nil {
		for k, _ := range r.Headers {
			req.Header.Set(k, r.Headers.Get(k))
		}
	}

	resp, err := r.Api.Client.Do(req)
	if err != nil {
		return r, err
	}

	r.Raw = resp

	if resp.StatusCode >= 400 {
		return r, nil
	}

	defer resp.Body.Close()

	r.Response, err = simplejson.NewFromReader(resp.Body)
	if err != nil {
		return r, err
	}

	return r, nil
}

// Sets Payload for current Resource
func (r *Resource) SetPayload(args interface{}) io.Reader {
	var b []byte
	b, err := json.Marshal(args)
	if err != nil {
		panic(err)
	}
	r.SetHeader("Content-Type", "application/json")
	return bytes.NewBuffer(b)
}

// Sets Headers
func (r *Resource) SetHeader(key string, value string) {
	r.Headers.Add(key, value)
}

// Overwrites the client that will be used for requests.
// For example if you want to use your own client with OAuth2
func (r *Resource) SetClient(c *http.Client) {
	r.Api.Client = c
}
