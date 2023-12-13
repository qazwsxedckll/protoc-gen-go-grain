package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/grain.tmpl
var grainTemplate string

type serviceDesc struct {
	Name    string // Greeter
	Methods []*methodDesc
}

type methodDesc struct {
	Name   string
	Input  string
	Output string
	Index  int
}

func (s *serviceDesc) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("grain").Parse(strings.TrimSpace(grainTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return strings.Trim(buf.String(), "\r\n")
}
