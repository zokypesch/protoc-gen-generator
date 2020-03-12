package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplFullMethod = `
// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}
package {{ .GoPackage }}

import  (
	core "github.com/zokypesch/proto-lib/core"
)

const (
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	{{ $service.Name }}_{{ $method.Name }} = "/{{ $.Package }}.{{ $service.Name }}/{{ $method.Name }}"
{{- end}}
{{- end}}
)

var (
	FullMethods = []string{{ "{" }}
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
		{{ $service.Name }}_{{ $method.Name }},
{{- end}}
{{- end}}
	{{ "}" }}
)

func InitCallGRPC() {
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	forward_{{ $service.Name }}_{{ $method.Name }}_0 = core.LocalForward
{{- end}}
{{- end}}
	runtime.HTTPError = core.CustomHTTPError
}
`

var ListFullMethod = lib.List{
	FileType: ".custom.pb.go",
	Template: tmplFullMethod,
	Location: "",
	Lang:     "go",
}
