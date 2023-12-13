package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const deprecationComment = "// Deprecated: Do not use."

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_grain.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	generateHeader(gen, g, file)
	generateContent(gen, g, file)
}

func generateHeader(gen *protogen.Plugin, g *protogen.GeneratedFile, file *protogen.File) {
	g.P("// Code generated by protoc-gen-grain. DO NOT EDIT.")
	g.P("// versions:")
	g.P("//  protoc-gen-grain ", version)
	protocVersion := "(unknown)"
	if v := gen.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		if s := v.GetSuffix(); s != "" {
			protocVersion += "-" + s
		}
	}
	g.P("//  protoc           ", protocVersion)
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
}

func generateContent(gen *protogen.Plugin, g *protogen.GeneratedFile, file *protogen.File) {
	g.P("package ", file.GoPackageName)
	g.P()

	if len(file.Services) == 0 {
		return
	}

	for _, service := range file.Services {
		generateService(service, file, g)
	}
}

func generateService(service *protogen.Service, file *protogen.File, g *protogen.GeneratedFile) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	sd := &serviceDesc{
		Name: service.GoName,
	}

	for i, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}

		md := &methodDesc{
			Name:   method.GoName,
			Input:  g.QualifiedGoIdent(method.Input.GoIdent),
			Output: g.QualifiedGoIdent(method.Output.GoIdent),
			Index:  i,
		}

		sd.Methods = append(sd.Methods, md)
	}

	if len(sd.Methods) != 0 {
		g.P(sd.execute())
	}
}
