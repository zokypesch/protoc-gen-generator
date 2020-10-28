package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplMainGov4 = `package main

// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
// File Location: {{ ucfirst (getFirstService .Services).Name }}.main.go

import (
	"fmt"
	"log"
	{{ ucdown (getFirstService .Services).Name }} "{{ .Src }}/{{ ucdown (getFirstService .Services).Name }}"
	sv "{{ .Src }}/handler"

	pb "{{ .Src }}/grpc/pb/{{ .GoPackage }}"
	core "github.com/zokypesch/proto-lib/core"
	config "{{ .Src }}/config"

	"net"
	
	{{- if .Elastic }}
	domain "{{ .Src }}/{{ ucdown (getFirstService .Services).Name }}"
	{{- end}}
	middleware "gitlab.com/prakerja/shared/middleware"
	shared "gitlab.com/prakerja/shared/client"
)

func main() {

	cfg := config.Get()

	userClient := shared.NewUserClientService(cfg.UserApiAddress, cfg.UserApiPassword)

	db := core.InitDB(cfg.DBAddress, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBLog, cfg.DBMaxOpen, cfg.DBMaxIdle)
	lis, errList := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPORT))

	if errList != nil {
		log.Fatal(errList)
	}

	masterRepo := {{ ucdown (getFirstService .Services).Name }}.NewMasterRepoService(db)
	
{{- range $msg := .Messages }}
{{- if $msg.IsElastic }}
	es{{ ucfirst $msg.Name }} := core.NewEsCore(cfg.ESAddress, "{{ $msg.Name }}ing", domain.Mapping{{ ucfirst $msg.Name }}, "{{ $msg.Name }}")
{{- end}}
{{- end}}

{{- if .Elastic }}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo,{{ .MessageAll }})
{{- else}}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo)
{{- end}}
	
	handler := sv.New{{ ucfirst (getFirstService .Services).Name }}(masterService)

	auth := middleware.NewAuthInterceptor(cfg.INTERNALPASSWORD, cfg.InternalKey, userClient)
	interceptor := auth.GetUnaryCustom([]string{
		{{- range $w := .WhiteList }}
			pb.{{ $w.ServiceName }}_{{ $w.Name }},
		{{- end }}
	})

	svcName := fmt.Sprintf("%s_%s", cfg.ServiceName, cfg.Env)
	server := core.RegisterGRPCWithPrometh(svcName, interceptor)
	
	pb.Register{{ ucfirst (getFirstService .Services).Name }}Server(server, handler)
	core.RegisterPrometheus(server, 9092)
	
	go func() {{ unescape "{" }}
		if err := server.Serve(lis); err != nil {{ unescape "{" }}
			log.Fatalf("failed to serve: %v", err)
		{{ unescape "}" }}
	{{ unescape "}" }}()

	log.Println("starting server")

	if err := core.RunHTTPWithCustomMatcher(pb.InitCallGRPC, pb.Register{{ ucfirst (getFirstService .Services).Name }}HandlerFromEndpoint, middleware.CustomMatcher([]string{cfg.InternalKey}), cfg.GRPCClient, cfg.GRPCPORT, cfg.PORT); err != nil {
		log.Fatal(err)
	}
}

`

var ListMainv4 = lib.List{
	FileType: ".main.go",
	Template: tmplMainGov4,
	Location: "./",
	Lang:     "go",
}
