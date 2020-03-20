package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplPostman = `{
	"info": {
		"_postman_id": "{{ generateGuid }}",
		"name": "{{ (getFirstService .Services).Name }}",
		"description": "API FOR {{ ucfirst .FileName }}",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [
		{{- range $service := .Services }} {{- range $method := $service.Methods }}
		{
			"name": "http://147.139.166.237{{ strReplaceParam $method.URLPath }}",
			"request": {
				"method": "{{ toupper $method.HttpMode }}",
				"header": [
					{
						"key": "Authorization",
						"value": "INTERNALAPIPASSWORD",
						"type": "text"
					}{{- if eq $method.HttpMode "get" }} {{- else }},
					{
						"key": "Content-Type",
						"value": "application/json",
						"type": "text"
					}
					{{- end }}
				],
				{{- if eq $method.HttpMode "get" }}
				"body": {
					"mode": "raw",
					"raw": ""
				},
				{{- else }}
				"body": {
					"mode": "raw",
					"raw": "{{ unescape "{" }}{{ range $field := $method.InputMessage.Fields }}\n\t\"{{ ucdown $field.NameGo }}\": \"\"{{if $field.ExtraComma }},{{- end }}{{- end }}}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				{{- end }}
				"url": {
					"raw": "http://147.139.166.237{{ strReplaceParam $method.URLPath }}",
					"protocol": "http",
					"host": [
						"147",
						"139",
						"166",
						"237"
					],
					"path": [
						{{- range $opt := $method.PathPostman }}
						"{{ $opt.Name }}"{{ $opt.Extra }}
						{{- end }}
					]
				}
			},
			"response": []
		}{{- if $method.ExtraComma }},{{- end }}
		{{- end }}
		{{- end}}
	]
}

`

var ListPostman = lib.List{
	FileType:     ".postman_collection.json",
	Template:     tmplPostman,
	Location:     "./",
	Lang:         "json",
	ReplaceQuote: true,
}
