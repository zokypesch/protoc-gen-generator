package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplReadme = `# {{ ucfirst (getFirstService .Services).Name }}

## How to build to docker ?
docker build -t {your_username}/{repo_name} .

## login to docker ?
username: {your_username}
password: {your_password}

## push to cloud ?
docker push {your_username}/{repo_name}

## how to auto generate proto using db ??
proto-tools -cmd=gen-proto-db host=localhost name={{ ucdown (getFirstService .Services).Name }} user=root password=

## How to regenerate ??
please split your own code within auto generate code, because auto generated code will be replace your own code or you can mark it by using type in code
sangkuriang grpc/proto/{{ ucdown (getFirstService .Services).Name }} {{ ucdown (getFirstService .Services).Name }} grpc/pb/{{ ucdown (getFirstService .Services).Name }} 

# ENDPOINTS
'''
{{- range $service := .Services }}
{{- range $method := $service.Methods }}

Name {{ $method.Name }}
Endpoint {{ $method.URLPath }}
HTTP {{ $method.HttpMode }}
Request: {{ $method.InputMessage.Name }}
{
	{{ range $field := $method.InputMessage.Fields }}
	{{- if allowRequest $field.Name }}
	"{{ ucdown $field.NameGo }}": {{ $field.TypeDataGo }}  {{- if eq $field.Tag "" }} {{- else}} '{{ unescape $field.Tag }}' {{- end}},
	{{- end }}
	{{- end}}
}
Response: {{ $method.OutputMessage.Name }}
{
	{{ range $field := $method.OutputMessage.Fields }}
	{{- if $field.IsFieldMessage }}
	{{- if eq $field.MessageTo.Name "" }}
	"{{ ucdown $field.NameGo }}": {{ $field.TypeDataGo }},
	{{- else }}
	"{{ ucdown $field.NameGo }}": {
		{{ range $sub := $field.MessageTo.Fields }}
		"{{ ucdown $sub.NameGo }}": {{ $sub.TypeDataGo }},
		{{- end}}
	}
	{{- end }}
	{{- else }}
	"{{ ucdown $field.NameGo }}": {{ $field.TypeDataGo }},
	{{- end }}
	{{- end}}
}

{{- end }}
{{- end }}

'''


# TABLES
'''
{{- range $msg := .Messages }}
{{- if $msg.IsRepository }}

// {{ ucfirst $msg.Name }} for struct info
type {{ ucfirst $msg.Name }} struct {
{{- range $field := $msg.Fields }}
	{{ $field.NameGo }}			{{ $field.TypeDataGo }} {{- if eq $field.Tag "" }} {{- else}} '{{ unescape $field.Tag }}' {{- end}}
{{- end}}
}

{{- end }}
{{- end}}
'''

`

var ListReadme = lib.List{
	FileType:     "Readme",
	Template:     tmplReadme,
	Location:     "./",
	Lang:         "readme",
	ReplaceQuote: true,
}
