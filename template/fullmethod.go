package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplFullMethod = `
// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}
package {{ .GoPackage }}

import  (
	core "github.com/zokypesch/proto-lib/core"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	{{ $service.Name }}_{{ $method.Name }} = "/{{ $.Package }}.{{ ucdown $service.Name }}/{{ $method.Name }}"
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
	Pattern_health_check_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"_health"}, ""))
	// Deprecated on v2 grpc gateway
	// Pattern_health_check_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"_health"}, "", runtime.AssumeColonVerbOpt(true)))
)

func InitCallGRPC() error {
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	forward_{{ $service.Name }}_{{ $method.Name }}_0 = core.LocalForward
{{- end}}
{{- end}}
	// Deprecated on v2 grpc gateway
	// runtime.HTTPError = CustomHTTPError

	return nil
}
`

var ListFullMethod = lib.List{
	FileType: ".custom.pb.go",
	Template: tmplFullMethod,
	Location: "",
	Lang:     "go",
}
