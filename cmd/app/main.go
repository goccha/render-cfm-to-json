package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/goccha/render-cfm-to-json/pkg/cloudformation"
	"github.com/goccha/render-cfm-to-json/pkg/temporaries"
	"os"
)

type Option struct {
	Name       string `validate:"required"`
	Cfm        string `validate:"required"`
	Parameters string
}

var opt = &Option{}

func init() {
	flag.StringVar(&opt.Name, "name", "", "resource name")
	flag.StringVar(&opt.Parameters, "json-params", "", "CloudFormation Parameters")
	flag.StringVar(&opt.Cfm, "cfm", "", "CloudFormation file path")
	flag.Parse()
}

func main() {
	var err error
	var params map[string]interface{}
	if opt.Parameters != "" {
		if err = json.Unmarshal([]byte(opt.Parameters), &params); err != nil {
			abort(err)
		}
	}
	renderer := &cloudformation.Renderer{
		Name:   opt.Name,
		Params: params,
	}
	var indent, min []byte
	if indent, min, err = renderer.Render(opt.Cfm); err != nil {
		abort(err)
	}
	printGroup("json-body", string(indent))
	setOutput("json-body", string(min))

	var tmpFile *temporaries.File
	if tmpFile, err = temporaries.Open(opt.Cfm); err != nil {
		abort(err)
	} else {
		if _, err = tmpFile.Write(indent); err != nil {
			abort(err)
		}
		setOutput("json-file", tmpFile.Name())
	}
}

func abort(err error) {
	fmt.Printf("%+v\n", err)
	os.Exit(1)
}

func printGroup(name, value string) {
	fmt.Printf("::group::%s\n", name)
	fmt.Printf("%s\n", value)
	fmt.Println("::endgroup::")
}

func setOutput(name, value string) {
	fmt.Printf("::set-output name=%s::%s\n", name, value)
}
