package handler

import (
	"context"
	"fmt"
	avatar "ihome/service/avatar/proto"
	"ihome/web/dao/fastdfs"
)

type Avatar struct{}

// Return a new handler
func New() *Avatar {
	return &Avatar{}
}

func (e *Avatar) SaveBuffer(ctx context.Context, req *avatar.AvatarRequest, rsp *avatar.AvatarResponse) error {
	fmt.Println("remote success")
	fdfs, err := fastdfs.Get()
	if err != nil {
		rsp.Err = err.Error()
		return err
	}
	fmt.Println("before upload, exeName: ", req.ExeName)
	fileID, err := fdfs.UploadByBuffer(req.Buf, req.ExeName)
	fmt.Println("after upload, fileID:", fileID)
	if err != nil {
		rsp.Err = err.Error()
		return err
	}
	rsp.FileID = fileID
	return nil
}
