// +build !dev

//go:generate gopherjs build -m -o assets/js/app.js github.com/cryptix/trakting/frontend
//go:generate go run assets_gen.go assets.go

package main

const production = true
