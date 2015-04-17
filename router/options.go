package router

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"gopkg.in/errgo.v1"
)

type Option func(*Router) error

func Delay(d time.Duration) Option {
	return func(r *Router) error {
		if d < time.Millisecond*50 {
			return errgo.New("delay very small.. change me if you really want this")

		}
		r.delay = d
		return nil
	}
}

func Root(root string) Option {
	return func(r *Router) error {
		r.root = "/" + clearSlashes(root) + "/"
		return nil
	}
}

func Mode(mode string) Option {
	var history = js.Global.Get("history")
	return func(r *Router) error {
		switch mode {
		case "history":
			if history == nil {
				return errgo.New("history api is not supported")
			}
			r.mode = true
			r.history = history
		case "hash":
			r.mode = false
		default:
			return errgo.Newf("unknown Mode: %s", mode)
		}
		return nil
	}
}
