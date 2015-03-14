package store

import (
	"errors"
	"time"

	"github.com/cryptix/trakting/types"
	"github.com/jmoiron/modl"
)

func init() {
	DB.AddTable(types.Track{}).SetKeys(true, "id")
	createSql = append(createSql, `alter table track alter added set default now()`)
}

type TrackStore struct {
	dbh modl.SqlExecutor
}

var _ types.Tracker = (*TrackStore)(nil)

func NewTrackStore() (*TrackStore, error) {
	if DBH == nil {
		return nil, errors.New("connect db first")
	}

	return &TrackStore{DBH}, nil
}

func (t *TrackStore) Add(tr types.Track) error {
	tr.Added = time.Now()
	return t.dbh.Insert(&tr)
}

func (t *TrackStore) Get(boomID string) (types.Track, error) {
	var track types.Track
	err := t.dbh.SelectOne(&track, `SELECT * FROM "track" WHERE BoomID = $1`, boomID)
	return track, err
}

func (t *TrackStore) All() ([]types.Track, error) {
	var tracks []types.Track
	err := t.dbh.Select(&tracks, `SELECT * FROM "track" ORDER BY added DESC`)
	return tracks, err
}

func (t *TrackStore) ByUserName(name string) ([]types.Track, error) {
	var tracks []types.Track
	err := t.dbh.Select(&tracks, `SELECT * FROM "track" WHERE by = $1 ORDER BY added DESC`, name)
	return tracks, err
}
