// router is a minimal browser router for a gopherjs experiement
//
// inspired by http://krasimirtsonev.com/blog/article/A-modern-JavaScript-router-in-100-lines-history-api-pushState-hash-url
package router

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/neelance/dom"
	"gopkg.in/errgo.v1"
	"honnef.co/go/js/console"
)

type Renderer interface {
	Render() dom.Aspect
}

type Router struct {
	mu           sync.Mutex
	routes       map[string]Renderer
	defaultRoute Renderer

	root  string
	mode  bool // true == history API
	delay time.Duration

	location  *js.Object
	window    *js.Object
	history   *js.Object
	decodeURI *js.Object
}

func New(def Renderer, options ...Option) (*Router, error) {
	r := new(Router)

	r.defaultRoute = def
	r.delay = time.Millisecond * 50
	r.root = "/"
	r.routes = make(map[string]Renderer)

	// our helpers from the browser
	r.location = js.Global.Get("location")
	if r.location == nil {
		return nil, errgo.New("location global object not found")
	}

	r.window = js.Global.Get("window")
	if r.location == nil {
		return nil, errgo.New("window global object not found")
	}

	// hope we dont need this actually
	r.decodeURI = js.Global.Get("decodeURI")
	if r.location == nil {
		return nil, errgo.New("decodeURI func not defined")
	}

	// apply our options
	for i, opt := range options {
		if err := opt(r); err != nil {
			return nil, errgo.Notef(err, "option %d failed", i)
		}
	}

	return r, nil
}

// Listen watches the browser fragment and calls the fn when a match is encountered
func (r *Router) Listen(fn func(string, Renderer)) {
	var curr = r.getFragment()
	for {
		time.Sleep(r.delay)
		if curr != r.getFragment() {
			curr = r.getFragment()
			console.Log("changed:", curr)
			if ren, ok := r.Match(curr); ok {
				console.Warn("matched:", curr)
				fn(curr, ren)
			} else {
				console.Warn("not found, defaulting")
				fn("index", r.defaultRoute)
			}
		}
	}
}

// Add adds an renderer to the routes with name
func (r *Router) Add(name string, ren Renderer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[name] = ren
}

// Match checks whether f matches one of the routes
func (r *Router) Match(f string) (Renderer, bool) {
	if f == "" {
		f = r.getFragment()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ren, found := r.routes[f]
	return ren, found
}

// Navigate changes the current url to path
func (r *Router) Navigate(path string) {
	console.Log("naviagting to", path)
	if r.mode {
		r.history.Call("pushState", nil, nil, r.root+clearSlashes(path))
	} else {
		u, err := r.getFragmentFromLocation()
		if err != nil {
			console.Error("getFragmentFromLocation failed:", err)
		}
		u.Fragment = path
		r.window.Set("location", u.String())
	}
}

func (r *Router) getFragment() string {
	var fragment string
	if r.mode {
		r.location = js.Global.Get("location")

		fragment = r.location.Get("pathname").String() + r.location.Get("search").String()
		// console.Log("first", fragment)
		fragment = r.decodeURI.Invoke(fragment).String()
		// console.Log("decoded", fragment)
		fragment = clearSlashes(fragment)
		// console.Log("cleared", fragment)

		// TODO: remove GET parameters

		if r.root != "/" {
			fragment = strings.Replace(fragment, r.root, "", 1)
			console.Log("rooted", fragment)
		}
	} else {
		u, err := r.getFragmentFromLocation()
		if err != nil {
			console.Error("getFragmentFromLocation failed:", err)
		}
		fragment = u.Fragment
	}

	return clearSlashes(fragment)
}

func (r *Router) getFragmentFromLocation() (*url.URL, error) {
	href := r.window.Get("location").Get("href").String()
	if href == "" {
		return nil, errgo.New("empty href..!?")
	}
	u, err := url.Parse(href)
	if err != nil {
		return nil, errgo.Notef(err, "url.Parse failed")
	}
	return u, nil
}

func clearSlashes(p string) string {
	//console.Log("before:", p)
	p = strings.TrimLeft(p, "/")
	p = strings.TrimRight(p, "/")
	//console.Log("after:", p)
	return p

}
