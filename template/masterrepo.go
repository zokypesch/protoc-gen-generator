package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplMasterRepoGo = `package repo

// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
// File Location: repo/{{ ucfirst (getFirstService .Services).Name }}.master.repository.go

import (
	"github.com/jinzhu/gorm"
)

// MasterRepository for factory repository
type MasterRepository struct {
{{- range $msg := .Messages }}
{{- if $msg.IsRepository}} 
	{{ ucfirst $msg.Name }}		{{ ucfirst $msg.Name }}Service
{{- end}}
{{- end}}
}
var masterRepo *MasterRepository

// NewMasterRepoService for new repository event service
func NewMasterRepoService(db *gorm.DB) *MasterRepository {
	if masterRepo == nil {
		masterRepo = &MasterRepository{
		{{- range $msg := .Messages }}
		{{- if $msg.IsRepository}} 
			{{ ucfirst $msg.Name }}:		New{{ ucfirst $msg.Name }}Service(db),
		{{- end}}
		{{- end}}
		}
	}
	return masterRepo
}
`

var ListMasterRepoGolang = lib.List{
	FileType: ".master.repository.go",
	Template: tmplMasterRepoGo,
	Location: "./repo/",
	Lang:     "go",
}
