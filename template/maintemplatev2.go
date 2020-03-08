package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplMainGov2 = `package main

// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
// File Location: r{{ ucfirst (getFirstService .Services).Name }}.main.go

import (
	"github.com/jinzhu/gorm"
	"log"
	{{ ucdown (getFirstService .Services).Name }} "{{ .Src }}/{{ ucdown (getFirstService .Services).Name }}"
	sv "{{ .Src }}/handler"

	_ "github.com/go-sql-driver/mysql"
	pb "{{ .Src }}/grpc/pb/{{ .GoPackage }}"

	"google.golang.org/grpc"
	"os"
	"time"
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"net"
	"net/http"
	"fmt"
	{{- if .Elastic }}
	core "{{ .Src }}/core"
	config "{{ .Src }}/config"
	{{- end}}
)

func runHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterSimpleHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8081", mux)
}


// CustomLogger for custome logger
func CustomLogger(code codes.Code) logrus.Level {
	if code == codes.OK {
		return logrus.InfoLevel
	}

	return logrus.WarnLevel
}

func InitDB(address string, dbName string) *gorm.DB {
	dbUser := "root"
	dbPass := "ErGeRj45"
	dbEndpoint := address
	dbPort := "3306"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbEndpoint, dbPort, dbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	db.LogMode(true)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(10))
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(50)
	db.SingularTable(true)
	return db
}

func main() {

	{{- if .Elastic }}
	cfg := config.Get()
	{{- end}}

	db := InitDB("fec-ticketing-stag-mysql.statefulset.svc.cluster.local", "ticketing")
	lis, errList := net.Listen("tcp", ":8080")

	if errList != nil {
		log.Fatal(errList)
	}

	masterRepo := {{ ucdown (getFirstService .Services).Name }}.NewMasterRepoService(db)
	
{{- range $msg := .Messages }}
{{- if $msg.IsElastic }}
	es{{ ucfirst $msg.Name }} := core.NewEsCore(cfg.ESAddress, "{{ $msg.Name }}ing", {{ ucdown (getFirstService .Services).Name }}.Mapping{{ ucfirst $msg.Name }}, "{{ $msg.Name }}")
{{- end}}
{{- end}}

{{- if .Elastic }}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo,{{ .MessageAll }})
{{- else}}
	masterService := {{ ucdown (getFirstService .Services).Name }}.New{{ ucfirst (getFirstService .Services).Name }}Service(masterRepo)
{{- end}}
	
	handler := sv.New{{ ucfirst (getFirstService .Services).Name }}(masterService)
	
	logger := &logrus.Logger{}
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	customFunc := CustomLogger

	logrusEntry := logrus.NewEntry(logger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
			),
		),
	)

	pb.Register{{ ucfirst (getFirstService .Services).Name }}Server(server, handler)

	go func() {{ unescape "{" }}
		if err := server.Serve(lis); err != nil {{ unescape "{" }}
			log.Fatalf("failed to serve: %v", err)
		{{ unescape "}" }}
	{{ unescape "}" }}()

	log.Println("starting server")

	if err := runHTTP(); err != nil {
		log.Fatal(err)
	}
}

`

var ListMainv2 = lib.List{
	FileType: ".main.go",
	Template: tmplMainGov2,
	Location: "./",
	Lang:     "go",
}
