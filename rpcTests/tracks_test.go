package rpcTests

import (
	"testing"

	"github.com/cryptix/trakting/types"
	"github.com/stretchr/testify/require"
)

func TestTrackAdd(t *testing.T) {
	want := types.Track{}
	want.ID = 123
	want.Name = "testTrack"
	e := testClient.Tracks.Add(want)
	require.Nil(t, e)
	require.Equal(t, 1, fakeTracker.AddCallCount())
	require.Equal(t, want, fakeTracker.AddArgsForCall(0))
}

func TestTrackAdd_emptyName(t *testing.T) {
	want := types.Track{}
	want.Name = ""
	e := testClient.Tracks.Add(want)
	require.NotNil(t, e)
	require.Equal(t, types.ErrEmptyTrackName, e)
}

func TestTrackAll(t *testing.T) {
	want := []types.Track{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	fakeTracker.AllReturns(want, nil)
	got, e := testClient.Tracks.All()
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, fakeTracker.AllCallCount())
}

func TestTrackByUserName(t *testing.T) {
	want := []types.Track{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	fakeTracker.ByUserNameReturns(want, nil)
	got, e := testClient.Tracks.ByUserName("herb")
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, fakeTracker.ByUserNameCallCount())
	require.Equal(t, "herb", fakeTracker.ByUserNameArgsForCall(0))
}

func TestTrackGet(t *testing.T) {
	var want = types.Track{ID: 1, Name: "123"}
	fakeTracker.GetReturns(want, nil)
	got, e := testClient.Tracks.Get("555")
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, fakeTracker.GetCallCount())
	require.Equal(t, "555", fakeTracker.GetArgsForCall(0))
}
