# gocrayons - REST api consumer with simplejson
[![GoDoc](https://godoc.org/github.com/cryptix/gocrayons?status.svg)](https://godoc.org/github.com/cryptix/gocrayons)
[![Build Status](https://travis-ci.org/cryptix/gocrayons.svg?branch=master)](https://travis-ci.org/cryptix/gocrayons)

## Summary
gocrayons is a fork of [bndr/gopencils](https://github.com/bndr/gopencils). The main change is the use of [github.com/bitly/go-simplejson](bitly/go-simplejson) for JSON handling.
## Install

    go get github.com/cryptix/gocrayons

## Usage example

Please see the `examples` directory.

## Why?

simplejson makes it easier to work with strange and irregular json that can't be simple handled by `encoding/json`'s `Unmarshal()`

## Is it ready?

It is more beta than bndrs original.

## Contribute

All Contributions are welcome. The todo list is on the bottom of this README. Feel free to send a pull request.

## License

Apache License 2.0

## TODO
1. Add examples
