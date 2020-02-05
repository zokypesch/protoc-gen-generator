package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplMainGo = `package main

// Code generated by sangkuriang  protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
// File Location: r{{ ucfirst (getFirstService .Services).Name }}.main.go

import (
	// "github.com/jinzhu/gorm"
	"log"

	"{{ .Src }}/config"
	// core "{{ .Src }}/core"
	// model "{{ .Src }}/model"
	{{ ucdown (getFirstService .Services).Name }} "{{ .Src }}/{{ ucdown (getFirstService .Services).Name }}"
	sv "{{ .Src }}/handler"
	"gitlab.com/ruangguru/source/shared-lib/go/conn/mysql"
	logInt "gitlab.com/ruangguru/source/shared-lib/go/middleware/grpc/logger"
	recInt "gitlab.com/ruangguru/source/shared-lib/go/middleware/grpc/recovery"
	morse "gitlab.com/ruangguru/source/shared-lib/go/morse"

	_ "github.com/go-sql-driver/mysql"
	pb "{{ .Src }}/grpc/pb/{{ .GoPackage }}"
	{{- if .Elastic }}
	core "{{ .Src }}/core"
	model "{{ .Src }}/model"
	{{- end}}
)

func main() {
	cfg := config.Get()

	db := mysql.Init(cfg.PREFIX)

	masterRepo := {{ ucdown (getFirstService .Services).Name }}.NewMasterRepoService(db)
{{- range $msg := .Messages }}
{{- if $msg.IsElastic }}
	es{{ ucfirst $msg.Name }} := core.NewEsCore(cfg.ESAddress, "{{ $msg.Name }}ing", model.Mapping{{ ucfirst $msg.Name }}, "{{ $msg.Name }}")
{{- end}}
{{- end}}

{{- if .Elastic }}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo,{{ .MessageAll }})
{{- else}}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo)
{{- end}}

	handler := sv.New{{ ucfirst (getFirstService .Services).Name }}(masterService)
	// masterService.Auth.GetAll(context.TODO())
	svc := morse.NewService(
		morse.GRPCPort(cfg.GRPCPORT),
		morse.RESTPort(cfg.PORT),
		morse.EnablePrometheus(),
	)

	loggerInt := logInt.UnaryServerInterceptor()
	recoverInt := recInt.UnaryServerInterceptor()

	// Assign unary interceptor to server
	svc.UseServerUnaryInterceptor(loggerInt, recoverInt)
	pb.Register{{ ucfirst (getFirstService .Services).Name }}(svc, handler)

	if err := {{ unescape "<-" }}svc.RunServers(); err != nil {
		log.Fatal("err", err)
	}
}

`

var ListMain = lib.List{
	FileType: ".main.go",
	Template: tmplMainGo,
	Location: "./",
	Lang:     "go",
}
