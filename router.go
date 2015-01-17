package main

import "github.com/gorilla/mux"

const (
	Start      = "start"
	List       = "list"
	UploadForm = "uploadForm"
	Upload     = "upload"
	Listen     = "listen"
	Fetch      = "fetch"

	UserProfile = "user:profile"

	AuthLogin  = "auth:login"
	AuthLogout = "auth:logout"
)

func App() *mux.Router {
	m := mux.NewRouter()

	m.Path("/start").Methods("GET").Name(Start)
	m.Path("/list").Name(List)
	m.Path("/upload").Methods("GET").Name(UploadForm)
	m.Path("/upload").Methods("POST").Name(Upload)
	m.Path("/listen").Methods("GET").Name(Listen)
	m.Path("/fetch/{id}").Methods("GET").Name(Fetch)

	m.Path("/profile").Methods("GET").Name(UserProfile)

	m.Path("/auth/login").Methods("POST").Name(AuthLogin)
	m.Path("/auth/logout").Methods("GET").Name(AuthLogout)

	return m
}
