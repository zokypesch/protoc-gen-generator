package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplConfig = `package config

import "github.com/kelseyhightower/envconfig"

// Config struct of configuration
type Config struct {
	ServiceName      string   'envconfig:"SERVICE_NAME" default:"enterprise-app-{{ ucdown .GoPackage }}"'
	PORT             string   'envconfig:"PORT" default:"80"'
	GRPCPORT         string   'envconfig:"GRPC_PORT" default:"8080"'
	{{- if .Elastic }}
	ESAddress        string   'envconfig:"ES_FEC_APPS" default:"http://fec-ticketing-stag-es.statefulset.svc.cluster.local:9200"'
	{{- end}}
	PREFIX           string   'envconfig:"PREFIX" default:"enterprise{{ ucdown .GoPackage }}"'
	INTERNALPASSWORD string   'envconfig:"INTERNAL_PASSWORD" default:"INTERNALAPIPASSWORD"'
}

// singleton of data
var data *Config

// Get configuration of data
func Get() *Config {
	if data == nil {
		data = &Config{}
		envconfig.MustProcess("", data)
	}

	// returing configuration
	return data
}

`

var ListConfig = lib.List{
	FileType:     ".config.go",
	Template:     tmplConfig,
	Location:     "./config/",
	Lang:         "go",
	ReplaceQuote: true,
}
