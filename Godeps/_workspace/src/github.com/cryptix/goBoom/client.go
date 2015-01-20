// package goBoom implements clients for the oBoom API. Docs: https://www.oboom.com/api
package goBoom

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/cryptix/gocrayons"
	"github.com/mitchellh/mapstructure"
)

const (
	libraryVersion = "1.0"
	defaultBaseURL = "https://api.oboom.com/1.0"
	userAgent      = "goBoom/" + libraryVersion

	defaultAccept    = "application/json"
	defaultMediaType = "application/octet-stream"
)

var (
	ErrUnknwonFourResponseType = errors.New("Can only handle three responses for InformationService.Ls()")
)

type ErrStatusCodeMissmatch struct{ Http, Api int }

func (e ErrStatusCodeMissmatch) Error() string {
	return fmt.Sprintf("ErrStatusCodeMissmatch: http[%d] != api[%d]", e.Http, e.Api)
}

// A Client manages communication with the oBoom API.
type Client struct {
	// new rest api client
	api *gocrayons.Resource

	// HTTP client used to communicate with the API.
	c *http.Client

	// Base URL for API requests.  baseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// User agent used when communicating with the REST API.
	userAgent string

	User *UserService
	Info *InformationService
	FS   *FilesystemService
}

// NewClient returns a new REST API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	httpClient.Jar = jar
	client := &Client{c: httpClient, baseURL: baseURL, userAgent: userAgent}

	client.api = gocrayons.Api(defaultBaseURL)
	client.api.SetClient(httpClient)

	client.User = newUserService(client)
	client.Info = newInformationService(client)
	client.FS = newFilesystemService(client)

	return client
}

func processResponse(resp *gocrayons.Resource, err error) ([]interface{}, error) {
	if err != nil {
		fmt.Printf("%+v\n", resp.Raw)
		return nil, err
	}

	if err := checkResponse(resp.Raw); err != nil {
		return nil, err
	}

	arr, err := resp.Response.Array()
	if err != nil {
		return nil, err
	}

	if len(arr) < 1 {
		return nil, ErrorResponse{resp.Raw, "Illegal oBoom response"}
	}

	statusCode, err := resp.Response.GetIndex(0).Int()
	if err != nil {
		return nil, err
	}

	if statusCode != resp.Raw.StatusCode {
		return nil, ErrorResponse{resp.Raw, fmt.Sprintf("StatusCode missmatch. %d vs %d\n%+v",
			statusCode,
			resp.Raw.StatusCode,
			resp.Response)}
	}

	return arr, nil
}

func decodeInto(t interface{}, input interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           t,
	}

	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return errors.New("NewDecoder Error:" + err.Error())
	}

	if err = dec.Decode(input); err != nil {
		return errors.New("Decode Error:" + err.Error())
	}

	return nil
}

/*
An ErrorResponse reports one or more errors caused by an API request.
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Body     string
}

func (r ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %q",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Body)
}

// checkResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = string(data)
	}
	return errorResponse
}
