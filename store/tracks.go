package store

import (
	"errors"
	"fmt"

	"github.com/jmoiron/modl"
)

func init() {
	DB.AddTable(Track{}).SetKeys(true, "id")
}

type Track struct {
	ID     int64
	By     string
	Name   string
	BoomID string
}

func (t Track) String() string {
	return fmt.Sprintf("%q (by %s)", t.Name, t.By)
}

type TrackStore struct {
	dbh modl.SqlExecutor
}

func NewTrackStore() (*TrackStore, error) {
	if DBH == nil {
		return nil, errors.New("connect db first")
	}

	return &TrackStore{DBH}, nil
}

func (t *TrackStore) Add(tr Track) error {
	return t.dbh.Insert(&tr)
}

func (t *TrackStore) Get(boomID string) (Track, error) {
	var track Track
	err := t.dbh.SelectOne(&track, `SELECT * FROM "track" WHERE BoomID = $1`, boomID)
	return track, err
}

func (t *TrackStore) All() ([]Track, error) {
	var tracks []Track
	err := t.dbh.Select(&tracks, `SELECT * FROM "track"`)
	return tracks, err
}
