package py

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"github.com/pharosnet/swagen/zlog"
	"io/ioutil"
	"strings"
	"time"
)

func Generate(sgs []*spec.Swagger, output string) (err error) {

	for _, sg := range sgs {
		title := sg.Info.Title
		if title == "" {
			err = errors.New("generate failed, title in swagger doc is not defined")
			return
		}
		zlog.Log().With("title", sg.Info.Title).Debug("begin to generate")

		outputFile := strings.Join([]string{output, fmt.Sprintf("%s.py", title)}, "/")

		p, genErr := generate0(sg)

		if genErr != nil {
			err = fmt.Errorf("generate failed, generate code failed, %v", genErr)
			return
		}

		wErr := ioutil.WriteFile(outputFile, p, 0666)
		if wErr != nil {
			err = fmt.Errorf("generate failed, write to code file failed, %v", wErr)
			return
		}

	}

	return
}


func generate0(sg *spec.Swagger) (p []byte, err error) {

	pkg := strings.TrimSpace(strings.ToLower(strcase.ToSnake(sg.Info.Title)))
	version := strings.TrimSpace(strings.ToLower(sg.Info.Version))

	bb := bytes.NewBufferString("")
	bb.WriteString("#!/usr/bin/env python\n")
	bb.WriteString("# -*-coding:utf-8-*-\n")
	bb.WriteString(`"""`)
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString("!!!DO NOT EDIT!!!")
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString("++++++++++++++++++++++++++++++++++++++++++++++")
	bb.WriteString("\n")
	bb.WriteString("+ author   +     swagen                      +")
	bb.WriteString("\n")
	bb.WriteString("++++++++++++++++++++++++++++++++++++++++++++++")
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf("+ datetime +     %s   +", time.Now().Format(time.RFC3339)))
	bb.WriteString("\n")
	bb.WriteString("++++++++++++++++++++++++++++++++++++++++++++++")
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString("***")
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString(sg.Info.Description)
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString("***")
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString(`"""`)
	bb.WriteString("\n")
	bb.WriteString("\n")


	bb.WriteString("import apis\n")
	bb.WriteString("import datetime\n")
	bb.WriteString("import json\n")
	bb.WriteString(fmt.Sprintf("import %s\n", pkg))

	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf(`service_name = "%s@%s"`, pkg, version))
	bb.WriteString("\n")
	bb.WriteString("\n")
	bb.WriteString("\n")



	// definitions
	defs, defErr := generateDefinitions(pkg, sg)
	if defErr != nil {
		zlog.Log().With("definitions", "failed").Errorf("generate definitions failed, %v", defErr)
		err = defErr
		return
	}
	zlog.Log().With("definitions", "succeed").Debug("generate definitions succeed")

	bb.Write(defs)

	// paths
	paths, pathErr := generatePaths(pkg, sg)
	if pathErr != nil {
		zlog.Log().With("paths", "failed").Errorf("generate paths failed, %v", pathErr)
		err = defErr
		return
	}
	zlog.Log().With("paths", "succeed").Debug("generate paths succeed")

	bb.Write(paths)

	p = bb.Bytes()

	return
}
