package source

import (
	"fmt"
	"github.com/go-openapi/spec"
)

func Fetch(kind string, path string) (sgs []*spec.Swagger, err error) {

	switch kind {
	case "file":
		b, fetchErr := fetchByFile(path)
		if fetchErr != nil {
			err = fetchErr
			return
		}
		bb := [][]byte{b}
		sgs, err = parse(bb)

	case "dir":
		b, fetchErr := fetchByDir(path)
		if fetchErr != nil {
			err = fetchErr
			return
		}
		sgs, err = parse(b)

	case "http":
		b, fetchErr := fetchByHttp(path)
		if fetchErr != nil {
			err = fetchErr
			return
		}
		bb := [][]byte{b}
		sgs, err = parse(bb)

	default:
		err = fmt.Errorf("read source swagger file failed, bad source_type (%s)", kind)
	}

	return
}

func parse(bb [][]byte) (sgs []*spec.Swagger, err error) {
	if bb == nil || len(bb) == 0 {
		return
	}

	sgs = make([]*spec.Swagger, 0, 1)
	for _, b := range bb {
		if b == nil || len(b) == 0 {
			continue
		}

		sg := &spec.Swagger{}

		parseErr := sg.UnmarshalJSON(b)
		if parseErr != nil {
			err = parseErr
			return
		}
		sgs = append(sgs, sg)
	}

	return
}
