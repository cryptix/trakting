package goBoom

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type UserService struct {
	c *Client

	session string
}

func newUserService(c *Client) *UserService {
	u := &UserService{}
	if c == nil {
		u.c = NewClient(nil)
	} else {
		u.c = c
	}

	return u
}

type loginResponse struct {
	Cookie  string `json:"cookie"`
	Session string `json:"session"`
	User    struct {
		ApiKey      string      `json:"api_key"`
		Balance     interface{} `json:"balance"`
		Email       string      `json:"email"`
		ExternalID  string      `json:"external_id"`
		FtpUsername string      `json:"ftp_username"`
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Partner     string      `json:"partner"`
		PartnerLast interface{} `json:"partner_last"`
		Pro         string      `json:"pro"`
		Settings    struct {
			Ddl              float64 `json:"ddl"`
			RewriteBehaviour float64 `json:"rewrite_behaviour"`
		} `json:"settings"`
		Traffic struct {
			Current  float64 `json:"current"`
			Increase float64 `json:"increase"`
			Last     float64 `json:"last"`
			Max      float64 `json:"max"`
		} `json:"traffic"`
		Webspace float64 `json:"webspace"`
	} `json:"user"`
}

// Login sends a login request to the service with name and passw as credentials
func (u *UserService) Login(name, passw string) (*loginResponse, error) {

	derived := pbkdf2.Key([]byte(passw), []byte(reverse(passw)), 1000, 16, sha1.New)

	reqParams := url.Values{
		"auth": []string{name},
		"pass": []string{fmt.Sprintf("%x", derived)},
	}

	oldHost := u.c.api.Api.BaseUrl.Host
	u.c.api.Api.BaseUrl.Host = strings.Replace(oldHost, "api.oboom.com", "www.oboom.com", 1)

	res := u.c.api.Res("login")
	res.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	res.Payload = strings.NewReader(reqParams.Encode())
	resp, err := res.FormPost(nil)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	u.c.api.Api.BaseUrl.Host = oldHost

	var liResp loginResponse
	if err = decodeInto(&liResp, arr[1]); err != nil {
		return nil, err
	}

	u.session = liResp.Session

	return &liResp, nil
}
