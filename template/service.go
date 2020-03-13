package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplService = `package {{ ucdown (getFirstService .Services).Name }}

// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
// File Location: handler/{{ ucfirst (getFirstService .Services).Name }}.service.go

import  (
	pb "{{ .Src }}/grpc/pb/{{ .GoPackage }}"
	"context"
	// "fmt"
	empty "github.com/golang/protobuf/ptypes/empty"
	{{- if .Elastic }}
	core "{{ .Src }}/core"
	{{- end}}
	{{- if .TimeStamp }}
	ptypes "github.com/golang/protobuf/ptypes"
	{{- end}}
)

{{- range $service := .Services }}
type {{ ucfirst $service.Name }}Service struct{
	repo *MasterRepository
	{{if $service.Elastic }}
	{{- range $msg := $service.AllMessage }}
	{{- if $msg.IsElastic }}
	es{{ ucfirst $msg.Name }} core.ESModule
	{{- end}}
	{{- end}}
	{{- end}}
}

// {{ ucfirst $service.Name }}Svc for service singleton
var {{ $service.Name }}Svc *{{ ucfirst $service.Name }}Service

// New{{ ucfirst $service.Name }}Service for new repository service
{{if $service.Elastic }}
func New{{ ucfirst $service.Name }}Service(repo *MasterRepository,{{ $service.MessageAllEs }}) *{{ ucfirst $service.Name }}Service {
{{- else}}
func New{{ ucfirst $service.Name }}Service(repo *MasterRepository) *{{ ucfirst $service.Name }}Service {
{{- end}}
	if {{ $service.Name }}Svc == nil {
		{{ $service.Name }}Svc = &{{ ucfirst $service.Name }}Service{
			repo,
			{{- range $msg := $service.AllMessage }}
			{{- if $msg.IsElastic }}
			es{{ ucfirst $msg.Name }},
			{{- end}}
			{{- end}}
		}
	}
	return {{ $service.Name }}Svc
}

{{- range $method := $service.Methods }}
// {{ ucfirst $method.Name }} method declare by generated code
{{- if eq $method.Input "empty"}}
func (svc *{{ ucfirst $service.Name }}Service) {{ ucfirst $method.Name }}(ctx context.Context, in *empty.Empty) (*pb.{{ ucfirst $method.Output }}, error) {
{{- else}}
func (svc *{{ ucfirst $service.Name }}Service) {{ ucfirst $method.Name }}(ctx context.Context, in *pb.{{ ucfirst $method.Input }}) (*pb.{{ ucfirst $method.Output }}, error) {
{{- end}}

{{- if eq $method.Input "empty"}}
	return &pb.{{ ucfirst $method.Output }}{}, nil
{{- else}}
{{- if $method.IsAgregator}}
	model := &{{ ucfirst $method.AgregatorMessage.Name }}{}
{{- range $field := $method.InputWithAgregator.Fields }}
{{- if eq $field.IgnoreGorm false}}
{{- if eq $field.TypeDataGo "time.Time"}}
	time{{ ucfirst $field.Name }}, errTime{{ ucfirst $field.Name }} := ptypes.Timestamp(in.{{ ucfirst $field.Name }})

	if errTime{{ ucfirst $field.Name }} == nil {
		model.{{ ucfirst $field.Name }} = time{{ ucfirst $field.Name }}
	}
{{- else }}
	model.{{ ucfirst $field.Name }} = in.{{ ucfirst $field.Name }}
{{- end}}
{{- end}}
{{- end}}
{{- if eq $method.AgregatorFunction "GetAll"}}
	_, err := svc.repo.{{ ucfirst $method.AgregatorMessage.Name }}.{{ $method.AgregatorFunction }}(model)
{{- else }}
	res, err := svc.repo.{{ ucfirst $method.AgregatorMessage.Name }}.{{ $method.AgregatorFunction }}(model)
{{- end}}

	resp := &pb.{{ ucfirst $method.Output }}{}
{{- if eq $method.AgregatorFunction "GetAll"}}
{{- else}}
{{- range $field := $method.IO.Fields }}
{{- if eq $field.IgnoreGorm false }}

{{- if eq $field.TypeDataGo "time.Time"}}
	protoTimeResp{{ ucfirst $field.Name }}, errProtoTimeResp{{ ucfirst $field.Name }} := ptypes.TimestampProto(res.{{ ucfirst $field.Name }})

	if errProtoTimeResp{{ ucfirst $field.Name }} == nil {
		resp.{{ ucfirst $field.Name }} = protoTimeResp{{ ucfirst $field.Name }}
	}
{{- else }}
	resp.{{ ucfirst $field.Name }} = res.{{ ucfirst $field.Name }}
{{- end}}

{{- end}}
{{- end}}
{{- end}}
	return resp, err

{{- else}}
	return &pb.{{ ucfirst $method.Output }}{}, nil
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}

`

var ListService = lib.List{
	FileType: ".service.go",
	Template: tmplService,
	Location: "./%s/",
	Lang:     "go",
}
