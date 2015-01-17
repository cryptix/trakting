package goBoom

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSession = "testSession"

func TestInformationService_Info(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	want := []ItemStat{
		{
			Iname: "trash",
			Root:  "1C",
			State: "online",
			User:  298814,
			Type:  "folder",
			ID:    "1C",
		},
		{
			Iname: "public",
			ID:    "1",
			Root:  "1",
			State: "online",
			User:  298814,
			Type:  "folder",
		},
		{
			Ctime: "2014-06-21 23:23:46.615535",
			Mtime: "2014-10-07 16:44:59.208856",
			Root:  "1",
			Iname: "pdfs",
			ID:    "99QJ0C6Y",
			State: "online",
			User:  298814,
			Type:  "folder",
		},
	}

	mux.HandleFunc("/1.0/info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		assert.Equal(t, r.URL.Query().Get("items"), "a,b,c")
		cpJson(t, w, "_tests/info.json")
	})

	resp, err := info.Info("a", "b", "c")
	assert.Nil(t, err)

	assert.IsType(t, resp, []ItemStat{})
	assert.Len(t, resp, len(want), "resp has incorrect length")
	for i := range want {
		assert.Equal(t, want[i], resp[i], "resp[%d] differs", i)
	}
}

func TestInformationService_Du(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	mux.HandleFunc("/1.0/du", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		cpJson(t, w, "_tests/du.json")
	})

	resp, err := info.Du()
	assert.Nil(t, err)

	dummyMap := make(map[string]ItemSize)
	assert.IsType(t, resp, dummyMap)
	assert.Len(t, resp, 3)

	assert.Equal(t, resp["total"], ItemSize{1, 2893557})
	assert.Equal(t, resp["1"], ItemSize{1, 2893557})
	assert.Equal(t, resp["1C"], ItemSize{0, 0})
}

func TestInformationService_Ls(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	wantPwd := ItemStat{
		ID:    "99QJ0C6Y",
		Root:  "1",
		Iname: "pdfs",
		State: "online",
		Type:  "folder",
		User:  298814,
		Ctime: "2014-06-21 23:23:46.615535",
		Mtime: "2014-10-07 16:44:59.208856",
	}
	wantItems := []ItemStat{
		{
			Isize:     2893557,
			Ctime:     "2014-06-16 09:49:07.808925",
			Parent:    "99QJ0C6Y",
			Type:      "file",
			Downloads: 1,
			Mtime:     "2014-06-21 23:23:49.391054",
			State:     "online",
			Mime:      "application/pdf",
			User:      298814,
			Owner:     true,
			Atime:     "2014-06-16 09:49:44.725223",
			Root:      "1",
			ID:        "3TRL28BM",
			Iname:     "gobook.pdf",
		},
		{
			Isize:     368005,
			Ctime:     "2014-06-21 23:23:12.143749",
			Parent:    "99QJ0C6Y",
			Type:      "file",
			Downloads: 3,
			Mtime:     "2014-12-25 01:42:58.747235",
			State:     "online",
			Mime:      "application/pdf",
			User:      298814,
			Owner:     true,
			Atime:     "2014-12-25 01:42:58.747235",
			Root:      "1",
			ID:        "DUT1V03Z",
			Iname:     "ds160.pdf",
		},
		{
			Isize:     1062822,
			Ctime:     "2014-06-21 23:23:32.754727",
			Parent:    "99QJ0C6Y",
			Type:      "file",
			Downloads: 1,
			Mtime:     "2014-08-08 21:25:07.721670",
			State:     "online",
			Mime:      "application/pdf",
			User:      298814,
			Owner:     true,
			Atime:     "2014-08-08 21:25:07.721670",
			Root:      "1",
			ID:        "L8PHHR89",
			Iname:     "calculus-indexed.pdf",
		},
	}

	mux.HandleFunc("/1.0/ls", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		assert.Equal(t, r.URL.Query().Get("item"), "pdfs")
		cpJson(t, w, "_tests/ls.json")
	})

	resp, err := info.Ls("pdfs")
	assert.Nil(t, err)

	assert.Equal(t, wantPwd, resp.Pwd)
	assert.Len(t, resp.Items, len(wantItems), "resp has incorrect length")
	for i := range wantItems {
		assert.Equal(t, wantItems[i], resp.Items[i], "resp[%d] differs", i)
	}
}

func TestInformationService_Tree(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	wantTree := []ItemStat{
		{
			Iname: "public",
			ID:    "1",
			Root:  "1",
			State: "online",
			User:  298814,
			Type:  "folder",
		},
		{
			Ctime:     "2014-12-25 12:54:56.734294",
			Parent:    "1C",
			Root:      "1C",
			Downloads: 0,
			State:     "online",
			User:      298814,
			Mtime:     "2014-12-25 12:55:07.626095",
			Type:      "folder",
			ID:        "TCNCL2X2",
			Iname:     "test1",
		},
		{
			Ctime:     "2014-12-25 12:54:01.294567",
			Parent:    "1C",
			Root:      "1C",
			Downloads: 1,
			Iname:     "File.txt",
			State:     "online",
			Mime:      "text/plain",
			User:      298814,
			Mtime:     "2014-12-25 12:54:17.324197",
			Atime:     "2014-12-25 12:54:09.696367",
			Type:      "file",
			ID:        "MB06AK12",
			Isize:     1,
			Owner:     true,
		},
	}
	wantRevs := map[string]string{
		"298814": "j3NWVmSdLlB",
	}
	mux.HandleFunc("/1.0/tree", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		cpJson(t, w, "_tests/tree.json")
	})

	tree, revs, err := info.Tree("")
	assert.Nil(t, err)

	assert.Equal(t, wantRevs, revs)
	assert.Len(t, tree, len(wantTree), "resp has incorrect length")
	for i := range wantTree {
		assert.Equal(t, wantTree[i], tree[i], "resp[%d] differs", i)
	}
}

func cpJson(t *testing.T, w io.Writer, path string) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cpJson Open failed: %s", err)
	}

	_, err = io.Copy(w, f)
	if err != nil {
		t.Fatalf("cpJson Copy failed: %s", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("cpJson Close failed: %s", err)
	}
}
