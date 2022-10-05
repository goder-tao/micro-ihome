package fastdfs

import (
	"github.com/tedcy/fdfs_client"
	"sync"
)

var cli *fdfs_client.Client = nil
var mu sync.Mutex

func Init() error {
	c, err := fdfs_client.NewClientWithConfig("/home/tao/Data/Software/project/go/project/IHome/configs/fdfs_client.conf")
	if err != nil {
		return err
	}
	cli = c
	return nil
}

func Get() (*fdfs_client.Client, error) {
	if cli == nil {
		err := Init()
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}
