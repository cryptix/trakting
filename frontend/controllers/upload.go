package controllers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/neelance/dom"
	"github.com/neelance/dom/bind"
	"gopkg.in/errgo.v1"
	"honnef.co/go/js/console"
	hdom "honnef.co/go/js/dom"

	"github.com/cryptix/trakting/frontend/model"
	"github.com/cryptix/trakting/frontend/views"
	"github.com/cryptix/trakting/rpcClient"
)

//convenience:
var jQuery = jquery.NewJQuery

func NewUpload(c *rpcClient.Client) (*views.Upload, error) {
	formDataConstr := js.Global.Get("FormData")
	if formDataConstr == js.Undefined {
		return nil, errgo.New("NewUpload: FormData not available")
	}

	xhrConstructor := js.Global.Get("XMLHttpRequest")
	if xhrConstructor == js.Undefined {
		return nil, errgo.New("NewUpload: XMLHttpRequest not available")
	}

	m := &model.Upload{
		Scope: bind.NewScope(),
	}

	var bar = jQuery(".progress .progress-bar")

	var reset = func() {
		time.Sleep(5 * time.Second)
		m.FormSuccess = false
		m.UploadSuccess = false
		m.UploadErr = nil
		bar.SetAttr("data-transitiongoal", 0)
		bar.Underlying().Call("progressbar", js.M{"display_text": "fill"})
		m.Scope.Digest()
	}

	lis := &views.UploadListeners{}
	lis.Send = func(ctx *dom.EventContext) {
		defer m.Scope.Digest()
		m.FormErr = nil

		el := hdom.GetWindow().Document().QuerySelector("#tt-file")
		inputEl, ok := el.(*hdom.HTMLInputElement)
		if !ok {
			m.FormErr = errgo.New("#tt-file is not an input element")
			console.Error(m.FormErr)
			return
		}

		files := inputEl.Files()
		if len(files) != 1 {
			m.FormErr = errgo.Newf("wanted 1 file got %d", len(files))
			console.Error(m.FormErr)
			return
		}

		// check form data
		var name = files[0].Get("name").String()
		if name == "" {
			m.FormErr = errgo.New("empty filename...")
			console.Error(m.FormErr)
			return
		}

		var size = files[0].Get("size").Uint64()
		if size <= 0 {
			m.FormErr = errgo.New("weird size..")
			console.Error(m.FormErr)
			return
		}

		m.Status = fmt.Sprintf("Uploading %s (Size: %s)",
			name,
			humanize.Bytes(size),
		)

		data := formDataConstr.New()
		data.Call("append", "file", files[0])

		m.FormSuccess = true
		m.Scope.Digest()

		m.UploadSuccess = false
		m.UploadInflight = false

		var checkUl = func(err error) {
			if err != nil {
				m.UploadSuccess = false
				m.UploadErr = err
				console.Error(err)
				m.Scope.Digest()
			}
		}

		xhr := xhrConstructor.New()

		respCh := make(chan *http.Response)
		errCh := make(chan error)

		// progress callback
		xhr.Get("upload").Call("addEventListener", "progress", func(ctx *dom.EventContext) {
			if ctx.Node.Get("lengthComputable").Bool() {
				p := ctx.Node.Get("loaded").Int() * 100 / ctx.Node.Get("total").Int()
				bar.SetAttr("data-transitiongoal", p) //Math.ceil(evt.loaded*100/evt.total)
				bar.Underlying().Call("progressbar", js.M{"display_text": "fill"})
				m.Scope.Digest()
			} else {
				console.Error("event: length not computable")
				console.Dir(ctx)
			}
		}, true)

		xhr.Set("onload", func() {
			defer m.Scope.Digest()
			m.UploadInflight = true

			header, err := textproto.NewReader(bufio.NewReader(bytes.NewReader([]byte(xhr.Call("getAllResponseHeaders").String() + "\n")))).ReadMIMEHeader()
			if err != nil {
				checkUl(err)
				return
			}
			body := js.Global.Get("Uint8Array").New(xhr.Get("response")).Interface().([]byte)

			respCh <- &http.Response{
				Status:        xhr.Get("status").String() + " " + xhr.Get("statusText").String(),
				StatusCode:    xhr.Get("status").Int(),
				Header:        http.Header(header),
				ContentLength: int64(len(body)),
				Body:          ioutil.NopCloser(bytes.NewReader(body)),
				// Request:       req,
			}
			m.UploadInflight = false
		})

		xhr.Set("onerror", func(e *js.Object) {
			errCh <- errors.New("net/http: XMLHttpRequest failed")
		})

		xhr.Set("onabort", func(e *js.Object) {
			errCh <- errors.New("net/http: request canceled")
		})

		xhr.Call("open", "POST", "/upload", true)
		xhr.Set("responseType", "arraybuffer") // Needs to happen after "open" due to bug in Firefox, see https://github.com/gopherjs/gopherjs/pull/213.

		m.UploadInflight = true
		m.Scope.Digest()

		xhr.Call("send", data)

		select {
		case resp := <-respCh:
			body, err := ioutil.ReadAll(resp.Body)
			checkUl(err)

			if resp.StatusCode != http.StatusCreated {
				console.Error("Error:", string(body))
				return
			}

			m.Status = string(body)
			m.UploadSuccess = true

		case err := <-errCh:
			console.Warn("error received")
			checkUl(err)
		}

		m.UploadInflight = false

		m.Scope.Digest()
		reset()

	}

	return views.NewUpload(m, lis), nil
}
