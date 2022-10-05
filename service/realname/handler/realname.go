package handler

import (
	"context"

	realname "ihome/service/realname/proto"
)

var fakeData = map[string]string{
	"张三": "123456789012345000",
	"李四": "123456789012345001",
	"王五": "123456789012345002",
}

type Realname struct{}

// Return a new handler
func New() *Realname {
	return &Realname{}
}

func (e *Realname) AuthRealName(ctx context.Context, req *realname.RealNameRequest, rsp *realname.RealNameResponse) error {
	idn, ok := fakeData[req.Realname]
	if !ok {
		rsp.Errmsg = "no such people"
		return nil
	}
	if idn != req.IdNumber {
		rsp.Errmsg = "mismatch name and id number"
		return nil
	}
	return nil
}
