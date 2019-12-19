package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplYaml = `service: enterprise-app-{{ ucdown .GoPackage }}
stack: go
ignorePath:
  - vendor
apiVersion: v1
disableHealthCheck: false
dependencies:
  databases:
    - name: "enterprise{{ ucdown .GoPackage }}"
      type: mysql
`

var ListYaml = lib.List{
	FileType: ".yaml",
	Template: tmplYaml,
	Location: "./",
	Lang:     "yaml",
}
