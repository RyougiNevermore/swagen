package source

import (
	"github.com/pharosnet/swagen/zlog"
	"io/ioutil"
)

func fetchByFile(path string) (b []byte, err error) {

	b, err = ioutil.ReadFile(path)

	if err != nil {
		zlog.Log().With("status", "failed").Debugf("read file %s", path)
	} else {
		zlog.Log().With("status", "succeed").Debugf("read file %s", path)
	}

	return
}
