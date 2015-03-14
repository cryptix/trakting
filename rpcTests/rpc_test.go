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
)

var (
	lis         net.Listener
	testServer  *rpc.Server
	testClient  *rpcClient.Client
	fakeTracker = new(FakeTracker)
	fakeUserer  = new(FakeUserer)
)

func TestMain(m *testing.M) {
	testServer = rpc.NewServer()

	ts, e := rpcServer.NewTrackService(types.User{}, fakeTracker)
	if e != nil {
		log.Fatal("NewTrackService error:", e)
	}
	testServer.RegisterName("TrackService", ts)

	us, e := rpcServer.NewUserService(types.User{}, fakeUserer)
	if e != nil {
		log.Fatal("NewUserService error:", e)
	}
	testServer.RegisterName("UserService", us)

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

	tcClient, e := rpcClient.NewTracksClient(rpcc)
	if e != nil {
		log.Fatal("rpcClient.NewTracksClient:", e)
	}

	tuClient, e := rpcClient.NewUsersClient(rpcc)
	if e != nil {
		log.Fatal("rpcClient.NewUsersClient:", e)
	}

	testClient = rpcClient.New(tcClient, tuClient)

	ret := m.Run()

	if e := lis.Close(); e != nil {
		log.Fatal(e)
	}
	testServer = nil

	os.Exit(ret)
}
