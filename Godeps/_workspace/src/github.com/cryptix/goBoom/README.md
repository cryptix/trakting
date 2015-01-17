goBoom
======
[![GoDoc](https://godoc.org/github.com/cryptix/goBoom?status.svg)](https://godoc.org/github.com/cryptix/goBoom)
[![Build Status](https://travis-ci.org/cryptix/goBoom.svg?branch=master)](https://travis-ci.org/cryptix/goBoom)

golang API for the oboom.com service


## Todo
- [ ] Add `LoginWithHash()` method
- [ ] Add [context](https://golang.org/x/net/context)
- [ ] Add lock to BaseURL to prevent races (login/upload use different hosts)
- [ ] Add tests for Upload()
- [ ] Add Delete Call
- [x] Add Tree Call
- [x] Add map[name]id for FileSystem
- [ ] Add Parent to Upload()
- [ ] Add Mock for service for integration testing of FileSystem
- [ ] Better timeout control