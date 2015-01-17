package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {

	setup()
	defer teardown()

	user := newUserService(client)

	mux.HandleFunc("/1.0/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Nil(t, r.ParseForm())

		assert.Equal(t, r.PostForm.Get("auth"), "test@mail.com")
		assert.Equal(t, r.PostForm.Get("pass"), "94406d8b3a3876308552d168e56a42f9")
		fmt.Fprint(w, `[200,{"cookie":"1000000000:efcb5ef3efec97aa50c33e1efb183e223633a3bf","user":{"id":"1000000000","name":"johndoe","email":"john@example.com","api_key":"d09272c7412aba77d1b06795bf9d8f701ee0171e","pro":"0000-00-00T00:00:00.000Z","webspace": 0.523,"traffic":{"current":0.532,"increase":0.532,"last":0.532,"max":0.532 },"balance":0,"settings":{"rewrite_behaviour":1,"ddl":0},"external_id":"EXTERNAL_OAUTH_PROVIDER_ID","ftp_username":"johndoe","partner":"3","partner_last":1392728258},"session":"cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a"}]`)
	})

	resp, err := user.Login("test@mail.com", "1234")
	assert.Nil(t, err)
	assert.Equal(t, resp.Cookie, "1000000000:efcb5ef3efec97aa50c33e1efb183e223633a3bf")
	assert.Equal(t, resp.User.Name, "johndoe")
	assert.Equal(t, resp.Session, "cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a")
	assert.Equal(t, user.session, "cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a")
}
