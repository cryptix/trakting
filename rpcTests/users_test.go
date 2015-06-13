package rpcTests

import (
	"testing"

	"github.com/cryptix/trakting/types"
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
	require.Equal(t, int64(0), id)
	require.Equal(t, "testPW", passw)
}

func TestUserList(t *testing.T) {
	want := []types.User{
		{ID: 1, Name: "Hans"},
		{ID: 2, Name: "Franz"},
		{ID: 3},
	}
	fakeUserer.ListReturns(want, nil)
	got, e := testClient.Users.List()
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, fakeUserer.ListCallCount())
}
