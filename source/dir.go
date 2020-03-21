package source

import (
	"github.com/pharosnet/swagen/zlog"
	"io/ioutil"
	"strings"
)

func fetchByDir(path string) (b [][]byte, err error) {

	files, dirErr := ioutil.ReadDir(path)
	if dirErr != nil {
		zlog.Log().With("status", "failed").Debugf("read dir %s", path)
		err = dirErr
		return
	}
	zlog.Log().With("status", "succeed").Debugf("read dir %s", path)

	b = make([][]byte, 0, 1)

	for _, f := range files {
		if f.IsDir() {
			b0, b0Err := fetchByDir(strings.Join([]string{path, f.Name()}, "/"))
			if b0Err != nil {
				err = b0Err
				return
			}
			if b0 == nil || len(b0) == 0 {
				continue
			}
			for _, b00 := range b0 {
				b = append(b, b00)
			}
		}
		if strings.Contains(strings.ToLower(f.Name()), "json") {
			b0, b0Err := fetchByFile(strings.Join([]string{path, f.Name()}, "/"))
			if b0Err != nil {
				err = b0Err
				return
			}
			b = append(b, b0)
		} else if strings.Contains(strings.ToLower(f.Name()), "yaml") || strings.Contains(strings.ToLower(f.Name()), "yml") {
			b0, b0Err := fetchByFile(strings.Join([]string{path, f.Name()}, "/"))
			if b0Err != nil {
				err = b0Err
				return
			}
			b = append(b, b0)
		}
	}

	return
}
