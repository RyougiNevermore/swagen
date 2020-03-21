package main

import (
	"errors"
	"github.com/pharosnet/swagen/py"
	"github.com/pharosnet/swagen/source"
	"github.com/pharosnet/swagen/zlog"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

const version = "v0.0.1"

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source_type",
				Value:    "",
				Required: true,
				Usage:    "swagger doc source type, --source_type=file, --source_type=dir, --source_type=http",
			},

			&cli.StringFlag{
				Name:     "source_path",
				Value:    "",
				Required: true,
				Usage:    "swagger doc source path, --source_path=/foo/bar.json --source_path=/foo/bar --source_path=http://foo.com/bar ",
			},

			//&cli.GenericFlag{
			//	Name:  "content_type",
			//	Value: &ContentGenericFlag{},
			//	Usage: "swagger doc file type, json or yaml",
			//},

			&cli.StringFlag{
				Name:     "output",
				Value:    "",
				Required: true,
				Usage:    "generated code's output dir",
			},


			&cli.StringFlag{
				Name:     "language",
				Value:    "",
				Required: true,
				Usage:    "generated code's language, such as python, golang, java",
			},

			&cli.BoolFlag{
				Name:     "debug",
				Required: false,
				Usage:    "--debug=true",
			},
		},

		Version: version,
		Action:  execute,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

func execute(ctx *cli.Context) (err error) {

	debug := ctx.Bool("debug")

	zlog.New(debug, version)

	sgs, sgsErr := source.Fetch(ctx.String("source_type"), ctx.String("source_path"))
	if sgsErr != nil {
		zlog.Log().With("fetch", "failed").Error(sgsErr)
		return
	}

	output := ctx.String("output")
	language := strings.ToLower(ctx.String("language"))

	switch language {
	case "python":
		genErr := py.Generate(sgs, output)
		if genErr != nil {
			zlog.Log().With("generate", "failed").Error(genErr)
			return
		}

	default:
		zlog.Log().With("generate", "failed").Errorf("the %s language is not supported", language)
	}

	return
}

type ContentGenericFlag struct {
	v string
}

func (f *ContentGenericFlag) Set(value string) (err error) {

	value = strings.ToLower(strings.TrimSpace(value))

	if value != "json" && value != "yaml" {
		err = errors.New("content_type must be json or yaml")
		return
	}
	f.v = value

	return
}
func (f *ContentGenericFlag) String() (v string) {
	v = f.v
	return
}
