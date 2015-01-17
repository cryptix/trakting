package goBoom

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cryptix/gocrayons"
	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/1.0")
	client.baseURL = url
	client.api = gocrayons.Api(server.URL + "/1.0")
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	a := assert.New(t)

	c := NewClient(nil)
	a.NotNil(c)

	a.Equal(c.baseURL.String(), defaultBaseURL, "wrong BaseURL")
	a.Equal(c.userAgent, userAgent, "wrong userAgent")

	a.IsType(c.User, &UserService{})
	a.NotNil(c.User)

	a.NotNil(c.Info, &InformationService{})
	a.NotNil(c.Info)

	a.IsType(c.FS, &FilesystemService{})
	a.NotNil(c.FS)
}

func TestClient_NewHttpFS(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/1.0/tree", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		cpJson(t, w, "_tests/tree.json")
	})

	fs, err := client.NewHTTPFS()
	assert.Nil(t, err)

	isFS := func(http.FileSystem) {}
	isFS(fs)

	_, err = fs.Open("/")
	assert.Nil(t, err)
}
