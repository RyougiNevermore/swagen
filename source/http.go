package source

import (
	"fmt"
	"github.com/pharosnet/swagen/zlog"
	"io/ioutil"
	"net/http"
)

func fetchByHttp(path string) (b []byte, err error) {

	resp, respErr := http.Get(path)

	if respErr != nil {
		zlog.Log().With("status", "failed").Debugf("read http %s", path)
		err = respErr
		return
	}

	defer func(response *http.Response) {
		_ = response.Body.Close()
	}(resp)

	if resp.StatusCode != 200 {
		zlog.Log().With("status", "failed").Debugf("read http %s", path)
		err = fmt.Errorf("read http %s failed, status is %d", path, resp.StatusCode)
		return
	}

	if resp.ContentLength == 0 {
		zlog.Log().With("status", "failed").Debugf("read http %s", path)
		err = fmt.Errorf("read http %s failed, content length is zero", path)
		return
	}

	d, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		zlog.Log().With("status", "failed").Debugf("read http %s", path)
		err = respErr
		return
	}

	b = d

	zlog.Log().With("status", "succeed").Debugf("read http %s", path)

	return
}
