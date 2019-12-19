package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplEnv = `export PORT=80
export GRPCPORT=8080
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_HOST=enterprise-app-{{ ucdown .GoPackage }}-stag-mysql.statefulset.svc.cluster.local
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_PORT=3306
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_USERNAME=root
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_PASSWORD=ErGeRj45
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_DB_NAME=enterprise{{ ucdown .GoPackage }}
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_MAX_IDLE_CONNECTION=5
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_MAX_LIFETIME_CONNECTION=200
export ENTERPRISE{{ upper .GoPackage }}_MYSQL_MAX_OPEN_CONNECTION=10

`

var ListEnv = lib.List{
	FileType: ".env",
	Template: tmplEnv,
	Location: "./",
	Lang:     "env",
}
