package gocrayons

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var (
	testMux *http.ServeMux
	testSrv *httptest.Server
)

func init() {
	testMux = http.NewServeMux()
	testSrv = httptest.NewServer(testMux)
}

type respStruct struct {
	Login   string
	Id      int
	Name    string
	Message string
}

type httpbinResponse struct {
	Args    string
	Headers map[string]string
	Url     string
	Json    map[string]interface{}
}

func TestResouceUrl(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Api.BaseUrl.String(), "https://test-url.com",
		"Parsed Url Should match")
	api.SetQuery(map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, api.QueryValues.Encode(), "key1=value1&key2=value2",
		"Parsed QueryString Url Should match")
	assert.Equal(t, api.Url, "", "Base Url Should be empty")
}

func TestCanUsePathInResourceUrl(t *testing.T) {
	testMux.HandleFunc("/path/to/api/resname/id123",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, `{"Test":"Okay"}`)
		})

	res := Api(testSrv.URL+"/path/to/api", nil)

	resp, err := res.Res("resname").Id("id123").Get(nil)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "Okay", resp.Response.Get("Test").MustString(), "resp should be Okay")
}

func TestCanUseAuthForApi(t *testing.T) {
	api := Api("https://test-url.com", &BasicAuth{"username", "password"})
	assert.Equal(t, api.Api.BasicAuth.Username, "username",
		"Username should match")
	assert.Equal(t, api.Api.BasicAuth.Password, "password",
		"Password should match")
}

func TestCanGetResource(t *testing.T) {
	// github stubs
	testMux.HandleFunc("/users/bndr",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, readJson("_tests/github_bndr.json"))
		})
	testMux.HandleFunc("/users/torvalds",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, readJson("_tests/github_torvalds.json"))
		})

	api := Api(testSrv.URL)

	// Users endpoint
	users := api.Res("users")

	usernames := []string{"bndr", "torvalds"}

	// construct mapstructure decoder
	r := new(respStruct)
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           r}
	dec, err := mapstructure.NewDecoder(config)
	assert.NoError(t, err, "Error constructing mapstruct")

	for _, username := range usernames {
		// Get user with id i into the newly created response struct
		res, err := users.Id(username).Get(nil)
		assert.NoError(t, err, "Error Getting Data from Test API")

		err = dec.Decode(res.Response.Interface())
		assert.NoError(t, err, "Error decoding mapstruct")

		assert.Equal(t, r.Message, "", "Error message must be empty")
		assert.Equal(t, r.Login, username, "Username should be equal")

	}

	res, err := api.Res("users").Id("bndr").Get(nil)
	assert.NoError(t, err, "Error Getting Data from Test API")

	err = dec.Decode(res.Response.Interface())
	assert.NoError(t, err, "Error decoding mapstruct")

	assert.Equal(t, r.Login, "bndr")

	res, err = api.Res("users").Id("bndr").Get(nil)
	assert.NoError(t, err, "Error Getting Data from Test API")

	err = dec.Decode(res.Response.Interface())
	assert.NoError(t, err, "Error decoding mapstruct")

	assert.Equal(t, r.Login, "bndr")
}

func TestCanCreateResource(t *testing.T) {
	testMux.HandleFunc("/post",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "POST", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/post", "unexpected Path")
			assert.Equal(t, req.Header.Get("Content-Type"), "application/json",
				"Expected json content type")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(httpbinResponse)
	res, err := api.Res("post").Post(payload)
	assert.NoError(t, err, "Error Getting Data from httpbin API")
	r.Json, err = res.Response.Get("json").Map()
	assert.NoError(t, err, "Could not convert responses to map")
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestCanPutResource(t *testing.T) {
	testMux.HandleFunc("/put",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "PUT", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/put", "unexpected Path")
			assert.Equal(t, req.Header.Get("Content-Type"), "application/json",
				"Expected json content type")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(httpbinResponse)
	res, err := api.Res("put").Put(payload)
	assert.NoError(t, err, "Error Getting Data from httpbin API")
	r.Json, err = res.Response.Get("json").Map()
	assert.NoError(t, err, "Could not convert responses to map")
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestCanDeleteResource(t *testing.T) {
	testMux.HandleFunc("/delete",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "DELETE", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/delete", "unexpected Path")
			fmt.Fprintln(rw, readJson("_tests/delete_response.json"))
		})

	api := Api(testSrv.URL)
	r := new(httpbinResponse)
	res, err := api.Res("delete").Delete(nil)
	assert.NoError(t, err, "Error Getting Data from httpbin API")
	r.Url, err = res.Response.Get("url").String()
	assert.NoError(t, err, "Could not convert responses to map")
	assert.Equal(t, r.Url, "https://httpbin.org/delete", "Url must match")
}

func TestPathSuffix(t *testing.T) {
	testMux.HandleFunc("/item/32.json",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "GET", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/item/32.json", "unexpected Path")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL, ".json")
	r := new(httpbinResponse)
	res, err := api.Res("item").Id(32).Get(nil)
	assert.NoError(t, err, "Error requesting item")

	config := &mapstructure.DecoderConfig{Result: r}
	dec, err := mapstructure.NewDecoder(config)
	assert.NoError(t, err, "Error constructing mapstructure")

	err = dec.Decode(res.Response.Interface())
	assert.NoError(t, err, "Error decoding response")

	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestPathSuffixWithQueryParam(t *testing.T) {
	testMux.HandleFunc("/item/42.json",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "GET", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/item/42.json", "unexpected Path")
			assert.Equal(t, req.URL.Query().Get("param"), "test", "unexpected QueryParam")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL, ".json")
	r := new(httpbinResponse)
	res, err := api.Res("item").Id(42).Get(map[string]string{"param": "test"})
	assert.NoError(t, err, "Error requesting item")

	config := &mapstructure.DecoderConfig{Result: r}
	dec, err := mapstructure.NewDecoder(config)
	assert.NoError(t, err, "Error constructing mapstructure")

	err = dec.Decode(res.Response.Interface())
	assert.NoError(t, err, "Error decoding response")
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestResourceId(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Res("users").Id("test").Url, "users/test",
		"Url should match")
	assert.Equal(t, api.Res("users").Id(123).Res("items").Id(111).Url,
		"users/123/items/111", "Multilevel Url should match")
	assert.Equal(t, api.Res("users").Id(int64(9223372036854775807)).Url, "users/9223372036854775807",
		"int64 id should work")
}

func TestDoNotDecodeBodyOnErr(t *testing.T) {
	tests := []int{
		400, 401, 500, 501,
	}

	actualData := strings.TrimSpace(readJson("_tests/error.json"))

	// will be changed later
	code := 0
	testMux.HandleFunc("/error",
		func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(code)
			fmt.Fprintln(rw, actualData)
		})

	api := Api(testSrv.URL)

	for _, code = range tests {
		resp := make(map[string]interface{})
		r, err := api.Id("error").Get(nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, map[string]interface{}{}, resp,
			fmt.Sprintf("response should be unparsed: %d", code))

		respData, err := ioutil.ReadAll(r.Raw.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, actualData, strings.TrimSpace(string(respData)),
			fmt.Sprintf("response body is not accessible: %d", code))
	}
}

func readJson(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
