package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

func init() {
	gob.Register(Track{})
}

const trackBucket = "tracks"

type Track struct {
	By     string
	Name   string
	BoomID string
}

func (t Track) String() string {
	return fmt.Sprintf("%q (by %s)", t.Name, t.By)
}

type TrackStore struct {
	db *bolt.DB
}

func NewTrackStore(db *bolt.DB) (*TrackStore, error) {
	if db == nil {
		return nil, errors.New("Init() first")
	}

	return &TrackStore{db}, nil
}

func (t *TrackStore) Add(tr Track) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		var (
			err error
			buf bytes.Buffer
		)

		err = gob.NewEncoder(&buf).Encode(tr)
		if err != nil {
			return err
		}

		tx.CreateBucket([]byte(trackBucket))
		b := tx.Bucket([]byte(trackBucket))
		b.Put([]byte(tr.BoomID), buf.Bytes())

		return nil
	})
}

func (t *TrackStore) Get(boomID string) (Track, error) {
	var track Track
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(trackBucket))
		if b == nil {
			return errors.New("no tracks at all")
		}

		v := b.Get([]byte(boomID))
		if v == nil {
			return errors.New("no such track at all")
		}

		return gob.NewDecoder(bytes.NewReader(v)).Decode(&track)
	})
	return track, err
}

func (t *TrackStore) All() ([]Track, error) {
	var tracks []Track
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(trackBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tr Track
			err := gob.NewDecoder(bytes.NewReader(v)).Decode(&tr)
			if err != nil {
				return err
			}
			tracks = append(tracks, tr)

		}

		return nil
	})

	return tracks, err
}
