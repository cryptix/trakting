package rpcServer

import "github.com/cryptix/trakting/types"

type TrackService struct {
	user types.User
	db   types.Tracker
}

func NewTrackService(u types.User, db types.Tracker) (*TrackService, error) {
	return &TrackService{
		user: u,
		db:   db,
	}, nil
}

func (ts *TrackService) Add(args *types.Track, _ *struct{}) error {
	return ts.db.Add(*args)
}

func (ts *TrackService) All(args *string, reply *[]types.Track) (err error) {
	*reply, err = ts.db.All()
	return
}

func (ts *TrackService) ByUserName(args *string, reply *[]types.Track) (err error) {
	*reply, err = ts.db.ByUserName(*args)
	return
}

func (ts *TrackService) Get(args *string, reply *types.Track) (err error) {
	*reply, err = ts.db.Get(*args)
	return
}
