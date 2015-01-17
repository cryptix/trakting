package goBoom

import "strings"

type InformationService struct {
	c *Client
}

func newInformationService(c *Client) *InformationService {
	i := &InformationService{}
	if c == nil {
		i.c = NewClient(nil)
	} else {
		i.c = c
	}

	return i
}

func (i InformationService) Info(ids ...string) ([]ItemStat, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"items": strings.Join(ids, ","),
	}

	resp, err := i.c.api.Res("info").Get(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	var infoResp []ItemStat
	if err = decodeInto(&infoResp, arr[1]); err != nil {
		return nil, err
	}

	return infoResp, nil
}

type ItemSize struct {
	Num  int64 `json:"num"`
	Size int64 `json:"size"`
}

func (i InformationService) Du() (map[string]ItemSize, error) {

	params := map[string]string{
		"token": i.c.User.session,
	}

	resp, err := i.c.api.Res("du").Get(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	duResp := make(map[string]ItemSize)
	if err = decodeInto(&duResp, arr[1]); err != nil {
		return nil, err
	}

	return duResp, nil
}

type LsInfo struct {
	Pwd   ItemStat
	Items []ItemStat
}

func (i InformationService) Ls(item string) (*LsInfo, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"item":  item,
	}

	resp, err := i.c.api.Res("ls").Get(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, err
	}

	var lsResp LsInfo
	if err = decodeInto(&lsResp.Pwd, arr[1]); err != nil {
		return nil, err
	}

	if err = decodeInto(&lsResp.Items, arr[2]); err != nil {
		return nil, err
	}

	return &lsResp, nil
}

func (i *InformationService) Tree(rev string) ([]ItemStat, map[string]string, error) {
	params := map[string]string{
		"token": i.c.User.session,
	}

	if rev != "" {
		params["revision"] = rev
	}

	resp, err := i.c.api.Res("tree").Get(params)
	arr, err := processResponse(resp, err)
	if err != nil {
		return nil, nil, err
	}

	var items []ItemStat
	err = decodeInto(&items, arr[1])
	if err != nil {
		return nil, nil, err
	}

	var revs map[string]string
	err = decodeInto(&revs, arr[2])
	if err != nil {
		return nil, nil, err
	}

	return items, revs, nil
}
