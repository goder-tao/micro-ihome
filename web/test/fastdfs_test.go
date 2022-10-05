package test

import (
	"fmt"
	"github.com/tedcy/fdfs_client"
	"os"
	"testing"
)

func TestFdfs(t *testing.T) {
	fmt.Println(os.Getwd())
	cli, err := fdfs_client.NewClientWithConfig("../config/fdfs_client.conf")
	if err != nil {
		t.Fatalf("fdfs new client: %s", err.Error())
	}
	v := "hello"
	fileID, err := cli.UploadByBuffer([]byte(v), "txt")
	if err != nil {
		t.Fatalf("fdfs client upload buffer: %s", err.Error())
	}

	// 取数据是否相等
	buf := make([]byte, len(v))
	if err := cli.DownloadToAllocatedBuffer(fileID, buf, 0, int64(len(v))); err != nil {
		t.Fatalf("fdfs client download buffer: %s", err.Error())
	}
	if v != string(buf) {
		t.Fatalf("invalid store data, orignal: %s, store: %s", v, string(buf))
	}

	if err := cli.DeleteFile(fileID); err != nil {
		t.Fatalf("fdfs client delete: %s", err.Error())
	}
}
