package rpcTests

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"testing"

	"github.com/cryptix/trakting/rpcClient"
	"github.com/cryptix/trakting/rpcServer"
	"github.com/cryptix/trakting/types"
	"github.com/stretchr/testify/require"
)

var (
	lis        net.Listener
	testServer *rpc.Server
	ft         = new(FakeTracker)
	testClient types.Tracker
)

func TestMain(m *testing.M) {
	// create server with fake db
	ts, e := rpcServer.NewTrackService(types.User{}, ft)
	if e != nil {
		log.Fatal("NewTrackService error:", e)
	}
	testServer = rpc.NewServer()

	testServer.RegisterName("TrackService", ts)
	lis, e = net.Listen("tcp", ":0")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("listening on", lis.Addr())
	go http.Serve(lis, testServer)

	// dial client to it and create test client
	rpcc, e := rpc.DialHTTP("tcp", lis.Addr().String())
	if e != nil {
		log.Fatal("rpc.DialHTTP:", e)
	}

	testClient, e = rpcClient.NewTracksClient(rpcc)
	if e != nil {
		log.Fatal("rpcClient.NewTracksClient:", e)
	}

	ret := m.Run()

	if e := lis.Close(); e != nil {
		log.Fatal(e)
	}
	testServer = nil

	os.Exit(ret)
}

// Add(Track) error
func TestAdd(t *testing.T) {
	want := types.Track{}
	want.ID = 123
	want.Name = "testTrack"
	e := testClient.Add(want)
	require.Nil(t, e)
	require.Equal(t, 1, ft.AddCallCount())
	require.Equal(t, want, ft.AddArgsForCall(0))
}

func TestAdd_emptyName(t *testing.T) {
	want := types.Track{}
	want.Name = ""
	e := testClient.Add(want)
	require.NotNil(t, e)
	require.Equal(t, types.ErrEmptyTrackName, e)
}

// All() ([]Track, error)
func TestAll(t *testing.T) {
	want := []types.Track{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	ft.AllReturns(want, nil)
	got, e := testClient.All()
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, ft.AllCallCount())
}

// ByUserName(name string) ([]Track, error)
func TestByUserName(t *testing.T) {
	want := []types.Track{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	ft.ByUserNameReturns(want, nil)
	got, e := testClient.ByUserName("herb")
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, ft.ByUserNameCallCount())
	require.Equal(t, "herb", ft.ByUserNameArgsForCall(0))
}

// Get(id string) (Track, error)
func TestGet(t *testing.T) {
	var want = types.Track{ID: 1, Name: "123"}
	ft.GetReturns(want, nil)
	got, e := testClient.Get("555")
	require.Nil(t, e)
	require.Equal(t, want, got)
	require.Equal(t, 1, ft.GetCallCount())
	require.Equal(t, "555", ft.GetArgsForCall(0))
}
