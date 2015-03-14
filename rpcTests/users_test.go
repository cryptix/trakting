package rpcTests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAdd(t *testing.T) {
	e := testClient.Users.Add("testUser", "testPw", 1)
	require.Nil(t, e)
	require.Equal(t, 1, fakeUserer.AddCallCount())
	gotUser, gotPW, gotLvl := fakeUserer.AddArgsForCall(0)
	require.Equal(t, "testUser", gotUser)
	require.Equal(t, "testPw", gotPW)
	require.Equal(t, 1, gotLvl)
}

func TestUserChangePassw(t *testing.T) {
	e := testClient.Users.ChangePassword(0, "testPW")
	require.Nil(t, e)
	require.Equal(t, 1, fakeUserer.ChangePasswordCallCount())
	id, passw := fakeUserer.ChangePasswordArgsForCall(0)
	require.Equal(t, 0, id)
	require.Equal(t, "testPW", passw)
}
