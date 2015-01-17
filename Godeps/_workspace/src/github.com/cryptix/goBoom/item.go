package goBoom

import (
	"fmt"
	"os"
	"time"
)

type ItemStat struct {
	Atime     string `mapstructure:"atime"`
	Ctime     string `mapstructure:"ctime"`
	Downloads int64  `mapstructure:"downloads"`
	ID        string `mapstructure:"id"`
	Mtime     string `mapstructure:"mtime"`
	Iname     string `mapstructure:"name"`
	Parent    string `mapstructure:"parent"`
	Root      string `mapstructure:"root"`
	State     string `mapstructure:"state"`
	Type      string `mapstructure:"type"`
	User      int64  `mapstructure:"user"`
	Isize     int64  `mapstructure:"size"`
	DDL       bool   `mapstructure:"ddl"`
	Mime      string `mapstructure:"mime"`
	Owner     bool   `mapstructure:"owner"`
}

func (i ItemStat) IsDir() bool {
	return i.Type == "folder"
}

func (i ItemStat) ModTime() time.Time {
	const format = `2006-01-02 15:04:05.000000`
	var (
		t   = time.Now()
		err error
	)
	switch {
	case i.Mtime != "":
		t, err = time.Parse(format, i.Mtime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	case i.Ctime != "":
		t, err = time.Parse(format, i.Mtime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	case i.Atime != "":
		t, err = time.Parse(format, i.Atime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	default:
		return t
	}

	return t
}

func (i ItemStat) Mode() os.FileMode {
	return os.ModePerm
}

func (i ItemStat) Name() string {
	return i.Iname
}
func (i ItemStat) Size() int64 {
	return i.Isize
}

func (i ItemStat) Sys() interface{} {
	return nil
}
