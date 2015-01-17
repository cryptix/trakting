package main

import (
	"fmt"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/cryptix/gocrayons"
)

// Oauth example taken from https://godoc.org/code.google.com/p/goauth2/oauth
// The authenticated client must be passed to gopencils for OAuth to work.

type respStruct struct {
	Args          map[string]string
	Headers       map[string]string
	Origin        string
	Url           string
	Authorization string
}

// Specify your configuration. (typically as a global variable)
var config = &oauth.Config{
	ClientId:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	Scope:        "https://www.googleapis.com/auth/buzz",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
	RedirectURL:  "http://you.example.org/handler",
}

// A landing page redirects to the OAuth provider to get the auth code.
func landing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL("foo"), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func handler(w http.ResponseWriter, r *http.Request) {
	t := &oauth.Transport{Config: config}
	t.Exchange(r.FormValue("code"))
	// The Transport now has a valid Token. Create an *http.Client
	// with which we can make authenticated API requests.
	c := t.Client()

	// Now you can pass the authenticated Client to gopencils, and
	// it will be used to make all the requests
	api := gocrayons.Api("http://your-api-url.com/api/", c)

	// Maybe some payload to send along with the request?
	payload := map[string]interface{}{"Key1": "Value1"}

	// Perform a GET request
	// URL Requested: http://your-api-url.com/api/users
	resp, err := api.Res("users").Get(nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%+v\n", resp.Response)

	// Perform a GET request with Querystring
	querystring := map[string]string{"page": "100"}
	// URL Requested: http://your-api-url.com/api/users/123/items?page=100
	resp, err = api.Res("users").Id(123).Res("items").Get(querystring)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%+v\n", resp.Response)

	// Or perform a POST Request
	// URL Requested: http://your-api-url.com/api/items/123 with payload as json Data
	resp, err = api.Res("items").Id(123).Post(payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%+v\n", resp.Response)
}

func main() {
	http.HandleFunc("/", landing)
	http.HandleFunc("/in", handler)
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
